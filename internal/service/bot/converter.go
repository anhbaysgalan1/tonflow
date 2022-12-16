package bot

import (
	telegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ton-flow-bot/internal/service/bot/model"
)

func toFlowUser(u *telegramBotAPI.User) *model.User {
	return &model.User{
		ID:           u.ID,
		Username:     u.UserName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		LanguageCode: u.LanguageCode,
	}
}
