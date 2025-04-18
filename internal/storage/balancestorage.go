package storage

import (
	"context"
	"database/sql"

	"github.com/rycln/loyalsys/internal/models"
)

type BalanceStorage struct {
	db *sql.DB
}

func NewBalanceStorage(db *sql.DB) *BalanceStorage {
	return &BalanceStorage{
		db: db,
	}
}

func (s *BalanceStorage) GetBalanceByUserID(ctx context.Context, uid models.UserID) (*models.Balance, error) {
	row := s.db.QueryRowContext(ctx, sqlGetBalanceByUserID, uid)
	var balance models.Balance
	balance.UserID = uid
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}
