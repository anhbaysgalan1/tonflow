package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"park-wallet/internal/storage"
	"time"
)

type DB struct {
	*redis.Client
}

type Config struct {
	Host string
	Port string
}

func strToStage(s string) (storage.Stage, error) {
	log.Debug().Msg(s)
	switch s {
	case "unset":
		return storage.StageUnset, nil
	case "walletWaiting":
		return storage.StageWalletWaiting, nil
	case "amountWaiting":
		return storage.StageAmountWaiting, nil
	default:
		return "", errors.New("no matching string and stage options")
	}
}

func NewRedisClient(cfg *Config) (storage.TemporaryStorage, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: "",
			DB:       0,
		},
	)

	// verify redis connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &DB{redisClient}, nil
}

func (db *DB) GetStage(ctx context.Context, userID string) (storage.Stage, error) {
	stringCmd := db.Get(ctx, userID)

	result, err := stringCmd.Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if err == redis.Nil {
		return storage.StageUnset, nil
	}

	out, err := strToStage(result)
	if err != nil {
		return "", err
	}
	return out, nil
}
func (db *DB) SetStage(ctx context.Context, userID string, stage storage.Stage) error {
	statusCmd := db.Set(ctx, userID, stage.String(), time.Hour*720)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}

	return nil
}
