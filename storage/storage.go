package storage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/context"
	"tonflow/bot/model"
	"tonflow/tonclient"
)

type Storage interface {
	AddUser(ctx context.Context, user *tgbotapi.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)

	AddWallet(ctx context.Context, wallet *tonclient.Wallet, userID int64) error
	GetWallet(ctx context.Context, address string) (*tonclient.Wallet, error)
}

type Cache interface {
	SetUserCache(ctx context.Context, cache *model.UserCache) error
	GetUserCache(ctx context.Context, userID int64) (*model.User, error)
}
