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

func orderNumbersGenerator(ctx context.Context, nums []string) <-chan string {
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

func (worker *orderGetWorker) ordersFanOut(ctx context.Context, inputNumCh <-chan string) []<-chan updateOrderResult {
	channels := make([]<-chan updateOrderResult, worker.cfg.fanOutPool)

	for i := 0; i < worker.cfg.fanOutPool; i++ {
		resultCh := worker.getUpdatedOrderByNum(ctx, inputNumCh)
		channels[i] = resultCh
	}

	return channels
}

func (worker *orderGetWorker) getUpdatedOrderByNum(ctx context.Context, inputNumCh <-chan string) <-chan updateOrderResult {
	resultCh := make(chan updateOrderResult)

	go func() {
		defer close(resultCh)

		for num := range inputNumCh {
			var orderDB *models.OrderDB

			ctxAPI, cancel := context.WithTimeout(ctx, worker.cfg.timeout)
			defer cancel()

			orderAccrual, err := worker.api.GetOrderFromAccrual(ctxAPI, num)
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
			case <-ctx.Done():
				return
			case resultCh <- result:
			}
		}
	}()

	return resultCh
}

func ordersFanIn(ctx context.Context, channels []<-chan updateOrderResult) <-chan updateOrderResult {
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

func ordersResultDispatcher(ctx context.Context, resultCh <-chan updateOrderResult, ordersCh chan<- *models.OrderDB) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		for result := range resultCh {
			select {
			case <-ctx.Done():
				return
			default:
				if result.err != nil {
					errCh <- result.err
				} else {
					ordersCh <- result.order
				}
			}
		}
	}()

	return errCh
}
