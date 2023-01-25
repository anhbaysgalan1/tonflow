package bot

import (
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type logger struct {
}

func (l logger) Println(v ...interface{}) {
	log.Debugln(v)
}

func (l logger) Printf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func (bot *Bot) sendErr(err error, desc string) {
	text := ""

	switch desc {
	case "":
		text = fmt.Sprintf("%v", err.Error())
	default:
		text = fmt.Sprintf("%v: %v ", desc, err.Error())
	}

	msg := tgBotAPI.NewMessage(-1001638881880, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = false

	_, er := bot.api.Send(msg)
	if er != nil {
		log.Error(err)
	}
}
