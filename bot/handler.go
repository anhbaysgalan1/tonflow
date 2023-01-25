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
	"tonflow/bot/model"
	"tonflow/bot/template"
	"tonflow/pkg"
	"tonflow/tonclient"
)

// Returns user object from cache or database in case of no cache.
// Also returns "true" if user is existed already else creates and returns new user with "false" second value.
func (bot *Bot) getTonflowUser(ctx context.Context, tgUser *tgBotAPI.User) (*model.User, bool, error) {
	user, err := bot.redis.GetUserCache(ctx, tgUser.ID)
	if err != nil && err != redis.Nil {
		log.Error(err)
		return nil, false, err
	}

	wallet := &tonclient.Wallet{}
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

			wlt.Seed, err = pkg.EncryptAES([]byte(bot.cryptoKey), wlt.Seed)
			if err != nil {
				log.Error(err)
				return nil, false, err
			}

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

		userCache := &model.UserCache{
			UserID: tgUser.ID,
			Data:   user,
		}

		err = bot.redis.SetUserCache(ctx, userCache)
		if err != nil {
			log.Error(err)
			return nil, true, err
		}
	}

	return user, isExisted, nil
}

func (bot *Bot) inlineReceiveCoins(user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	qr, err := qrcode.Encode(address, qrcode.Medium, 512)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "generate qr")
		return
	}

	caption := fmt.Sprintf("<code>%s</code>", address) + "\n\n" + template.ReceiveInstruction
	if err := bot.sendPhoto(chatID, qr, caption, mainInlineKeyboard); err != nil {
		log.Error(err)
		bot.sendErr(err, "send wallet address and qr")
	}
}

func (bot *Bot) inlineSendCoins(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	balance, err := bot.ton.GetWalletBalance(address)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "get balance")
		return
	}

	if balance == "0" {
		if err = bot.sendText(chatID, template.NoFunds, nil); err != nil {
			log.Error(err)
			bot.sendErr(err, "not enough message")
			return
		}
		return
	}

	user.StageData.Stage = model.AddressWait

	userCache := &model.UserCache{
		UserID: user.ID,
		Data:   user,
	}

	if err = bot.redis.SetUserCache(ctx, userCache); err != nil {
		log.Error(err)
		bot.sendErr(err, "set wallet waiting stage")
		return
	}

	if err = bot.sendText(chatID, template.AskWallet, nil); err != nil {
		log.Error(err)
		bot.sendErr(err, "ask amount to send")
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "")
	}
}

func (bot *Bot) inlineBalance(update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address
	chatID := user.ID

	balance, err := bot.ton.GetWalletBalance(address)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "get wallet balance")
		return
	}

	if err = bot.sendText(chatID, fmt.Sprintf(template.Balance, balance), mainInlineKeyboardCheckBalance); err != nil {
		log.Error(err)
		bot.sendErr(err, "send wallet balance")
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "")
	}
}

func (bot *Bot) inlineUpdateBalance(update tgBotAPI.Update, user *model.User) {
	balance, err := bot.ton.GetWalletBalance(user.Wallet.Address)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "get wallet balance")
		return
	}

	newText := fmt.Sprintf(template.Balance, balance)

	notification := "Balance is up to date"

	if newText != update.CallbackQuery.Message.Text {
		msg := tgBotAPI.EditMessageTextConfig{
			BaseEdit: tgBotAPI.BaseEdit{
				ChatID:          user.ID,
				ChannelUsername: "",
				MessageID:       update.CallbackQuery.Message.MessageID,
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
			log.Error(err)
			bot.sendErr(err, "")
			return
		}

		notification = "Balance updated"
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, notification)
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "")
	}

}

func (bot *Bot) inlineCancel(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	text := update.CallbackQuery.Message.Text
	//userID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)

	//err := bot.redis.SetStage(ctx, userID, storage.StageUnset)
	//if err != nil {
	//	bot.err(err, "set \"unset\" stage")
	//	return
	//}

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

	user.StageData.Stage = model.ZeroStage

	userCache := &model.UserCache{
		UserID: user.ID,
		Data:   user,
	}

	err := bot.redis.SetUserCache(ctx, userCache)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = bot.api.Request(msg)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "")
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "")
	}
}

func (bot *Bot) cmdStart(update tgBotAPI.Update, user *model.User, isExist bool) {
	message := update.Message
	wallet := user.Wallet.Address

	qr, err := qrcode.Encode(wallet, qrcode.Medium, 512)
	if err != nil {
		log.Error(err)
		bot.sendErr(err, "generate qr")
		return
	}

	caption := ""
	switch isExist {
	case true:
		caption = fmt.Sprintf(template.WelcomeExistedUser, message.From.FirstName, wallet)
	case false:
		caption = fmt.Sprintf(template.WelcomeNewUser, message.From.FirstName, wallet)
	}

	if err := bot.sendPhoto(message.Chat.ID, qr, caption, mainInlineKeyboard); err != nil {
		bot.sendErr(err, "send wallet address and qr")
	}
}

func (bot *Bot) validateAmount(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	wallet := user.Wallet.Address

	amount := strings.ReplaceAll(text, " ", "")
	amount = strings.ReplaceAll(amount, ",", ".")
	rw := "xxx-xxx-xxx-xxx"

	amountInCoins, err := tlb.FromTON(amount)
	if err != nil {
		if err := bot.sendText(chatID, template.InvalidAmount, struct{}{}); err != nil {
			bot.sendErr(err, "send invalid amount")
			return
		}
		bot.sendErr(err, "convert string to coins")
		return
	}

	balance, err := bot.ton.GetWalletBalance(wallet)
	if err != nil {
		bot.sendErr(err, "get wallet balance")
		return
	}

	balanceInCoins, err := tlb.FromTON(balance)
	if err != nil {
		bot.sendErr(err, "convert string to coins")
		return
	}

	if balanceInCoins.NanoTON().Cmp(amountInCoins.NanoTON()) < 0 {
		txt := fmt.Sprintf(template.NotEnoughFunds, balance)
		if err := bot.sendText(chatID, txt, struct{}{}); err != nil {
			bot.sendErr(err, "send not enough amount")
			return
		}
		return
	}

	//err = bot.redis.SetStage(ctx, userID, storage.StageWalletWaiting)
	//if err != nil {
	//	bot.err(err, "set wallet waiting stage")
	//	return
	//}

	txt := fmt.Sprintf(template.SendingConfirmation, rw)
	if err := bot.sendText(chatID, txt, confirmInlineKeyboard); err != nil {
		bot.sendErr(err, "send validate fail")
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

func (bot *Bot) parseQR(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	message := update.Message

	index, size := 0, 0
	for i, v := range message.Photo {
		if v.FileSize > size {
			index = i
		}
	}

	fileURL, err := bot.api.GetFileDirectURL(message.Photo[index].FileID)
	if err != nil {
		bot.sendErr(err, "get file direct url")
		return
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		bot.sendErr(err, "get data from url")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bot.sendErr(errors.New("bad status: "+resp.Status), "http status")
		return
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		bot.sendErr(err, "decode jpeg")
		return
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		bot.sendErr(err, "prepare binary bitmap")
		return
	}

	// decode image
	qrReader := qrScan.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		bot.sendErr(err, "decode qr")
		if err := bot.sendText(message.Chat.ID, template.InvalidQR, struct{}{}); err != nil {
			bot.sendErr(err, "send decode fail")
			return
		}
		return
	}

	receiverAddress := result.String()

	err = bot.ton.ValidateWallet(receiverAddress)
	if err != nil {
		bot.sendErr(err, "validate receiver wallet")
		if err := bot.sendText(message.Chat.ID, template.InvalidWallet, struct{}{}); err != nil {
			bot.sendErr(err, "send validate fail")
			return
		}
		return
	}

	//userID := strconv.FormatInt(mc.FromID, 10)
	//err = bot.redis.SetStage(ctx, userID, storage.StageAmountWaiting)
	//if err != nil {
	//	bot.err(err, "set wallet waiting stage")
	//	return
	//}

	user.StageData.Stage = model.AmountWait

	userCache := &model.UserCache{
		UserID: user.ID,
		Data:   user,
	}

	err = bot.redis.SetUserCache(ctx, userCache)
	if err != nil {
		log.Error(err)
		return
	}

	if err := bot.sendText(message.Chat.ID, template.AskAmount, struct{}{}); err != nil {
		bot.sendErr(err, "ask wallet")
		log.Error(err)
		return
	}

}
