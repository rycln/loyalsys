package db

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(database *sql.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(database, "migrations"); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return err
	}

	return nil
}
