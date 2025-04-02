package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
)

var (
	ErrConflict = errors.New("login already registered")
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (us *UserStorage) AddUser(ctx context.Context, user *models.UserDB) (models.UserID, error) {
	row := us.db.QueryRowContext(ctx, sqlAddUser, user.Login, user.PasswordHash)
	var uid int64
	err := row.Scan(&uid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return 0, ErrConflict
		}
		return 0, err
	}
	return models.UserID(uid), nil
}
