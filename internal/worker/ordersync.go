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
				if errors.Is(err, errNoOrderNums) {
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
	orderNums, err := worker.getOrderNums()
	if err != nil {
		return err
	}

	numsChan := orderNumbersGenerator(stopCh, orderNums)
	resultChans := worker.ordersFanOut(stopCh, numsChan)
	resultCh := ordersFanIn(stopCh, resultChans)

	updatedOrders, err := updateConsumer(stopCh, resultCh)
	if err != nil {
		return err
	}

	err = worker.updateOrders(updatedOrders)
	if err != nil {
		return err
	}
	return nil
}

func (worker *OrderSyncWorker) getOrderNums() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), worker.cfg.timeout)
	defer cancel()

	orderNums, err := worker.storage.GetInconclusiveOrderNums(ctx)
	if err != nil {
		return nil, err
	}
	if len(orderNums) == 0 {
		return nil, errNoOrderNums
	}
	return orderNums, nil
}

func updateConsumer(stopCh <-chan struct{}, resultCh chan updateOrderResult) ([]*models.OrderDB, error) {
	var updatedOrders []*models.OrderDB

	for result := range resultCh {
		select {
		case <-stopCh:
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

func (worker *OrderSyncWorker) updateOrders(updatedOrders []*models.OrderDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), worker.cfg.timeout)
	defer cancel()

	err := worker.storage.UpdateOrdersBatch(ctx, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
