package worker

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/worker/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_orderUpdateWorker_updateOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCfg := NewSyncWorkerConfigBuilder().
		WithTimeout(testTimeout).
		WithTickerPeriod(testTickerPeriod).
		Build()

	orders := []*models.OrderDB{
		{
			Number:  "123",
			Status:  "some status #1",
			Accrual: 0,
		},
		{
			Number:  "456",
			Status:  "some status #2",
			Accrual: 10,
		},
		{
			Number:  "789",
			Status:  "some status #2",
			Accrual: 20,
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mStrg := mocks.NewMockupdateStorager(ctrl)

		mStrg.EXPECT().UpdateOrdersBatch(gomock.Any(), orders).Return(nil)

		worker := newOrderUpdateWorker(mStrg, testCfg)

		err := worker.updateOrders(context.Background(), orders)
		assert.NoError(t, err)
	})
}
