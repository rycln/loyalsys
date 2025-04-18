package worker

import (
	"context"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

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

func (worker *OrderSyncWorker) Run(stopCh <-chan struct{}) chan struct{} {
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		ticker := time.NewTicker(worker.cfg.tickerPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				err := worker.updateOrdersBatch(stopCh)
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
	}()

	return doneCh
}

func (worker *OrderSyncWorker) updateOrdersBatch(stopCh <-chan struct{}) error {
	ctxDB, cancelDB := context.WithTimeout(context.Background(), worker.cfg.timeout)
	defer cancelDB()

	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctxDB)
	if err != nil {
		return err
	}
	if len(orderNums) == 0 {
		return nil
	}

	numsChan := orderNumbersGenerator(stopCh, orderNums)
	resultChans := worker.ordersFanOut(stopCh, numsChan)
	resultCh := ordersFanIn(stopCh, resultChans)

	var updatedOrders []*models.OrderDB

	for result := range resultCh {
		select {
		case <-stopCh:
			return nil
		case <-resultCh:
			if err, ok := result.err.(errRetryAfter); ok && err.IsErrRetryAfter() {
				return err
			}
			if result.err != nil {
				logger.Log.Debug("pipeline error", zap.Error(result.err))
				continue
			}
			updatedOrders = append(updatedOrders, result.order)
		}
	}

	ctxDB, cancelDB = context.WithTimeout(context.Background(), worker.cfg.timeout)
	defer cancelDB()

	err = worker.storage.UpdateOrdersBatch(ctxDB, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
