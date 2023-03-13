package storage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/context"
	"tonflow/model"
)

type Storage interface {
	AddUser(ctx context.Context, user *tgbotapi.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)

	AddWallet(ctx context.Context, wallet *model.Wallet, userID int64) error
	GetWallet(ctx context.Context, address string) (*model.Wallet, error)

	GetInMemoryWallets() map[string]int64

	SetLastSeqno(ctx context.Context, shards map[string]uint32) error
	GetLastSeqno(ctx context.Context) (map[string]uint32, error)
}

type Cache interface {
	SetUserCache(ctx context.Context, user *model.User) error
	GetUserCache(ctx context.Context, userID int64) (*model.User, error)
}
