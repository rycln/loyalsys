package main

import (
	"database/sql"
	"embed"
	"errors"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	uri := os.Getenv("DATABASE_URI")

	db, err := sql.Open("pgx", uri)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		log.Fatal(err)
	}
}
