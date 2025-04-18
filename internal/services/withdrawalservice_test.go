package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawalService_GetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockwithdrawalStorager(ctrl)
	mBalance := mocks.NewMockbalanceServicer(ctrl)
	s := NewWithdrawalService(mStrg, mBalance)

	t.Run("valid test", func(t *testing.T) {
		testWithdrawals := []*models.Withdrawal{
			{
				ID:          1,
				Order:       "123",
				UserID:      testUserID,
				Sum:         10,
				ProcessedAt: time.Now().String(),
			},
			{
				ID:          2,
				Order:       "456",
				UserID:      testUserID,
				Sum:         5,
				ProcessedAt: time.Now().String(),
			},
		}

		mStrg.EXPECT().GetWithdrawalsByUserID(context.Background(), testUserID).Return(testWithdrawals, nil)

		withdrawals, err := s.GetUserWithdrawals(context.Background(), testUserID)
		assert.Equal(t, testWithdrawals, withdrawals)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetWithdrawalsByUserID(context.Background(), testUserID).Return(nil, errTest)

		_, err := s.GetUserWithdrawals(context.Background(), testUserID)
		assert.Error(t, err)
	})
}

func TestWithdrawalService_WithdrawalProcessing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockwithdrawalStorager(ctrl)
	mBalance := mocks.NewMockbalanceServicer(ctrl)
	s := NewWithdrawalService(mStrg, mBalance)

	testBalance := &models.Balance{
		UserID:    testUserID,
		Current:   10,
		Withdrawn: 20,
	}

	t.Run("valid test", func(t *testing.T) {
		testWithdrawal := &models.Withdrawal{
			ID:          1,
			Order:       validLuhnString,
			UserID:      testUserID,
			Sum:         10,
			ProcessedAt: time.Now().String(),
		}

		mBalance.EXPECT().GetUserBalance(context.Background(), testWithdrawal.UserID).Return(testBalance, nil)
		mStrg.EXPECT().AddWithdrawal(context.Background(), testWithdrawal).Return(nil)

		err := s.WithdrawalProcessing(context.Background(), testWithdrawal)
		assert.NoError(t, err)
	})

	t.Run("luhn error", func(t *testing.T) {
		testWithdrawal := &models.Withdrawal{
			ID:          1,
			Order:       "12345",
			UserID:      testUserID,
			Sum:         10,
			ProcessedAt: time.Now().String(),
		}

		err := s.WithdrawalProcessing(context.Background(), testWithdrawal)
		assert.ErrorIs(t, err, ErrWrongOrderNum)
	})

	t.Run("balance error", func(t *testing.T) {
		testWithdrawal := &models.Withdrawal{
			ID:          1,
			Order:       validLuhnString,
			UserID:      testUserID,
			Sum:         10,
			ProcessedAt: time.Now().String(),
		}

		mBalance.EXPECT().GetUserBalance(context.Background(), testWithdrawal.UserID).Return(nil, errTest)

		err := s.WithdrawalProcessing(context.Background(), testWithdrawal)
		assert.Error(t, err)
	})

	t.Run("not enough currency", func(t *testing.T) {
		testWithdrawal := &models.Withdrawal{
			ID:          1,
			Order:       validLuhnString,
			UserID:      testUserID,
			Sum:         20,
			ProcessedAt: time.Now().String(),
		}

		mBalance.EXPECT().GetUserBalance(context.Background(), testWithdrawal.UserID).Return(testBalance, nil)

		err := s.WithdrawalProcessing(context.Background(), testWithdrawal)
		assert.ErrorIs(t, err, ErrNotEnoughCurrency)
	})

	t.Run("add withdrawal error", func(t *testing.T) {
		testWithdrawal := &models.Withdrawal{
			ID:          1,
			Order:       validLuhnString,
			UserID:      testUserID,
			Sum:         10,
			ProcessedAt: time.Now().String(),
		}

		mBalance.EXPECT().GetUserBalance(context.Background(), testWithdrawal.UserID).Return(testBalance, nil)
		mStrg.EXPECT().AddWithdrawal(context.Background(), testWithdrawal).Return(errTest)

		err := s.WithdrawalProcessing(context.Background(), testWithdrawal)
		assert.Error(t, err)
	})
}
