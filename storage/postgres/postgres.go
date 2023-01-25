package postgres

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
	"tonflow/bot/model"
	"tonflow/storage"
	"tonflow/tonclient"
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
	SSL       string
	Migration bool
}

func NewConnection(URI string) (storage.Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool, err := pgxpool.New(ctx, URI)
	if err != nil {
		return nil, err
	}

	return &DB{
		pool,
	}, nil
}

func (db *DB) AddUser(ctx context.Context, user *tgbotapi.User) error {
	query := `
		insert into tonflow.users (
		id,
		username,
		first_name,
		last_name,
		language_code,
		first_message_at)
		values ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(ctx, query,
		user.ID,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
		time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUser(ctx context.Context, id int64) (*model.User, error) {
	query := `
		select users.id, 
		       users.username,
		       users.first_name,
		       users.last_name,
		       users.language_code,
		       users.wallet,
		       users.first_message_at
		from tonflow.users
		where id = $1
	`
	user := &model.User{Wallet: &tonclient.Wallet{}}

	err := db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.LanguageCode,
		&user.Wallet.Address,
		&user.FirstMessageAt,
	)
	if err != nil {
		log.Error("GetUser(): %v", err)
		return nil, err
	}

	return user, nil
}

func (db *DB) AddWallet(ctx context.Context, wallet *tonclient.Wallet, userID int64) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)

		} else {
			tx.Commit(ctx)
		}
	}()

	query := `
		insert into tonflow.wallets (
		address, 
		version,
		seed,
		created_at)
		values ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, query, wallet.Address, wallet.Version, wallet.Seed, time.Now())
	if err != nil {
		return err
	}

	query = `
		update tonflow.users
		set wallet = $1
		where id = $2
	`
	_, err = tx.Exec(ctx, query, wallet.Address, userID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetWallet(ctx context.Context, address string) (*tonclient.Wallet, error) {
	query := `
		select address,
		       version,
		       seed 
		from tonflow.wallets 
		where address = $1
	`
	wallet := &tonclient.Wallet{}
	err := db.QueryRow(ctx, query, address).Scan(&wallet.Address, &wallet.Version, &wallet.Seed)
	if err != nil {
		log.Error("GetWallet(): %v", err)
		return nil, err
	}

	return wallet, nil
}
