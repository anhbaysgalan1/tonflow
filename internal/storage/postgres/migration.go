package postgres

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migration(databaseURL string) error {
	m, err := migrate.New("file://internal/storage/postgres/migration", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to init database migration: %w", err)
	}

	// uncomment to drop tables on migration
	//err = m.Down()
	//if err != nil && err != migrate.ErrNoChange {
	//	return fmt.Errorf("failed to apply all down migration: %w", err)
	//}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply all up migration: %w", err)
	}

	return nil
}
