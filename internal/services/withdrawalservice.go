package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type withdrawalStorager interface {
	GetWithdrawalsByUserID(context.Context, models.UserID) ([]*models.WithdrawalDB, error)
}

type WithdrawalService struct {
	strg withdrawalStorager
}

func NewWithdrawalService(strg withdrawalStorager) *WithdrawalService {
	return &WithdrawalService{strg: strg}
}

func (s *WithdrawalService) GetUserWithdrawals(ctx context.Context, uid models.UserID) ([]*models.WithdrawalDB, error) {
	withdrawals, err := s.strg.GetWithdrawalsByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
