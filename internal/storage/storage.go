package storage

import (
	"golang.org/x/net/context"
	"time"
	"ton-flow-bot/internal/service/model"
)

type Storage interface {
	CheckUser(ctx context.Context, user model.User) (bool, error)
	AddPicture(ctx context.Context, ID string, time time.Time) error
	GetRandomPicture(ctx context.Context) (string, error)
	Close()
}
