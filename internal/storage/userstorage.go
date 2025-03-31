package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrConflict = errors.New("login already registered")
)

type UserDB struct {
	ID           int64
	Login        string
	PasswordHash string
}

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (us *UserStorage) AddUser(ctx context.Context, user *UserDB) error {
	tx, err := us.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlAddUser, user.Login, user.PasswordHash)
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict
		}
		return err
	}
	return tx.Commit()
}
