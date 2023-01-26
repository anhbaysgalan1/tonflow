package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
	"tonflow/model"
	"tonflow/storage"
)

type DB struct {
	*redis.Client
}

func NewRedisClient(URI string) (storage.Cache, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     URI,
			Password: "",
			DB:       0,
		},
	)

	// verify redis connection
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &DB{redisClient}, nil
}

func (db *DB) SetUserCache(ctx context.Context, user *model.User) error {
	id := strconv.FormatInt(user.ID, 10)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = db.Set(ctx, id, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserCache(ctx context.Context, userID int64) (*model.User, error) {
	id := strconv.FormatInt(userID, 10)

	val, err := db.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	b := []byte(val)
	user := &model.User{Wallet: &model.Wallet{}}
	err = json.Unmarshal(b, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
