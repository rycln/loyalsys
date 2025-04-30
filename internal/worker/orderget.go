package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

var errNoOrderNums = errors.New("nothing to update")

type getAPI interface {
	GetOrderFromAccrual(context.Context, string) (*models.OrderAccrual, error)
}

type getStorager interface {
	GetInconclusiveOrderNums(context.Context) ([]string, error)
}

type errRetryAfter interface {
	error
	IsErrRetryAfter() bool
	GetRetryAfterDuration() time.Duration
}

type orderGetWorker struct {
	api     getAPI
	storage getStorager
	cfg     *SyncWorkerConfig
}

func newOrderGetWorker(api getAPI, storage getStorager, cfg *SyncWorkerConfig) *orderGetWorker {
	return &orderGetWorker{
		api:     api,
		storage: storage,
		cfg:     cfg,
	}
}

func (worker *orderGetWorker) run(ctx context.Context, wg *sync.WaitGroup, orderCh chan<- *models.OrderDB) {
	defer wg.Done()

	ticker := time.NewTicker(worker.cfg.tickerPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := worker.getOrders(ctx, orderCh)
			if err != nil {
				logger.Log.Debug("get orders error", zap.Error(err))
			}
			if e, ok := err.(errRetryAfter); ok && e.IsErrRetryAfter() {
				dur := e.GetRetryAfterDuration()
				ticker.Reset(dur)
				logger.Log.Info("worker retry after", zap.Duration("duration, sec:", dur))
			}
		}
	}
}

func (worker *orderGetWorker) getOrders(ctx context.Context, orderCh chan<- *models.OrderDB) error {
	ctxGet, cancelGet := context.WithCancel(ctx)
	defer cancelGet()

	orderNums, err := worker.getOrderNums(ctxGet)
	if err != nil {
		return err
	}

	numsChan := orderNumbersGenerator(ctxGet, orderNums)
	resultChans := worker.ordersFanOut(ctxGet, numsChan)
	resultCh := ordersFanIn(ctxGet, resultChans)
	errCh := ordersResultDispatcher(ctxGet, resultCh, orderCh)

	for err := range errCh {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err, ok := err.(errRetryAfter); ok && err.IsErrRetryAfter() {
				return err
			}
			if err != nil {
				logger.Log.Debug("pipeline error", zap.Error(err))
				continue
			}
		}
	}

	return nil
}

func (worker *orderGetWorker) getOrderNums(ctx context.Context) ([]string, error) {
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
