package worker

import (
	"context"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

type syncAPI interface {
	GetOrderFromAccrual(context.Context, string) (*models.OrderAccrual, error)
}

type syncStorager interface {
	GetInconclusiveOrderNums(context.Context) ([]string, error)
	UpdateOrdersBatch(context.Context, []*models.OrderDB) error
}

type OrderSyncWorker struct {
	api     syncAPI
	storage syncStorager
	pool    int
	timeout time.Duration
}

func NewOrderSyncWorker(api syncAPI, storage syncStorager, pool int, timeout time.Duration) *OrderSyncWorker {
	return &OrderSyncWorker{
		api:     api,
		storage: storage,
		pool:    pool,
		timeout: timeout,
	}
}

func (worker *OrderSyncWorker) Run(cancelCtx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-cancelCtx.Done():
			return
		case <-ticker.C:
			worker.updateOrdersBatch(cancelCtx)
		}
	}
}

func (worker *OrderSyncWorker) updateOrdersBatch(cancelCtx context.Context) {
	ctx, cancel := context.WithTimeout(cancelCtx, worker.timeout)
	defer cancel()
	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctx)
	if err != nil {
		logger.Log.Debug("worker error", zap.Error(err))
		return
	}
	if len(orderNums) == 0 {
		return
	}

	var updatedOrders []*models.OrderDB

	numsChan := orderNumbersGenerator(cancelCtx, orderNums)
	resultChans := worker.ordersFanOut(cancelCtx, numsChan)
	resultCh := ordersFanIn(cancelCtx, resultChans)

	for result := range resultCh {
		select {
		case <-cancelCtx.Done():
			return
		case <-resultCh:
			if result.err != nil {
				//добавить обработку ошибки retry-after и логгирование
				logger.Log.Debug("worker error", zap.Error(err))
				continue
			}
			updatedOrders = append(updatedOrders, result.order)
		}
	}

	ctx, cancel = context.WithTimeout(cancelCtx, worker.timeout)
	defer cancel()
	err = worker.storage.UpdateOrdersBatch(ctx, updatedOrders)
	if err != nil {
		logger.Log.Debug("worker error", zap.Error(err))
		return
	}
}
