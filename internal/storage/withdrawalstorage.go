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

func (s *WithdrawalStorage) AddWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error {
	_, err := s.db.ExecContext(ctx, sqlAddWithdrawal, withdrawal.Order, withdrawal.UserID, withdrawal.Sum)
	if err != nil {
		return err
	}
	return nil
}

func (s *WithdrawalStorage) GetWithdrawalsByUserID(ctx context.Context, uid models.UserID) ([]*models.Withdrawal, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetWithdrawalsByUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var withdrawals []*models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal
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
