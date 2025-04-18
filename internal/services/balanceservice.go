package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type balanceStorager interface {
	GetBalanceByUserID(context.Context, models.UserID) (*models.Balance, error)
}

type BalanceService struct {
	strg balanceStorager
}

func NewBalanceService(strg balanceStorager) *BalanceService {
	return &BalanceService{strg: strg}
}

func (s *BalanceService) GetUserBalance(ctx context.Context, uid models.UserID) (*models.Balance, error) {
	balance, err := s.strg.GetBalanceByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
