package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rycln/loyalsys/internal/models"

	"github.com/go-resty/resty/v2"
)

type OrderUpdateClient struct {
	client  *resty.Client
	baseURL string
}

func NewOrderUpdateClient(client *resty.Client, baseURL string, timeout time.Duration) *OrderUpdateClient {
	return &OrderUpdateClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *OrderUpdateClient) GetOrderFromAccrual(ctx context.Context, num string) (*models.OrderAccrual, error) {
	res, err := c.client.R().SetContext(ctx).SetPathParams(map[string]string{
		"orderNum": num,
	}).Get(c.baseURL + "/api/orders/{orderNum}")
	if err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}

	if res.StatusCode() == http.StatusOK {
		var order models.OrderAccrual
		err = json.Unmarshal(res.Body(), &order)
		if err != nil {
			return nil, fmt.Errorf("client error: %v", err)
		}
		return &order, nil
	}
	if res.StatusCode() == http.StatusNoContent {
		return nil, ErrNoContent
	}
	if res.StatusCode() == http.StatusTooManyRequests {
		dur, err := time.ParseDuration(res.Header().Get("Retry-After") + "s")
		if err != nil {
			return nil, fmt.Errorf("client error: %v", err)
		}
		return nil, newErrRetryAfter(dur, ErrTooManyRequests)
	}
	return nil, fmt.Errorf("client received an unexpected status code: %s", res.Status())
}
