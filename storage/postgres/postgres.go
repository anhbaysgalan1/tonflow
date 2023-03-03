package postgres

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"tonflow/model"
	"tonflow/storage"
)

type DB struct {
	*pgxpool.Pool
	memory map[string]int64
}

func NewConnection(URI string) (storage.Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool, err := pgxpool.New(ctx, URI)
	if err != nil {
		return nil, err
	}

	query := `
		select users.id,
		       users.wallet
		from tonflow.users
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	memory := make(map[string]int64)
	defer rows.Close()
	for rows.Next() {
		var (
			addr   string
			userID int64
		)

		err = rows.Scan(&userID, &addr)
		if err != nil {
			return nil, err
		}
		memory[addr] = userID
	}

	return &DB{
		pool,
		memory,
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
		created_at)
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
		       users.created_at
		from tonflow.users
		where id = $1
	`
	user := &model.User{Wallet: &model.Wallet{}}

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
		return nil, err
	}

	return user, nil
}

func (db *DB) AddWallet(ctx context.Context, wallet *model.Wallet, userID int64) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		tx.Commit(ctx)
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

	db.memory[wallet.Address] = userID
	return nil
}

func (db *DB) GetWallet(ctx context.Context, address string) (*model.Wallet, error) {
	query := `
		select address,
		       version,
		       seed 
		from tonflow.wallets 
		where address = $1
	`
	wallet := &model.Wallet{}
	err := db.QueryRow(ctx, query, address).Scan(
		&wallet.Address,
		&wallet.Version,
		&wallet.Seed,
	)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (db *DB) GetInMemoryWallets() map[string]int64 {
	return db.memory
}
