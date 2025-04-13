package services

import (
	"context"
	"errors"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
)

var (
	ErrWrongNum      = errors.New("luhn algorithm validation failed")
	ErrOrderExists   = errors.New("order already registered by user")
	ErrOrderConflict = errors.New("order already registered by other user")
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type orderStorager interface {
	AddOrder(context.Context, *models.Order) error
	GetOrderByNum(context.Context, string) (*models.OrderDB, error)
	GetOrdersByUserID(context.Context, models.UserID) ([]*models.OrderDB, error)
}

type OrderService struct {
	strg orderStorager
}

func NewOrderService(strg orderStorager) *OrderService {
	return &OrderService{strg: strg}
}

func (s *OrderService) SaveOrder(ctx context.Context, order *models.Order) error {
	err := goluhn.Validate(order.Number)
	if err != nil {
		return ErrWrongNum
	}
	checkOrder, err := s.strg.GetOrderByNum(ctx, order.Number)
	if errors.Is(err, storage.ErrNoOrder) {
		err := s.strg.AddOrder(ctx, order)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	if checkOrder.UserID == order.UserID {
		return ErrOrderExists
	}
	return ErrOrderConflict
}

func (s *OrderService) GetUserOrders(ctx context.Context, uid models.UserID) ([]*models.OrderDB, error) {
	orders, err := s.strg.GetOrdersByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
