package storage

import (
	"golang.org/x/net/context"
	"time"
	"ton-flow-bot/internal/service/bot/model"
	"ton-flow-bot/internal/service/ton"
)

type Storage interface {
	CheckUser(ctx context.Context, user *model.User) (bool, error)

	AddWallet(ctx context.Context, wallet *ton.Wallet, userID int64) error
	GetUserWallet(ctx context.Context, userID int64) (string, error)

	AddTransaction(ctx context.Context, tr ton.Transaction) error
	GetUserTransactions(ctx context.Context, userID int64) ([]*ton.Transaction, error)

	AddPicture(ctx context.Context, ID string, time time.Time) error
	GetRandomPicture(ctx context.Context) (string, error)
}
