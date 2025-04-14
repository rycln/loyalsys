package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"

	"github.com/go-resty/resty/v2"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrNoContent       = errors.New("no content")
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
		logger.Log.Debug("client error", zap.Error(err))
		return nil, err
	}

	//debug
	logger.Log.Debug("client response", zap.String("status code", res.Status()))

	if res.StatusCode() == http.StatusOK {
		var order models.OrderAccrual
		err = json.Unmarshal(res.Body(), &order)
		if err != nil {
			logger.Log.Debug("client error", zap.Error(err))
			return nil, err
		}
		return &order, nil
	}
	if res.StatusCode() == http.StatusNoContent {
		return nil, newErrorNoContent(ErrNoContent)
	}
	if res.StatusCode() == http.StatusTooManyRequests {
		dur, err := time.ParseDuration(res.Header().Get("Retry-After") + "s")
		if err != nil {
			logger.Log.Debug("client error", zap.Error(err))
			return nil, err
		}
		return nil, newErrorTooManyRequests(dur, ErrTooManyRequests)
	}
	return nil, fmt.Errorf("the client received an unexpected status code: %s", res.Status())
}
