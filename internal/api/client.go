package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"

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
		"baseURL":  c.baseURL,
		"orderNum": num,
	}).Get("{baseURL}/{orderNum}")
	if err != nil {
		logger.Log.Debug("client error", zap.Error(err))
		return nil, err
	}

	if res.StatusCode() == http.StatusOK {

	}
	var order models.OrderAccrual
	err = json.Unmarshal(res.Body(), &order)
	if err != nil {
		logger.Log.Debug("client error", zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
}
