package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
	"ton-flow-bot/internal/service/model"
)

func toUser(u *tgbotapi.User) model.User {
	return model.User{
		ID:             u.ID,
		Username:       u.UserName,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		LanguageCode:   u.LanguageCode,
		FirstMessageAt: time.Now(),
		LastMessageAt:  time.Now(),
		Wallet:         "",
	}
}
