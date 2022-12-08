package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"ton-flow-bot/internal/service/model"
	"ton-flow-bot/internal/storage"
)

type DB struct {
	*pgxpool.Pool
}

type Config struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Migration bool
}

func NewPGStorage(cfg *Config) (storage.Storage, error) {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=require",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	//conn, err := sqlx.Connect("postgres", url)
	//if err != nil {
	//	return nil, err
	//}

	if cfg.Migration {
		err = Migration(url)
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		pool,
	}, nil
}

func (db *DB) Close() {
	db.Close()
}

func (db *DB) CheckUser(ctx context.Context, u model.User) (bool, error) {
	query := `select id from users where id = $1`
	var userID int64
	err := db.QueryRow(ctx, query, u.ID).Scan(&userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		query := `insert into users (id, username, first_name, last_name, language_code, wallet, first_message_at, last_message_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := db.Exec(ctx, query, u.ID, u.Username, u.FirstName, u.LastName, u.LanguageCode, u.Wallet, u.FirstMessageAt, u.LastMessageAt)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	query = `update users set last_message_at = $2 where id = $1`
	_, err = db.Exec(ctx, query, u.ID, u.LastMessageAt)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) AddPicture(ctx context.Context, ID string, time time.Time) error {
	query := `insert into pictures (id, added_at) values ($1, $2)`
	_, err := db.Exec(ctx, query, ID, time)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRandomPicture(ctx context.Context) (string, error) {
	query := `select id from pictures order by random() limit 1`

	id := ""
	err := db.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return "", nil
	}

	return id, nil
}
