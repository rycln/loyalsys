package worker

import (
	"context"
	"sync"

	"github.com/rycln/loyalsys/internal/models"
)

type updateOrderResult struct {
	order *models.OrderDB
	err   error
}

func (worker *OrderSyncWorker) ordersFanOut(stopCh <-chan struct{}, inputNumCh <-chan string) []chan updateOrderResult {
	channels := make([]chan updateOrderResult, worker.cfg.fanOutPool)

	for i := 0; i < worker.cfg.fanOutPool; i++ {
		resultCh := worker.getUpdatedOrderByNum(stopCh, inputNumCh)
		channels[i] = resultCh
	}

	return channels
}

func (worker *OrderSyncWorker) getUpdatedOrderByNum(stopCh <-chan struct{}, inputNumCh <-chan string) chan updateOrderResult {
	resultCh := make(chan updateOrderResult)

	go func() {
		defer close(resultCh)

		for num := range inputNumCh {
			var orderDB *models.OrderDB

			ctx, cancel := context.WithTimeout(context.Background(), worker.cfg.timeout)
			defer cancel()
			orderAccrual, err := worker.api.GetOrderFromAccrual(ctx, num)
			if err == nil {
				orderDB = &models.OrderDB{
					Number:  orderAccrual.Number,
					Status:  orderAccrual.Status,
					Accrual: orderAccrual.Accrual,
				}
			}

			result := updateOrderResult{
				order: orderDB,
				err:   err,
			}

			select {
			case <-stopCh:
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func ordersFanIn(stopCh <-chan struct{}, channels []chan updateOrderResult) chan updateOrderResult {
	resultCh := make(chan updateOrderResult)

	var wg sync.WaitGroup

	for _, ch := range channels {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for result := range chClosure {
				select {
				case <-stopCh:
					return
				case resultCh <- result:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	return resultCh
}

func orderNumbersGenerator(stopCh <-chan struct{}, nums []string) chan string {
	inputNumCh := make(chan string)

	go func() {
		defer close(inputNumCh)

		for _, num := range nums {
			select {
			case <-stopCh:
				return
			case inputNumCh <- num:
			}
		}
	}()

	return inputNumCh
}
