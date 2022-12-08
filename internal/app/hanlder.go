package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"time"
	"ton-flow-bot/pkg"
)

func (app *App) handleUpdate(update tgbotapi.Update) {
	log.Debug().Msg(pkg.AnyPrint("UPDATE", update))

	ctx := context.Background()

	if update.Message != nil { // If we got a message

		user := toUser(update.SentFrom())
		_, err := app.storage.CheckUser(ctx, user)
		if err != nil {
			log.Error().Err(err).Send()
		}

		switch {
		case update.Message.From.ID == 903169: // if msg received from admin

			msg := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}
			if _, err := app.bot.Send(msg); err != nil {
				log.Error().Err(err)
			}

			if len(update.Message.Photo) != 0 {
				ID := update.Message.Photo[0].FileID
				err = app.storage.AddPicture(ctx, ID, time.Now())
				if err != nil {
					log.Error().Err(err).Send()
				}
			}

			if update.Message.Text == "778" {
				id, err := app.storage.GetRandomPicture(ctx)
				if err != nil {
					log.Error().Err(err).Send()
				}
				photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(id))
				_, err = app.bot.Send(photo)
				if err != nil {
					log.Error().Err(err)
				}

				time.AfterFunc(time.Second*30, func() {
					msg = tgbotapi.DeleteMessageConfig{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID + 1,
					}
					if _, err = app.bot.Send(msg); err != nil {
						log.Error().Err(err)
					}
				})
			}

		default: // from other users

		}
	}
}
