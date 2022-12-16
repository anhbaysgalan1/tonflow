package bot

import (
	"context"
	"errors"
	telegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
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
			bot.err(err, "")
		}
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update telegramBotAPI.Update) {
	_, err := bot.checkUser(ctx, update) //TODO: сделать надежную логику внутри этого метода
	if err != nil {
		bot.err(err, "failed to check user")
		return
	}

	switch {
	case update.Message.IsCommand():
		// bot.handleCommand(ctx, update)
	case update.Message.From.ID == bot.adminID:
		bot.handleAdminMessage(ctx, update)
	default:
		// bot.handleUserMessage(ctx, update)
	}
}

func (bot *Bot) checkUser(ctx context.Context, update telegramBotAPI.Update) (bool, error) {
	flowUser := toFlowUser(update.SentFrom())

	userExist, err := bot.storage.CheckUser(ctx, flowUser)
	if err != nil {
		return false, err
	}

	wlt, err := bot.storage.GetUserWallet(ctx, flowUser.ID)
	if err != nil {
		return false, err
	}

	if !userExist || (wlt == "" && err == nil) {
		wallet, err := bot.ton.NewWallet()
		if err != nil {
			return false, err
		}

		err = bot.storage.AddWallet(ctx, wallet, flowUser.ID)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (bot *Bot) handleAdminMessage(ctx context.Context, update telegramBotAPI.Update) {
	bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID)

	switch {
	case len(update.Message.Photo) != 0:
		ID := update.Message.Photo[0].FileID
		err := bot.storage.AddPicture(ctx, ID, time.Now())
		if err != nil {
			bot.err(err, "failed to add picture in storage")
			bot.sendText(update.Message.Chat.ID, "One of the pictures was not saved in the database", telegramBotAPI.ReplyKeyboardMarkup{})
		}
	case update.Message.Text == "778":
		bot.sendUploadingPhoto(update.Message.Chat.ID)

		fileID, err := bot.storage.GetRandomPicture(ctx)
		if err != nil {
			bot.err(err, "failed to get random pic")
			return
		}
		if fileID == "" {
			bot.err(errors.New("no pics in database"), "")
			return
		}

		bot.sendPhoto(update.Message.Chat.ID, fileID)

		time.AfterFunc(time.Second*5, func() {
			bot.deleteMessage(update.Message.Chat.ID, update.Message.MessageID+1)
		})
	default:

	}
}

//func (bot *Bot) cmdStart(ctx context.Context) {
//	user := ctx.Value("user")
//	chatID := ctx.Value("chatID")
//
//	wallet, err := bot.storage.GetUserWallet(ctx, user)
//	if err != nil {
//		log.Error().Err(err).Send()
//		break
//	}
//
//	var text string
//	if !isRegistered {
//		text = fmt.Sprintf(startNewUser, update.Message.From.FirstName)
//	} else {
//		text = fmt.Sprintf(startRegisteredUser, update.Message.From.FirstName)
//	}
//	bot.sendText(chatID, text, api.ReplyKeyboardMarkup{})
//	text = "<pre>" + wallet + "</pre>"
//	bot.sendText(chatID, text, mainKeyboard)
//}

func (bot *Bot) sendText(chatID int64, text string, markup telegramBotAPI.ReplyKeyboardMarkup) {
	msg := telegramBotAPI.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = true
	msg.ReplyMarkup = markup

	_, err := bot.api.Send(msg)
	if err != nil {
		bot.err(err, "failed to send text message")
	}
}

func (bot *Bot) sendPhoto(chatID int64, fileID string) {
	photo := telegramBotAPI.NewPhoto(chatID, telegramBotAPI.FileID(fileID))
	_, err := bot.api.Send(photo)
	if err != nil {
		bot.err(err, "failed to send photo")
	}
}

func (bot *Bot) sendTyping(chatID int64) {
	action := telegramBotAPI.ChatActionConfig{
		BaseChat: telegramBotAPI.BaseChat{ChatID: chatID},
		Action:   telegramBotAPI.ChatTyping,
	}
	_, err := bot.api.Request(action)
	if err != nil {
		bot.err(err, "failed to send typing action")
	}
}

func (bot *Bot) sendUploadingPhoto(chatID int64) {
	action := telegramBotAPI.ChatActionConfig{
		BaseChat: telegramBotAPI.BaseChat{ChatID: chatID},
		Action:   telegramBotAPI.ChatUploadPhoto,
	}
	_, err := bot.api.Request(action)
	if err != nil {
		bot.err(err, "failed to send uploading photo action")
	}
}

func (bot *Bot) deleteMessage(chatID int64, messageID int) {
	deleteConfig := telegramBotAPI.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}
	_, err := bot.api.Request(deleteConfig)
	if err != nil {
		bot.err(err, "failed to delete message")
	}
}
