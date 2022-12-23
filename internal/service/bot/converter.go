package bot

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tonflow/internal/service/bot/model"
)

func toFlowUser(u *tgBotAPI.User) *model.User {
	return &model.User{
		ID:           u.ID,
		Username:     u.UserName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		LanguageCode: u.LanguageCode,
	}
}
