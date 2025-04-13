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

func (worker *OrderSyncWorker) GetUpdatedOrderByNum(ctx context.Context, inputNumCh <-chan string) chan updateOrderResult {
	resultCh := make(chan updateOrderResult)

	go func() {
		defer close(resultCh)

		for num := range inputNumCh {
			timeoutCtx, cancel := context.WithTimeout(ctx, worker.timeout)
			order, err := worker.api.GetOrderFromAccrual(timeoutCtx, num)
			cancel()

			orderDB := &models.OrderDB{
				Number:  order.Number,
				Status:  order.Status,
				Accrual: order.Accrual,
			}

			result := updateOrderResult{
				order: orderDB,
				err:   err,
			}

			select {
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func (worker *OrderSyncWorker) ordersFanOut(ctx context.Context, inputNumCh <-chan string) []chan updateOrderResult {
	channels := make([]chan updateOrderResult, worker.pool)

	for i := 0; i < worker.pool; i++ {
		resultCh := worker.GetUpdatedOrderByNum(ctx, inputNumCh)
		channels[i] = resultCh
	}

	return channels
}

func ordersFanIn(ctx context.Context, channels []chan updateOrderResult) chan updateOrderResult {
	resultCh := make(chan updateOrderResult)

	var wg sync.WaitGroup

	for _, ch := range channels {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for result := range chClosure {
				select {
				case <-ctx.Done():
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

func orderNumbersGenerator(ctx context.Context, nums []string) chan string {
	inputNumCh := make(chan string)

	go func() {
		defer close(inputNumCh)

		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case inputNumCh <- num:
			}
		}
	}()

	return inputNumCh
}
