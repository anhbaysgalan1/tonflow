package bot

import (
	"context"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"tonflow/internal/storage"
)

func (bot *Bot) handleUpdate(ctx context.Context, update tgBotAPI.Update) {
	switch {
	case update.CallbackQuery != nil:
		bot.handleCallbackQuery(ctx, update)
	case update.Message != nil:
		bot.handleMessage(ctx, update)
	}
}

func (bot *Bot) handleCallbackQuery(ctx context.Context, update tgBotAPI.Update) {
	cbData := update.CallbackData()
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)

	_, wallet, _, err := bot.checkUser(ctx, update)
	if err != nil {
		bot.err(err, "check user")
		return
	}

	switch cbData {
	case "receive":
		bot.msgReceivingOptions(chatID, wallet)
	case "send":
		//msg := update.CallbackQuery.Message
		//
		//edit := tgBotAPI.EditMessageReplyMarkupConfig{
		//	BaseEdit: tgBotAPI.BaseEdit{
		//		ChatID:          chatID,
		//		ChannelUsername: "",
		//		MessageID:       msg.MessageID,
		//		InlineMessageID: "",
		//		ReplyMarkup:     nil,
		//	},
		//}
		//
		//_, err := bot.api.Request(edit)
		//if err != nil {
		//	log.Error().Err(err).Send()
		//}
		bot.msgAskAddress(ctx, chatID, userID, wallet)
	case "balance":
		bot.msgBalance(chatID, wallet)
	case "update balance":
		bot.msgUpdateBalance(update, wallet)
	case "cancel":
		bot.cancelInline(ctx, update)
		bot.msgBalance(chatID, wallet)
	}

	// SEND:
	// отправляем сообщение

	//msg := update.CallbackQuery.Message
	//
	//newText := "Some text"
	//
	//edit := tgBotAPI.EditMessageCaptionConfig{
	//	BaseEdit: tgBotAPI.BaseEdit{
	//		ChatID:      msg.Chat.ID,
	//		MessageID:   msg.MessageID,
	//		ReplyMarkup: &cancelInlineKeyboard,
	//	},
	//	Caption: newText,
	//}
	//
	//_, err := bot.api.Request(edit)
	//if err != nil {
	//	log.Error().Err(err).Send()
	//}
}

func (bot *Bot) handleMessage(ctx context.Context, update tgBotAPI.Update) {
	isExist, wallet, stage, err := bot.checkUser(ctx, update)
	if err != nil {
		bot.err(err, "check user")
		return
	}

	switch {
	case update.Message.IsCommand(): // command handling
		switch update.Message.Command() {
		case "start":
			bot.cmdStart(update, isExist, wallet)
		}
	default: // non-command user messages handling
		chatID := update.Message.Chat.ID
		messageID := update.Message.MessageID
		userID := strconv.FormatInt(toFlowUser(update.SentFrom()).ID, 10)
		text := update.Message.Text

		if stage == storage.StageWalletWaiting && len(update.Message.Photo) != 0 {
			bot.msgAddressCheck(ctx, update, chatID)
		}

		if stage == storage.StageAmountWaiting && update.Message.Text != "" {
			bot.msgAmountCheck(ctx, chatID, userID, wallet, text, messageID, "xxx-xxx-xxx-xxx")
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
