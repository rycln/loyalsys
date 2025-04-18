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
	var totalAccrual, totalWithdrawn float64
	err := row.Scan(&totalAccrual, &totalWithdrawn)
	if err != nil {
		return nil, err
	}
	current := totalAccrual - totalWithdrawn
	balance := &models.Balance{
		UserID:    uid,
		Current:   current,
		Withdrawn: totalWithdrawn,
	}
	return balance, nil
}
