package worker

import (
	"context"
	"sync"

	"github.com/rycln/loyalsys/internal/models"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

const (
	ordersChanBuffer = 1024
)

type syncAPI interface {
	getAPI
}

type syncStorager interface {
	getStorager
	updateStorager
}

type OrderSyncWorker struct {
	getter  *orderGetWorker
	updater *orderUpdateWorker
}

func NewOrderSyncWorker(api syncAPI, storage syncStorager, cfg *SyncWorkerConfig) *OrderSyncWorker {
	return &OrderSyncWorker{
		getter:  newOrderGetWorker(api, storage, cfg),
		updater: newOrderUpdateWorker(storage, cfg),
	}
}

func (worker *OrderSyncWorker) Run(ctx context.Context) chan struct{} {
	doneCh := make(chan struct{})

	var wg sync.WaitGroup

	ordersCh := make(chan *models.OrderDB, ordersChanBuffer)

	worker.getter.run(ctx, &wg, ordersCh)
	worker.updater.run(ctx, &wg, ordersCh)

	go func() {
		wg.Wait()
		defer close(doneCh)
	}()

	return doneCh
}
