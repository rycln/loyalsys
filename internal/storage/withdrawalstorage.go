package storage

import (
	"context"
	"database/sql"

	"github.com/rycln/loyalsys/internal/models"
)

type WithdrawalStorage struct {
	db *sql.DB
}

func NewWithdrawalStorage(db *sql.DB) *WithdrawalStorage {
	return &WithdrawalStorage{
		db: db,
	}
}

func (s *WithdrawalStorage) GetWithdrawalsByUserID(ctx context.Context, uid models.UserID) ([]*models.WithdrawalDB, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetWithdrawalsByUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var withdrawals []*models.WithdrawalDB
	for rows.Next() {
		var withdrawal models.WithdrawalDB
		withdrawal.UserID = uid
		err = rows.Scan(&withdrawal.ID, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, &withdrawal)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if withdrawals == nil {
		return nil, newErrNoWithdrawal(ErrNoWithdrawal)
	}
	return withdrawals, nil
}
