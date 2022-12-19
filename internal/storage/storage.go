package storage

import (
	"golang.org/x/net/context"
	"park-wallet/internal/service/bot/model"
	"park-wallet/internal/service/ton"
	"time"
)

type Storage interface {
	CheckUser(ctx context.Context, user *model.User) (bool, error)

	AddWallet(ctx context.Context, wallet *ton.Wallet, userID int64) error
	GetUserWallet(ctx context.Context, userID int64) (string, error)

	AddTransaction(ctx context.Context, tr ton.Transaction) error
	GetUserTransactions(ctx context.Context, userID int64) ([]*ton.Transaction, error)

	AddPicture(ctx context.Context, ID string, time time.Time) error
	GetRandomPicture(ctx context.Context) (string, error)
	GetAllPictures(ctx context.Context) ([]string, error)
}
