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
	s := NewWithdrawalService(mStrg)

	t.Run("valid test", func(t *testing.T) {
		testWithdrawals := []*models.WithdrawalDB{
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
