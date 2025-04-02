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
	ErrNoUser   = errors.New("user does not exist")
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) AddUser(ctx context.Context, user *models.UserDB) (models.UserID, error) {
	row := s.db.QueryRowContext(ctx, sqlAddUser, user.Login, user.PasswordHash)
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

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*models.UserDB, error) {
	row := s.db.QueryRowContext(ctx, sqlGetUserByLogin, login)
	var userDB models.UserDB
	err := row.Scan(&userDB.ID, &userDB.Login, &userDB.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoUser
	}
	if err != nil {
		return nil, err
	}
	return &userDB, nil
}
