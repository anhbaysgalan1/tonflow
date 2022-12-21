package bot

import (
	"context"
	"errors"
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/makiuchi-d/gozxing"
	qrScan "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"github.com/xssnick/tonutils-go/tlb"
	"image/jpeg"
	"net/http"
	"park-wallet/internal/service/bot/template"
	"park-wallet/internal/storage"
	"park-wallet/pkg"
	"strconv"
	"strings"
)

func (bot *Bot) handleUpdate(ctx context.Context, update tgBotAPI.Update) {
	switch {
	case update.Message == nil:
		bot.handleNilMessage(ctx, update)
	case update.Message != nil:
		bot.handleMessage(ctx, update)
	}
}

func (bot *Bot) handleNilMessage(_ context.Context, update tgBotAPI.Update) {
	if update.CallbackQuery != nil {
		callback := tgBotAPI.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		_, err := bot.api.Request(callback)
		if err != nil {
			bot.err(err, "send callback message")
		}
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update tgBotAPI.Update) {
	isExist, wallet, stage, err := bot.checkUser(ctx, update)
	if err != nil {
		bot.err(err, "check user")
		return
	}

	switch {
	case update.Message.IsCommand():
		bot.handleCommand(ctx, update, isExist, wallet)
	default:
		bot.handleUserMessage(ctx, update, isExist, wallet, stage)
	}
}

func (bot *Bot) checkUser(ctx context.Context, update tgBotAPI.Update) (bool, string, storage.Stage, error) {
	bot.sendTyping(update.Message.From.ID)

	flowUser := toFlowUser(update.SentFrom())

	userExist, err := bot.storage.CheckUser(ctx, flowUser)
	if err != nil {
		return false, "", "", err
	}

	userID := strconv.FormatInt(toFlowUser(update.SentFrom()).ID, 10)
	stage, err := bot.redis.GetStage(ctx, userID)
	if err != nil {
		return false, "", "", err
	}

	wlt, err := bot.storage.GetUserWallet(ctx, flowUser.ID)
	if err != nil {
		return false, "", "", err
	}

	if !userExist || (wlt == "" && err == nil) {
		wallet, err := bot.ton.NewWallet()
		if err != nil {
			return false, "", "", err
		}
		wallet.Seed, err = pkg.EncryptAES([]byte(bot.cryptoKey), wallet.Seed)
		if err != nil {
			return false, "", "", err
		}
		err = bot.storage.AddWallet(ctx, wallet, flowUser.ID)
		if err != nil {
			return false, "", "", err
		}
		wlt = wallet.Address
		return false, wlt, stage, nil
	}

	return true, wlt, stage, nil
}

func (bot *Bot) handleCommand(ctx context.Context, update tgBotAPI.Update, isExist bool, wallet string) {
	switch update.Message.Command() {
	case "start":
		qr, err := qrcode.Encode(wallet, qrcode.Medium, 512)
		if err != nil {
			bot.err(err, "generate qr")
			return
		}

		caption := ""
		switch isExist {
		case true:
			caption = fmt.Sprintf(template.WelcomeExistedUser, update.Message.From.FirstName, wallet)
		case false:
			caption = fmt.Sprintf(template.WelcomeNewUser, update.Message.From.FirstName, wallet)
		}

		if err := bot.sendPhoto(update.Message.Chat.ID, qr, caption, mainKeyboard); err != nil {
			bot.err(err, "send wallet address and qr")
		}

	}
}

func (bot *Bot) handleUserMessage(ctx context.Context, update tgBotAPI.Update, isExist bool, wallet string, stage storage.Stage) {
	chatID := update.Message.Chat.ID
	userID := strconv.FormatInt(toFlowUser(update.SentFrom()).ID, 10)
	text := update.Message.Text

	switch update.Message.Text {
	case template.BalanceButton:
		balance, err := bot.ton.GetWalletBalance(wallet)
		if err != nil {
			bot.err(err, "get wallet balance")
			return
		}

		if err = bot.sendText(chatID, fmt.Sprintf(template.Balance, balance), mainKeyboard); err != nil {
			bot.err(err, "send wallet balance")
		}

	case template.ReceiveButton:
		qr, err := qrcode.Encode(wallet, qrcode.Medium, 512)
		if err != nil {
			bot.err(err, "generate qr")
			return
		}

		caption := fmt.Sprintf("<pre>%s</pre>", wallet) + "\n\n" + template.ReceiveInstruction
		if err := bot.sendPhoto(chatID, qr, caption, mainKeyboard); err != nil {
			bot.err(err, "send wallet address and qr")
		}

	case template.SendButton:
		balance, err := bot.ton.GetWalletBalance(wallet)
		if err != nil {
			bot.err(err, "get balance")
			return
		}

		if balance == "0" {
			if err = bot.sendText(chatID, template.NoFunds, mainKeyboard); err != nil {
				bot.err(err, "not enough message")
				return
			}
			return
		}

		if err = bot.redis.SetStage(ctx, userID, storage.StageAmountWaiting); err != nil {
			bot.err(err, "set amount waiting stage")
			return
		}

		if err = bot.sendText(chatID, template.AskAmount, cancelKeyboard); err != nil {
			bot.err(err, "ask amount to send")
			return
		}

	default:
		switch text {
		case template.CancelButton:
			err := bot.redis.SetStage(ctx, userID, storage.StageUnset)
			if err != nil {
				bot.err(err, "set \"unset\" stage")
				return
			}

			if err = bot.sendText(chatID, template.Canceled, mainKeyboard); err != nil {
				bot.err(err, "send canceled message")
			}

		default:
			if stage == storage.StageAmountWaiting && update.Message.Text != "" {

				amount := strings.ReplaceAll(text, " ", "")
				amount = strings.ReplaceAll(amount, ",", ".")

				amountInCoins, err := tlb.FromTON(amount)
				if err != nil {
					if err := bot.sendText(chatID, template.InvalidAmount, cancelKeyboard); err != nil {
						bot.err(err, "send invalid amount")
						return
					}
					bot.err(err, "convert string to coins")
					return
				}

				balance, err := bot.ton.GetWalletBalance(wallet)
				if err != nil {
					bot.err(err, "get wallet balance")
					return
				}

				balanceInCoins, err := tlb.FromTON(balance)
				if err != nil {
					bot.err(err, "convert string to coins")
					return
				}

				if balanceInCoins.NanoTON().Cmp(amountInCoins.NanoTON()) < 0 {
					txt := fmt.Sprintf(template.NotEnoughFunds, balance)
					if err := bot.sendText(chatID, txt, cancelKeyboard); err != nil {
						bot.err(err, "send not enough amount")
						return
					}
					return
				}

				err = bot.redis.SetStage(ctx, userID, storage.StageWalletWaiting)
				if err != nil {
					bot.err(err, "set wallet waiting stage")
					return
				}

				if err := bot.sendText(chatID, template.AskWallet, cancelKeyboard); err != nil {
					bot.err(err, "ask wallet")
					return
				}
				return
			}

			if stage == storage.StageWalletWaiting && len(update.Message.Photo) != 0 {

				index := 0
				size := 0
				for i, v := range update.Message.Photo {
					if v.FileSize > size {
						index = i
					}
				}

				fileURL, err := bot.api.GetFileDirectURL(update.Message.Photo[index].FileID)
				if err != nil {
					bot.err(err, "get file direct url")
					return
				}

				log.Debug().Msg(fileURL)

				resp, err := http.Get(fileURL)
				if err != nil {
					bot.err(err, "get data from url")
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					bot.err(errors.New("bad status: "+resp.Status), "http status")
					return
				}

				img, err := jpeg.Decode(resp.Body)
				if err != nil {
					bot.err(err, "decode jpeg")
					return
				}

				// prepare BinaryBitmap
				bmp, err := gozxing.NewBinaryBitmapFromImage(img)
				if err != nil {
					bot.err(err, "prepare binary bitmap")
					return
				}

				// decode image
				qrReader := qrScan.NewQRCodeReader()
				result, err := qrReader.Decode(bmp, nil)
				if err != nil {
					bot.err(err, "decode qr")
					if err := bot.sendText(chatID, template.InvalidQR, cancelKeyboard); err != nil {
						bot.err(err, "send decode fail")
						return
					}
					return
				}

				receiverWallet := result.String()

				err = bot.ton.ValidateWallet(receiverWallet)
				if err != nil {
					bot.err(err, "validate receiver wallet")
					if err := bot.sendText(chatID, template.InvalidWallet, cancelKeyboard); err != nil {
						bot.err(err, "send validate fail")
						return
					}
					return
				}

				txt := fmt.Sprintf(template.SendingConfirmation, receiverWallet)
				if err := bot.sendText(chatID, txt, confirmKeyboard); err != nil {
					bot.err(err, "send validate fail")
					return
				}

			}

		}
	}
}

//func (bot *Bot) handleAdminMessage(ctx context.Context, update tgBotAPItgBotAPI.Update) {
//	bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
//
//	switch {
//	case len(update.Message.Photo) != 0:
//		ID := update.Message.Photo[0].FileID
//		err := bot.storage.AddPicture(ctx, ID, time.Now())
//		if err != nil {
//			log.Error().Err(err).Send()
//			bot.err(err, "failed to add picture in storage")
//			bot.sendText(update.Message.Chat.ID, "One of the pictures was not saved in the database", tgBotAPItgBotAPI.ReplyKeyboardMarkup{})
//		}
//	case update.Message.Text == "778":
//		bot.sendUploadingPhoto(update.Message.Chat.ID)
//
//		fileID, err := bot.storage.GetRandomPicture(ctx)
//		if err != nil {
//			log.Error().Err(err).Send()
//			bot.err(err, "failed to get random pic")
//			return
//		}
//
//		bot.sendPhoto(update.Message.Chat.ID, fileID)
//
//		time.AfterFunc(time.Second*5, func() {
//			bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID+1)
//		})
//	//case update.Message.Text == "55555":
//	//	IDs, err := bot.storage.GetAllPictures(ctx)
//	//	if err != nil {
//	//		bot.err(err, "failed to get random pic")
//	//		return
//	//	}
//	//
//	//	for _, v := range IDs {
//	//		bot.sendPhoto(update.Message.Chat.ID, v)
//	//		time.Sleep(time.Millisecond * 500)
//	//	}
//	default:
//
//	}
//}
