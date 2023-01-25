package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
	"tonflow/bot/model"
	"tonflow/storage"
	"tonflow/tonclient"
)

type DB struct {
	*redis.Client
}

type Config struct {
	Host string
	Port string
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

func (db *DB) SetUserCache(ctx context.Context, cache *model.UserCache) error {
	id := strconv.FormatInt(cache.UserID, 10)

	// log.Debugf("SetUserCache() cache.Data to marshal:\n%v", pkg.AnyPrint(cache.Data))

	data, err := json.Marshal(cache.Data)
	if err != nil {
		return err
	}

	// log.Debugf("SetUserCache() marshaled:\n%v", pkg.AnyPrint(string(data)))

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
	user := model.User{Wallet: &tonclient.Wallet{}}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
