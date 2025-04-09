package services

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services/mocks"
	"github.com/rycln/loyalsys/internal/storage"
	"github.com/stretchr/testify/assert"
)

const validLuhnString = "4512812345678909"

func TestOrderService_SaveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStrg := mocks.NewMockorderStorager(ctrl)

	t.Run("valid test", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, storage.ErrNoOrder)
		mStrg.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(nil)

		s := NewOrderService(mStrg)
		err := s.SaveOrder(context.Background(), testOrder)
		assert.NoError(t, err)
	})

	t.Run("luhn validation error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: "12345",
			UserID: testUserID,
		}

		s := NewOrderService(mStrg)
		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, ErrWrongNum)
	})

	t.Run("AddOrder error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, storage.ErrNoOrder)
		mStrg.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(errTest)

		s := NewOrderService(mStrg)
		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, errTest)
	})

	t.Run("GetOrderByNum error", func(t *testing.T) {
		testOrder := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mStrg.EXPECT().GetOrderByNum(gomock.Any(), testOrder.Number).Return(nil, errTest)

		s := NewOrderService(mStrg)
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

		s := NewOrderService(mStrg)
		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, ErrOrderExists)
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

		s := NewOrderService(mStrg)
		err := s.SaveOrder(context.Background(), testOrder)
		assert.Equal(t, err, ErrOrderConflict)
	})
}
