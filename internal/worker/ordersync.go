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

type errRetryAfter interface {
	error
	GetRetryAfterDuration() time.Duration
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

func (worker *OrderSyncWorker) Run(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := worker.updateOrdersBatch(ctx)
			if e, ok := err.(errRetryAfter); ok {
				dur := e.GetRetryAfterDuration()
				ticker.Reset(dur)
				logger.Log.Info("worker retry after", zap.Duration("duration, sec:", dur))
				continue
			}
			if err != nil {
				logger.Log.Debug("worker error", zap.Error(err))
			}
		}
	}
}

func (worker *OrderSyncWorker) updateOrdersBatch(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, worker.timeout)
	defer cancel()
	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctx)
	if err != nil {
		return err
	}
	if len(orderNums) == 0 {
		return nil
	}

	numsChan := orderNumbersGenerator(ctx, orderNums)
	resultChans := worker.ordersFanOut(ctx, numsChan)
	resultCh := ordersFanIn(ctx, resultChans)

	var updatedOrders []*models.OrderDB

	for result := range resultCh {
		select {
		case <-ctx.Done():
			return nil
		case <-resultCh:
			if err, ok := result.err.(errRetryAfter); ok {
				return err
			}
			if result.err != nil {
				logger.Log.Debug("pipeline error", zap.Error(err))
				continue
			}
			updatedOrders = append(updatedOrders, result.order)
		}
	}

	ctx, cancel = context.WithTimeout(ctx, worker.timeout)
	defer cancel()
	err = worker.storage.UpdateOrdersBatch(ctx, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
