package storage

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rycln/loyalsys/db"
)

const (
	maxOpenConns    = 0 //unlimited
	maxIdleConns    = 10
	maxIdleTime     = time.Duration(3) * time.Minute
	maxConnLifetime = 0 //unlimited
)

func NewDB(dsn string) (*sql.DB, error) {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(maxOpenConns)
	database.SetMaxIdleConns(maxIdleConns)
	database.SetConnMaxIdleTime(maxIdleTime)
	database.SetConnMaxLifetime(maxConnLifetime)

	goose.SetBaseFS(db.Migrations)
	err = goose.Up(database, "migrations")
	if err != nil {
		return nil, err
	}
	return database, nil
}
