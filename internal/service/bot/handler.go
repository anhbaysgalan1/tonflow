package bot

import (
	"context"
	"fmt"
	telegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

func (bot *Bot) handleUpdate(ctx context.Context, update telegramBotAPI.Update) {
	switch {
	case update.Message == nil:
		bot.handleNilMessage(ctx, update)
	case update.Message != nil:
		bot.handleMessage(ctx, update)
	}
}

func (bot *Bot) handleNilMessage(_ context.Context, update telegramBotAPI.Update) {
	if update.CallbackQuery != nil {
		callback := telegramBotAPI.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		_, err := bot.api.Request(callback)
		if err != nil {
			msg := "send callback message"
			log.Error().Err(err).Msg(msg)
			bot.err(err, msg)
		}
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update telegramBotAPI.Update) {
	isExist, wallet, err := bot.checkUser(ctx, update)
	if err != nil {
		msg := "check user"
		log.Error().Err(err).Msg(msg)
		bot.err(err, msg)
		return
	}

	switch {
	case update.Message.IsCommand():
		bot.handleCommand(ctx, update, isExist, wallet)
	case update.Message.From.ID == bot.adminID:
		// bot.handleAdminMessage(ctx, update)
		//default:
		bot.handleUserMessage(ctx, update, isExist, wallet)
	}
}

func (bot *Bot) checkUser(ctx context.Context, update telegramBotAPI.Update) (bool, string, error) {
	bot.sendTyping(update.Message.From.ID)

	flowUser := toFlowUser(update.SentFrom())

	userExist, err := bot.storage.CheckUser(ctx, flowUser)
	if err != nil {
		return false, "", err
	}

	wlt, err := bot.storage.GetUserWallet(ctx, flowUser.ID)
	if err != nil {
		return false, "", err
	}

	if !userExist || (wlt == "" && err == nil) {
		wallet, err := bot.ton.NewWallet()
		if err != nil {
			return false, "", err
		}

		err = bot.storage.AddWallet(ctx, wallet, flowUser.ID)
		if err != nil {
			return false, "", err
		}

		wlt = wallet.Address

		return false, wlt, nil
	}

	return true, wlt, nil
}

func (bot *Bot) handleCommand(ctx context.Context, update telegramBotAPI.Update, isExist bool, wallet string) {
	switch update.Message.Command() {
	case "start":
		text := ""
		if !isExist {
			text = fmt.Sprintf(startNewUser, update.Message.From.FirstName, wallet)
		} else {
			text = fmt.Sprintf(startRegisteredUser, update.Message.From.FirstName, wallet)
		}

		qr, err := qrcode.Encode(wallet, qrcode.Medium, 512)
		if err != nil {
			msg := "generate QR"
			log.Error().Err(err).Msg(msg)
			bot.err(err, msg)
		}

		qrBytes := telegramBotAPI.FileBytes{
			Bytes: qr,
		}

		startMsg := telegramBotAPI.NewPhoto(update.Message.Chat.ID, qrBytes)
		startMsg.ParseMode = "HTML"
		startMsg.Caption = text
		startMsg.DisableNotification = true

		buttons := mainKeyboard
		buttons.ResizeKeyboard = true
		buttons.InputFieldPlaceholder = ""

		startMsg.ReplyMarkup = buttons

		_, err = bot.api.Send(startMsg)
		if err != nil {
			msg := "send start message"
			log.Error().Err(err).Msg(msg)
			bot.err(err, msg)
		}
	}
}

func (bot *Bot) handleUserMessage(ctx context.Context, update telegramBotAPI.Update, isExist bool, wallet string) {
	switch update.Message.Text {
	case "💎 Balance":
		bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)

		balance, err := bot.ton.GetWalletBalance(wallet)
		if err != nil {
			msg := "get wallet balance"
			log.Error().Err(err).Msg(msg)
			bot.err(err, msg)
		}

		text := fmt.Sprintf("💎 Your balance is %s TON", balance)
		textMsg := telegramBotAPI.NewMessage(update.Message.Chat.ID, text)
		textMsg.DisableNotification = true

		_, err = bot.api.Send(textMsg)
		if err != nil {
			msg := "send balance message"
			log.Error().Err(err).Msg(msg)
			bot.err(err, msg)
		}
	}
}

//func (bot *Bot) handleAdminMessage(ctx context.Context, update telegramBotAPI.Update) {
//	bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)
//
//	switch {
//	case len(update.Message.Photo) != 0:
//		ID := update.Message.Photo[0].FileID
//		err := bot.storage.AddPicture(ctx, ID, time.Now())
//		if err != nil {
//			log.Error().Err(err).Send()
//			bot.err(err, "failed to add picture in storage")
//			bot.sendText(update.Message.Chat.ID, "One of the pictures was not saved in the database", telegramBotAPI.ReplyKeyboardMarkup{})
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

//func (bot *Bot) sendText(chatID int64, text string, markup telegramBotAPI.ReplyKeyboardMarkup) {
//	msg := telegramBotAPI.NewMessage(chatID, text)
//	msg.ParseMode = "HTML"
//	msg.DisableNotification = true
//
//	_, err := bot.api.Send(msg)
//	if err != nil {
//		log.Error().Err(err).Send()
//		bot.err(err, "failed to send text message")
//	}
//}

func (bot *Bot) sendTyping(chatID int64) {
	action := telegramBotAPI.ChatActionConfig{
		BaseChat: telegramBotAPI.BaseChat{ChatID: chatID},
		Action:   telegramBotAPI.ChatTyping,
	}
	_, err := bot.api.Request(action)
	if err != nil {
		msg := "send typing action"
		log.Error().Err(err).Msg(msg)
		bot.err(err, msg)
	}
}

//func (bot *Bot) sendUploadingPhoto(chatID int64) {
//	action := telegramBotAPI.ChatActionConfig{
//		BaseChat: telegramBotAPI.BaseChat{ChatID: chatID},
//		Action:   telegramBotAPI.ChatUploadPhoto,
//	}
//	_, err := bot.api.Request(action)
//	if err != nil {
//		log.Error().Err(err).Send()
//		bot.err(err, "failed to send uploading photo action")
//	}
//}

func (bot *Bot) deleteMessage(chatID int64, messageID int) {
	deleteConfig := telegramBotAPI.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}
	_, err := bot.api.Request(deleteConfig)
	if err != nil {
		msg := "auto delete message"
		log.Error().Err(err).Msg(msg)
		bot.err(err, msg)
	}
}
