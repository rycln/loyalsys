package storage

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConns    = 0 //unlimited
	maxIdleConns    = 10
	maxIdleTime     = time.Duration(3) * time.Minute
	maxConnLifetime = 0 //unlimited
)

func NewDB(uri string) (*sql.DB, error) {
	database, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(maxOpenConns)
	database.SetMaxIdleConns(maxIdleConns)
	database.SetConnMaxIdleTime(maxIdleTime)
	database.SetConnMaxLifetime(maxConnLifetime)

	return database, nil
}
