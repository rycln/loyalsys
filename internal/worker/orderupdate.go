package worker

import (
	"context"
	"sync"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

const (
	ordersMaxBufSize = 1024
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type updateStorager interface {
	UpdateOrdersBatch(context.Context, []*models.OrderDB) error
}

type orderUpdateWorker struct {
	storage updateStorager
	cfg     *SyncWorkerConfig
}

func newOrderUpdateWorker(storage updateStorager, cfg *SyncWorkerConfig) *orderUpdateWorker {
	return &orderUpdateWorker{
		storage: storage,
		cfg:     cfg,
	}
}

func (worker *orderUpdateWorker) run(ctx context.Context, wg *sync.WaitGroup, orderCh <-chan *models.OrderDB) {
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(worker.cfg.tickerPeriod)
		defer ticker.Stop()

		var updatedOrdersBuf = make([]*models.OrderDB, 0, ordersMaxBufSize)

		for {
			select {
			case <-ctx.Done():
				return
			case order, ok := <-orderCh:
				if !ok {
					return
				}
				if len(updatedOrdersBuf) == ordersMaxBufSize {
					err := worker.updateOrders(ctx, updatedOrdersBuf)
					if err != nil {
						logger.Log.Debug("update orders error", zap.Error(err))
						continue
					}
					updatedOrdersBuf = updatedOrdersBuf[:0]
				}
				updatedOrdersBuf = append(updatedOrdersBuf, order)
			case <-ticker.C:
				if len(updatedOrdersBuf) > 0 {
					err := worker.updateOrders(ctx, updatedOrdersBuf)
					if err != nil {
						logger.Log.Debug("update orders error", zap.Error(err))
						continue
					}
					updatedOrdersBuf = updatedOrdersBuf[:0]
				}
			}
		}
	}()
}

func (worker *orderUpdateWorker) updateOrders(ctx context.Context, updatedOrders []*models.OrderDB) error {
	ctxDB, cancel := context.WithTimeout(ctx, worker.cfg.timeout)
	defer cancel()

	err := worker.storage.UpdateOrdersBatch(ctxDB, updatedOrders)
	if err != nil {
		return err
	}
	return nil
}
