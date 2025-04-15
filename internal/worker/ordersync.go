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
	IsErrRetryAfter() bool
	GetRetryAfterDuration() time.Duration
}

type OrderSyncWorker struct {
	api     syncAPI
	storage syncStorager
	cfg     *SyncWorkerConfig
}

func NewOrderSyncWorker(api syncAPI, storage syncStorager, cfg *SyncWorkerConfig) *OrderSyncWorker {
	return &OrderSyncWorker{
		api:     api,
		storage: storage,
		cfg:     cfg,
	}
}

func (worker *OrderSyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(worker.cfg.tickerPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := worker.updateOrdersBatch(ctx)
			if e, ok := err.(errRetryAfter); ok && e.IsErrRetryAfter() {
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
	ctxDB, cancelDB := context.WithTimeout(ctx, worker.cfg.timeout)
	defer cancelDB()

	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctxDB)
	if err != nil {
		return err
	}
	if len(orderNums) == 0 {
		return nil
	}

	ctxPipeline, cancelPipeline := context.WithCancel(ctx)
	defer cancelPipeline()

	numsChan := orderNumbersGenerator(ctxPipeline, orderNums)
	resultChans := worker.ordersFanOut(ctxPipeline, numsChan)
	resultCh := ordersFanIn(ctxPipeline, resultChans)

	var updatedOrders []*models.OrderDB

	for result := range resultCh {
		select {
		case <-ctx.Done():
			return nil
		case <-resultCh:
			if err, ok := result.err.(errRetryAfter); ok && err.IsErrRetryAfter() {
				return err
			}
			if result.err != nil {
				logger.Log.Debug("pipeline error", zap.Error(err))
				continue
			}
			updatedOrders = append(updatedOrders, result.order)
		}
	}

	ctxDB, cancelDB = context.WithTimeout(ctx, worker.cfg.timeout)
	defer cancelDB()

	err = worker.storage.UpdateOrdersBatch(ctxDB, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
