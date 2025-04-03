package services

import (
	"context"

	"github.com/rycln/loyalsys/internal/models"
)

type orderStorager interface {
	AddOrder(context.Context, models.Order) error
	GetOrderByNum(context.Context, string) (*models.OrderDB, error)
}

type OrderService struct {
	strg orderStorager
}

func NewOrderService(strg orderStorager) *OrderService {
	return &OrderService{strg: strg}
}

func (s *OrderService) SaveOrder(ctx context.Context, order *models.Order) error {
	//сделать проверку Луном
}
