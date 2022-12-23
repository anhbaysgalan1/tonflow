package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"tonflow/internal/service/bot/model"
	"tonflow/internal/service/ton"
	"tonflow/internal/storage"
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

func NewStorage(cfg *Config) (storage.Storage, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSL)

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

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

func (db *DB) CheckUser(ctx context.Context, u *model.User) (bool, error) {
	query := `select id from users where id = $1`
	var userID int64
	err := db.QueryRow(ctx, query, u.ID).Scan(&userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		query = `insert into users (id, username, first_name, last_name, language_code, wallet, first_message_at, last_message_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(ctx, query, u.ID, u.Username, u.FirstName, u.LastName, u.LanguageCode, u.Wallet, time.Now(), time.Now())
		if err != nil {
			return false, err
		}
		return false, nil
	}

	query = `update users set last_message_at = $2 where id = $1`
	_, err = db.Exec(ctx, query, u.ID, time.Now())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) AddTransaction(ctx context.Context, tr ton.Transaction) error {
	query := `insert into transactions ("source", "hash", "value", "comment") values ($1, $2, $3, $4)`
	_, err := db.Exec(ctx, query, tr.Source, tr.Hash, tr.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddWallet(ctx context.Context, wallet *ton.Wallet, userID int64) error {
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

	query := `insert into wallets ("address", "version", "seed", "created_at") values ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, query, wallet.Address, wallet.Version, wallet.Seed, time.Now())
	if err != nil {
		return err
	}

	query = `update users set wallet = $1 where id = $2`
	_, err = tx.Exec(ctx, query, wallet.Address, userID)
	if err != nil {
		return err
	}

	query = `update users set last_message_at = $2 where id = $1`
	_, err = tx.Exec(ctx, query, userID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserWallet(ctx context.Context, userID int64) (string, error) {
	query := `select wallet from users where id = $1`
	var wallet string
	err := db.QueryRow(ctx, query, userID).Scan(&wallet)
	if err != nil {
		return "", err
	}

	return wallet, nil
}

func (db *DB) GetUserTransactions(ctx context.Context, userID int64) ([]*ton.Transaction, error) {
	wallet, err := db.GetUserWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet == "" {
		return nil, errors.New("user have no wallet")
	}

	query := `select * from transactions where source = $1`
	rows, err := db.Query(ctx, query, wallet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trs := make([]*ton.Transaction, 0)
	for rows.Next() {
		var r *ton.Transaction
		err = rows.Scan(&r.Source, &r.Hash, &r.Value, &r.Comment)
		if err != nil {
			return nil, err
		}
		trs = append(trs, r)
	}

	return trs, nil
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
		return "", err
	}
	return id, nil
}

func (db *DB) GetAllPictures(ctx context.Context) ([]string, error) {
	query := `select id from pictures`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	IDs := make([]string, 0)
	for rows.Next() {
		var ID string
		err = rows.Scan(&ID)
		if err != nil {
			return nil, err
		}
		IDs = append(IDs, ID)
	}
	return IDs, nil
}
