package worker

import (
	"context"
	"errors"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

var errNoOrderNums = errors.New("nothing to update")

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

func (worker *OrderSyncWorker) Run(ctx context.Context) chan struct{} {
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		ticker := time.NewTicker(worker.cfg.tickerPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateCtx, ctxCancel := context.WithCancel(ctx)

				err := worker.updateOrdersBatch(updateCtx)
				if err != nil {
					logger.Log.Debug("worker error", zap.Error(err))
				}
				if e, ok := err.(errRetryAfter); ok && e.IsErrRetryAfter() {
					dur := e.GetRetryAfterDuration()
					ticker.Reset(dur)
					logger.Log.Info("worker retry after", zap.Duration("duration, sec:", dur))
				}

				ctxCancel()
			}
		}
	}()

	return doneCh
}

func (worker *OrderSyncWorker) updateOrdersBatch(ctx context.Context) error {
	orderNums, err := worker.getOrderNums(ctx)
	if err != nil {
		return err
	}

	numsChan := orderNumbersGenerator(ctx, orderNums)
	resultChans := worker.ordersFanOut(ctx, numsChan)
	resultCh := ordersFanIn(ctx, resultChans)

	updatedOrders, err := updateConsumer(ctx, resultCh)
	if err != nil {
		return err
	}

	err = worker.updateOrders(ctx, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}

func (worker *OrderSyncWorker) getOrderNums(ctx context.Context) ([]string, error) {
	ctxDB, cancel := context.WithTimeout(ctx, worker.cfg.timeout)
	defer cancel()

	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctxDB)
	if err != nil {
		return nil, err
	}
	if len(orderNums) == 0 {
		return nil, errNoOrderNums
	}
	return orderNums, nil
}

func updateConsumer(ctx context.Context, resultCh chan updateOrderResult) ([]*models.OrderDB, error) {
	var updatedOrders []*models.OrderDB

	for result := range resultCh {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-resultCh:
			if err, ok := result.err.(errRetryAfter); ok && err.IsErrRetryAfter() {
				return nil, err
			}
			if result.err != nil {
				logger.Log.Debug("pipeline error", zap.Error(result.err))
				continue
			}
			updatedOrders = append(updatedOrders, result.order)
		}
	}

	return updatedOrders, nil
}

func (worker *OrderSyncWorker) updateOrders(ctx context.Context, updatedOrders []*models.OrderDB) error {
	ctxDB, cancel := context.WithTimeout(ctx, worker.cfg.timeout)
	defer cancel()

	err := worker.storage.UpdateOrdersBatch(ctxDB, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
