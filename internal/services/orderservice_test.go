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

func TestOrderService_SaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockorderStorager(ctrl)
	s := NewOrderService(mStrg)

	mErrNoOrder := mocks.NewMockerrNoOrder(ctrl)

	t.Run("valid test", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mErrNoOrder.EXPECT().IsErrNoOrder().Return(true)
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, mErrNoOrder)
		mStrg.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(nil)

		err := s.SaveOrder(context.Background(), testOrder)
		assert.NoError(t, err)
	})

	t.Run("luhn validation error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: "12345",
			UserID: testUserID,
		}

		err := s.SaveOrder(context.Background(), testOrder)
		assert.ErrorIs(t, err, ErrWrongNum)
	})

	t.Run("AddOrder error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mErrNoOrder.EXPECT().IsErrNoOrder().Return(true)
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, mErrNoOrder)
		mStrg.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(errTest)

		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, errTest)
	})

	t.Run("GetOrderByNum error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, errTest)

		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, errTest)
	})

	t.Run("order exists error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		testOrderDB := &models.OrderDB{
			UserID: testUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(testOrderDB, nil)

		err := s.SaveOrder(context.Background(), testOrder)
		assert.ErrorIs(t, err, ErrOrderExists)
	})

	t.Run("order conflict error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		testOrderDB := &models.OrderDB{
			UserID: testOtherUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(testOrderDB, nil)

		err := s.SaveOrder(context.Background(), testOrder)
		assert.ErrorIs(t, err, ErrOrderConflict)
	})
}

func TestOrderService_GetUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockorderStorager(ctrl)
	s := NewOrderService(mStrg)

	t.Run("valid test", func(t *testing.T) {
		testOrders := []*models.OrderDB{
			{
				Number:    "123",
				Status:    "some status",
				Accrual:   10,
				CreatedAt: time.Now().String(),
			},
			{
				Number:    "456",
				Status:    "some status",
				Accrual:   10,
				CreatedAt: time.Now().String(),
			},
		}

		mStrg.EXPECT().GetOrdersByUserID(context.Background(), testUserID).Return(testOrders, nil)

		orders, err := s.GetUserOrders(context.Background(), testUserID)
		assert.Equal(t, testOrders, orders)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mStrg.EXPECT().GetOrdersByUserID(context.Background(), testUserID).Return(nil, errTest)

		_, err := s.GetUserOrders(context.Background(), testUserID)
		assert.Error(t, err)
	})
}
