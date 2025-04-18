package worker

import (
	"errors"
	"testing"
	"time"

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

func TestOrderSyncWorker_Run(t *testing.T) {
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

	mAPI := mocks.NewMocksyncAPI(ctrl)
	mStrg := mocks.NewMocksyncStorager(ctrl)
	testCfg := NewSyncWorkerConfigBuilder().
		WithTimeout(testTimeout).
		WithTickerPeriod(testTickerPeriod).
		Build()
	worker := NewOrderSyncWorker(mAPI, mStrg, testCfg)
	stopCh := make(chan struct{})

	t.Run("valid test", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(testOrderNums, nil)
		for i, testOrder := range testOrders {
			mAPI.EXPECT().GetOrderFromAccrual(gomock.Any(), testOrderNums[i]).Return(testOrder, nil)
		}
		mStrg.EXPECT().UpdateOrdersBatch(gomock.Any(), gomock.Any()).Return(nil)

		err := worker.updateOrdersBatch(stopCh)
		assert.NoError(t, err)
	})

	t.Run("get orders num error", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(nil, errTest)

		err := worker.updateOrdersBatch(stopCh)
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("zero len nums list", func(t *testing.T) {
		var zeroOrderNums []string
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(zeroOrderNums, nil)

		err := worker.updateOrdersBatch(stopCh)
		assert.NoError(t, err)
	})

	t.Run("get orders num error", func(t *testing.T) {
		mStrg.EXPECT().GetInconclusiveOrderNums(gomock.Any()).Return(testOrderNums, nil)
		for i, testOrder := range testOrders {
			mAPI.EXPECT().GetOrderFromAccrual(gomock.Any(), testOrderNums[i]).Return(testOrder, nil)
		}
		mStrg.EXPECT().UpdateOrdersBatch(gomock.Any(), gomock.Any()).Return(errTest)

		err := worker.updateOrdersBatch(stopCh)
		assert.ErrorIs(t, err, errTest)
	})
}
