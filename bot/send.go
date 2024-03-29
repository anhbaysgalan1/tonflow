package bot

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) sendText(chatID int64, text string, markup interface{}) error {
	msg := tgBotAPI.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = true
	msg.ReplyMarkup = markup

	_, err := bot.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) sendNotification(chatID int64, text string, markup interface{}) error {
	msg := tgBotAPI.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = false
	msg.ReplyMarkup = markup

	_, err := bot.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) sendImage(chatID int64, data []byte, caption string, markup interface{}) error {
	photoBytes := tgBotAPI.FileBytes{
		Bytes: data,
	}

	msg := tgBotAPI.NewPhoto(chatID, photoBytes)
	msg.ParseMode = "HTML"
	msg.Caption = caption
	msg.DisableNotification = true
	msg.ReplyMarkup = markup

	_, err := bot.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) deleteMessage(chatID int64, messageID int) {
	deleteConfig := tgBotAPI.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}

	_, err := bot.api.Request(deleteConfig)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) sendTyping(chatID int64) {
	action := tgBotAPI.ChatActionConfig{
		BaseChat: tgBotAPI.BaseChat{ChatID: chatID},
		Action:   tgBotAPI.ChatTyping,
	}

	_, err := bot.api.Request(action)
	if err != nil {
		log.Error(err)
	}
}
