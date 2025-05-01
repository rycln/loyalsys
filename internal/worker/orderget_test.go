package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/worker/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	testTimeout      = time.Duration(1) * time.Second
	testTickerPeriod = time.Duration(1) * time.Second
)

var errTest = errors.New("test error")

func Test_orderGetWorker_getOrders(t *testing.T) {
	defer leaktest.Check(t)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testOrderNums := []string{
		"123",
		"456",
		"789",
	}

	testOrders := []*models.OrderAccrual{
		{
			Number:  testOrderNums[0],
			Status:  "some status #1",
			Accrual: 0,
		},
		{
			Number:  testOrderNums[1],
			Status:  "some status #2",
			Accrual: 10,
		},
		{
			Number:  testOrderNums[2],
			Status:  "some status #2",
			Accrual: 20,
		},
	}

	var testCh = make(chan *models.OrderDB, 10)

	mAPI := mocks.NewMockgetAPI(ctrl)
	mStrg := mocks.NewMockgetStorager(ctrl)
	testCfg := NewSyncWorkerConfigBuilder().
		WithTimeout(testTimeout).
		WithTickerPeriod(testTickerPeriod).
		Build()
	worker := newOrderGetWorker(mAPI, mStrg, testCfg)

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(testOrderNums, nil)
		for i, testOrder := range testOrders {
			mAPI.EXPECT().GetOrderFromAccrual(gomock.Any(), testOrderNums[i]).Return(testOrder, nil)
		}

		err := worker.getOrders(context.Background(), testCh)
		assert.NoError(t, err)

	})

	t.Run("get order nums some error", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(nil, errTest)

		err := worker.getOrders(context.Background(), testCh)
		assert.Error(t, err)
	})

	t.Run("get order no nums error", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(nil, nil)

		err := worker.getOrders(context.Background(), testCh)
		assert.ErrorIs(t, err, errNoOrderNums)
	})

	t.Run("retry after error", func(t *testing.T) {
		orderNum := []string{
			"123",
		}

		mErr := mocks.NewMockerrRetryAfter(ctrl)
		mErr.EXPECT().IsErrRetryAfter().Return(true)

		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(orderNum, nil)
		mAPI.EXPECT().GetOrderFromAccrual(gomock.Any(), gomock.Any()).Return(nil, mErr)

		err := worker.getOrders(context.Background(), testCh)
		assert.Error(t, err, mErr)
	})
}
