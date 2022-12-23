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
	"strconv"
	"strings"
	"tonflow/internal/service/bot/template"
	"tonflow/internal/storage"
	"tonflow/pkg"
)

func (bot *Bot) checkUser(ctx context.Context, update tgBotAPI.Update) (bool, string, storage.Stage, error) {
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

func (bot *Bot) msgReceivingOptions(chatID int64, wallet string) {
	qr, err := qrcode.Encode(wallet, qrcode.Medium, 512)
	if err != nil {
		bot.err(err, "generate qr")
		return
	}

	caption := fmt.Sprintf("<code>%s</code>", wallet) + "\n\n" + template.ReceiveInstruction
	if err := bot.sendPhoto(chatID, qr, caption, mainInlineKeyboard); err != nil {
		bot.err(err, "send wallet address and qr")
	}
}

func (bot *Bot) msgAskAddress(ctx context.Context, chatID int64, userID, wallet string) {
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

	if err = bot.redis.SetStage(ctx, userID, storage.StageWalletWaiting); err != nil {
		bot.err(err, "set wallet waiting stage")
		return
	}

	if err = bot.sendText(chatID, template.AskWallet, nil); err != nil {
		bot.err(err, "ask amount to send")
	}
}

func (bot *Bot) msgBalance(chatID int64, wallet string) {
	balance, err := bot.ton.GetWalletBalance(wallet)
	if err != nil {
		bot.err(err, "get wallet balance")
		return
	}

	if err = bot.sendText(chatID, fmt.Sprintf(template.Balance, balance), mainInlineKeyboardCheckBalance); err != nil {
		bot.err(err, "send wallet balance")
	}
}

func (bot *Bot) msgUpdateBalance(update tgBotAPI.Update, wallet string) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	balance, err := bot.ton.GetWalletBalance(wallet)
	if err != nil {
		bot.err(err, "get wallet balance")
		return
	}

	newText := fmt.Sprintf(template.Balance, balance)

	notification := "Balance is up to date"

	if newText != update.CallbackQuery.Message.Text {
		msg := tgBotAPI.EditMessageTextConfig{
			BaseEdit: tgBotAPI.BaseEdit{
				ChatID:          chatID,
				ChannelUsername: "",
				MessageID:       messageID,
				InlineMessageID: "",
				ReplyMarkup:     &mainInlineKeyboardCheckBalance,
			},
			Text:                  newText,
			ParseMode:             "HTML",
			Entities:              nil,
			DisableWebPagePreview: false,
		}

		_, err = bot.api.Request(msg)
		if err != nil {
			log.Error().Err(err).Send()
			bot.err(err, "")
			return
		}

		notification = "Balance updated"
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, notification)
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error().Err(err).Send()
		bot.err(err, "")
	}

}

func (bot *Bot) msgCancel(ctx context.Context, chatID int64, userID string) {
	err := bot.redis.SetStage(ctx, userID, storage.StageUnset)
	if err != nil {
		bot.err(err, "set \"unset\" stage")
		return
	}

	if err = bot.sendText(chatID, template.Canceled, mainKeyboard); err != nil {
		bot.err(err, "send canceled message")
	}
}

func (bot *Bot) cancelInline(ctx context.Context, update tgBotAPI.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	text := update.CallbackQuery.Message.Text
	userID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)

	err := bot.redis.SetStage(ctx, userID, storage.StageUnset)
	if err != nil {
		bot.err(err, "set \"unset\" stage")
		return
	}

	newText := text + "\n<b>✅ Canceled</b>"

	msg := tgBotAPI.EditMessageTextConfig{
		BaseEdit: tgBotAPI.BaseEdit{
			ChatID:          chatID,
			ChannelUsername: "",
			MessageID:       messageID,
			InlineMessageID: "",
			ReplyMarkup:     nil,
		},
		Text:                  newText,
		ParseMode:             "HTML",
		Entities:              nil,
		DisableWebPagePreview: false,
	}

	_, err = bot.api.Request(msg)
	if err != nil {
		log.Error().Err(err).Send()
		bot.err(err, "")
	}
}

func (bot *Bot) cmdStart(update tgBotAPI.Update, isExist bool, wallet string) {
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

	if err := bot.sendPhoto(update.Message.Chat.ID, qr, caption, mainInlineKeyboard); err != nil {
		bot.err(err, "send wallet address and qr")
	}
}

func (bot *Bot) msgAmountCheck(ctx context.Context, chatID int64, userID, wallet, text string, messageID int, rw string) {
	amount := strings.ReplaceAll(text, " ", "")
	amount = strings.ReplaceAll(amount, ",", ".")

	amountInCoins, err := tlb.FromTON(amount)
	if err != nil {
		if err := bot.sendText(chatID, template.InvalidAmount, struct{}{}); err != nil {
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
		if err := bot.sendText(chatID, txt, struct{}{}); err != nil {
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

	txt := fmt.Sprintf(template.SendingConfirmation, rw)
	if err := bot.sendText(chatID, txt, confirmInlineKeyboard); err != nil {
		bot.err(err, "send validate fail")
		return
	}

	//edit := tgBotAPI.EditMessageReplyMarkupConfig{
	//	BaseEdit: tgBotAPI.BaseEdit{
	//		ChatID:          chatID,
	//		ChannelUsername: "",
	//		MessageID:       messageID - 1,
	//		InlineMessageID: "",
	//		ReplyMarkup:     nil,
	//	},
	//}
	//
	//_, err = bot.api.Request(edit)
	//if err != nil {
	//	log.Error().Err(err).Send()
	//}

}

func (bot *Bot) msgAddressCheck(ctx context.Context, update tgBotAPI.Update, chatID int64) {
	userID := strconv.FormatInt(update.Message.Chat.ID, 10)

	//
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
		if err := bot.sendText(chatID, template.InvalidQR, struct{}{}); err != nil {
			bot.err(err, "send decode fail")
			return
		}
		return
	}

	receiverAddress := result.String()

	err = bot.ton.ValidateWallet(receiverAddress)
	if err != nil {
		bot.err(err, "validate receiver wallet")
		if err := bot.sendText(chatID, template.InvalidWallet, struct{}{}); err != nil {
			bot.err(err, "send validate fail")
			return
		}
		return
	}

	err = bot.redis.SetStage(ctx, userID, storage.StageAmountWaiting)
	if err != nil {
		bot.err(err, "set wallet waiting stage")
		return
	}

	if err := bot.sendText(chatID, template.AskAmount, struct{}{}); err != nil {
		bot.err(err, "ask wallet")
		return
	}

}
