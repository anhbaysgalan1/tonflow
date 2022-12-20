package bot

import (
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type logger struct {
}

func (l logger) Printf(_ string, v ...interface{}) {
	log.Debug().Msgf("%s", v)
}
func (l logger) Println(v ...interface{}) {
	l.Printf("%s", v)
}

func (bot *Bot) err(err error, desc string) {
	text := ""

	switch desc {
	case "":
		log.Error().Err(err).Send()
		text = fmt.Sprintf("%s", err.Error())
	default:
		log.Error().Err(err).Msg(desc)
		text = fmt.Sprintf("%s: %s ", desc, err.Error())
	}

	msg := tgBotAPI.NewMessage(-1001638881880, text)
	msg.ParseMode = "HTML"
	msg.DisableNotification = false

	_, er := bot.api.Send(msg)
	if er != nil {
		log.Error().Err(er).Send()
	}
}
