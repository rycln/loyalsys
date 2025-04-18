package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBalanceService_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockbalanceStorager(ctrl)
	s := NewBalanceService(mStrg)

	t.Run("valid test", func(t *testing.T) {
		testBalance := &models.Balance{
			UserID:    testUserID,
			Current:   10,
			Withdrawn: 20,
		}

		mStrg.EXPECT().GetBalanceByUserID(context.Background(), testUserID).Return(testBalance, nil)

		balance, err := s.GetUserBalance(context.Background(), testUserID)
		assert.Equal(t, testBalance, balance)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetBalanceByUserID(context.Background(), testUserID).Return(nil, errTest)

		_, err := s.GetUserBalance(context.Background(), testUserID)
		assert.Error(t, err)
	})
}
