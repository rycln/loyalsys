package services

import (
	"context"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/rycln/loyalsys/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type withdrawalStorager interface {
	GetWithdrawalsByUserID(context.Context, models.UserID) ([]*models.Withdrawal, error)
	AddWithdrawal(context.Context, *models.Withdrawal) error
}

type balanceServicer interface {
	GetUserBalance(context.Context, models.UserID) (*models.Balance, error)
}

type WithdrawalService struct {
	strg    withdrawalStorager
	balance balanceServicer
}

func NewWithdrawalService(strg withdrawalStorager, balance balanceServicer) *WithdrawalService {
	return &WithdrawalService{
		strg:    strg,
		balance: balance,
	}
}

func (s *WithdrawalService) WithdrawalProcessing(ctx context.Context, withdrawal *models.Withdrawal) error {
	err := goluhn.Validate(withdrawal.Order)
	if err != nil {
		return newErrWrongOrderNum(ErrWrongOrderNum)
	}

	balance, err := s.balance.GetUserBalance(ctx, withdrawal.UserID)
	if err != nil {
		return err
	}

	if balance.Current < withdrawal.Sum {
		return newErrNotEnoughCurrency(ErrNotEnoughCurrency)
	}

	err = s.strg.AddWithdrawal(ctx, withdrawal)
	if err != nil {
		return err
	}
	return nil
}

func (s *WithdrawalService) GetUserWithdrawals(ctx context.Context, uid models.UserID) ([]*models.Withdrawal, error) {
	withdrawals, err := s.strg.GetWithdrawalsByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
