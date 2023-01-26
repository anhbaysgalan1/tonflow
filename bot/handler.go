package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/makiuchi-d/gozxing"
	qrScan "github.com/makiuchi-d/gozxing/qrcode"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"github.com/xssnick/tonutils-go/tlb"
	"image/jpeg"
	"net/http"
	"strings"
	"tonflow/model"
)

// Returns user object from cache or database in case of no cache.
// Also returns "true" if user is existed already else creates and returns new user with "false" second value.
func (bot *Bot) getTonflowUser(ctx context.Context, tgUser *tgBotAPI.User) (*model.User, bool, error) {
	user, err := bot.redis.GetUserCache(ctx, tgUser.ID)
	if err != nil && err != redis.Nil {
		log.Error(err)
		return nil, false, err
	}

	wallet := &model.Wallet{}
	isExisted := true

	if err == redis.Nil {
		user, err = bot.storage.GetUser(ctx, tgUser.ID)
		switch {
		case err != nil && !errors.Is(err, pgx.ErrNoRows):
			log.Error(err)
			return nil, false, err
		case errors.Is(err, pgx.ErrNoRows):
			wlt, err := bot.ton.NewWallet()
			if err != nil {
				log.Error(err)
				return nil, false, err
			}

			log.Debugf("generated seed: %v", wlt.Seed)

			err = bot.storage.AddUser(ctx, tgUser)
			if err != nil {
				log.Error(err)
				return nil, false, err
			}

			err = bot.storage.AddWallet(ctx, wlt, tgUser.ID)
			if err != nil {
				log.Error(err)
				return nil, false, err
			}

			wallet = wlt
			isExisted = false
		default:
			wallet, err = bot.storage.GetWallet(ctx, user.Wallet.Address)
			if err != nil {
				log.Error(err)
				return nil, false, err
			}
		}

		user = &model.User{
			ID:           tgUser.ID,
			Username:     tgUser.UserName,
			FirstName:    tgUser.FirstName,
			LastName:     tgUser.LastName,
			LanguageCode: tgUser.LanguageCode,
			Wallet:       wallet,
			StageData:    &model.StageData{},
		}

		err = bot.redis.SetUserCache(ctx, user)
		if err != nil {
			log.Error(err)
			return nil, true, err
		}
	}

	return user, isExisted, nil
}

func (bot *Bot) cmdStart(update tgBotAPI.Update, user *model.User, isExisted bool) {
	chatID := update.Message.Chat.ID
	firstName := update.Message.From.FirstName
	address := user.Wallet.Address

	qr, err := qrcode.Encode(address, qrcode.Medium, 512)
	if err != nil {
		log.Error(err)
		return
	}

	caption := ""
	switch isExisted {
	case true:
		caption = fmt.Sprintf(WelcomeExistedUser, firstName, address)
	case false:
		caption = fmt.Sprintf(WelcomeNewUser, firstName, address)
	}

	if err = bot.sendPhoto(chatID, qr, caption, inlineMainKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) acceptSendingAddress(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.Message.Chat.ID
	photos := update.Message.Photo
	addr := update.Message.Text

	if len(photos) != 0 {
		index, size := 0, 0
		for i, v := range photos {
			if v.FileSize > size {
				index = i
			}
		}

		fileURL, err := bot.api.GetFileDirectURL(photos[index].FileID)
		if err != nil {
			log.Error(err)
			return
		}

		resp, err := http.Get(fileURL)
		if err != nil {
			log.Error(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error(err)
			return
		}

		img, err := jpeg.Decode(resp.Body)
		if err != nil {
			log.Error(err)
			return
		}

		// prepare BinaryBitmap
		bmp, err := gozxing.NewBinaryBitmapFromImage(img)
		if err != nil {
			log.Error(err)
			return
		}

		// decode image
		qrReader := qrScan.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)
		if err != nil {
			log.Warnf("failed to decode qr %v: %v", fileURL, err)
			if err = bot.sendText(chatID, InvalidQR, nil); err != nil {
				log.Error(err)
				return
			}
			return
		}

		addr = result.String()
	}

	err := bot.ton.ValidateWallet(addr)
	if err != nil {
		log.Warnf("failed to validate address %v: %v", addr, err)
		if err = bot.sendText(chatID, InvalidWallet, nil); err != nil {
			log.Error(err)
			return
		}
		return
	}

	// set user stage in cache
	user.StageData.Stage = model.AmountWait
	user.StageData.AddressToSend = addr
	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, AskAmount, nil); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) acceptSendingAmount(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	wallet := user.Wallet.Address

	amount := strings.ReplaceAll(text, " ", "")
	amount = strings.ReplaceAll(amount, ",", ".")

	amountInCoins, err := tlb.FromTON(amount)
	if err != nil {
		log.Warnf("failed to convert %v to coins: %v", amountInCoins, err)
		if err = bot.sendText(chatID, InvalidAmount, struct{}{}); err != nil {
			log.Error(err)
			return
		}
		return
	}

	balance, err := bot.ton.GetWalletBalance(wallet)
	if err != nil {
		log.Error(err)
		return
	}

	balanceInCoins, err := tlb.FromTON(balance)
	if err != nil {
		log.Error()
		return
	}

	if balanceInCoins.NanoTON().Cmp(amountInCoins.NanoTON()) < 0 {
		txt := fmt.Sprintf(NotEnoughFunds, balance)
		if err = bot.sendText(chatID, txt, nil); err != nil {
			return
		}
		return
	}

	user.StageData.Stage = model.ConfirmationWait
	user.StageData.AmountToSend = amount
	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	txt := fmt.Sprintf(SendingConfirmation, user.StageData.AddressToSend, amount)
	if err = bot.sendText(chatID, txt, inlineConfirmKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) acceptComment(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.Message.Chat.ID
	addr := user.StageData.AddressToSend
	amount := user.StageData.AmountToSend
	comment := update.Message.Text

	user.StageData.Stage = model.ConfirmationWait
	user.StageData.Comment = amount
	err := bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	txt := fmt.Sprintf(SendingConfirmation, addr, amount) + fmt.Sprintf(Comment, comment)
	if err = bot.sendText(chatID, txt, inlineConfirmWithCommentKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineReceiveCoins(update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	qr, err := qrcode.Encode(address, qrcode.Medium, 512)
	if err != nil {
		log.Error(err)
		return
	}

	caption := fmt.Sprintf("<code>%s</code>", address) + "\n\n" + ReceiveInstruction
	if err := bot.sendPhoto(chatID, qr, caption, inlineMainKeyboard); err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineSendCoins(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	balance, err := bot.ton.GetWalletBalance(address)
	if err != nil {
		log.Error(err)
		return
	}

	if balance == "0" {
		if err = bot.sendText(chatID, NoFunds, nil); err != nil {
			log.Error(err)
			return
		}
		return
	}

	user.StageData.Stage = model.AddressWait
	if err = bot.redis.SetUserCache(ctx, user); err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, AskWallet, nil); err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineBalance(update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	balance, err := bot.ton.GetWalletBalance(address)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, fmt.Sprintf(Balance, balance), inlineBalanceKeyboard); err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineBalanceUpdate(update tgBotAPI.Update, user *model.User) {
	balance, err := bot.ton.GetWalletBalance(user.Wallet.Address)
	if err != nil {
		log.Error(err)
		return
	}

	newText := fmt.Sprintf(Balance, balance)
	notification := BalanceUpToDate

	if newText != update.CallbackQuery.Message.Text {
		msg := tgBotAPI.EditMessageTextConfig{
			BaseEdit: tgBotAPI.BaseEdit{
				ChatID:      user.ID,
				MessageID:   update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &inlineBalanceKeyboard,
			},
			Text:      newText,
			ParseMode: "HTML",
		}

		_, err = bot.api.Request(msg)
		if err != nil {
			log.Error(err)
			return
		}

		notification = BalanceUpdated
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, notification)
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineCancel(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	text := update.CallbackQuery.Message.Text

	newText := text + Canceled

	msg := tgBotAPI.EditMessageTextConfig{
		BaseEdit: tgBotAPI.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      newText,
		ParseMode: "HTML",
	}

	user.StageData = &model.StageData{
		Stage:         model.ZeroStage,
		AddressToSend: "",
		AmountToSend:  "",
	}

	err := bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = bot.api.Request(msg)
	if err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) inlineAddComment(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.CallbackQuery.Message.Chat.ID

	user.StageData.Stage = model.CommentWait
	err := bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, AskComment, nil); err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}
