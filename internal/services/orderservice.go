package services

import (
	"context"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/rycln/loyalsys/internal/models"
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

type errNoOrder interface {
	error
	IsErrNoOrder() bool
}

func (s *OrderService) SaveOrder(ctx context.Context, order *models.Order) error {
	err := goluhn.Validate(order.Number)
	if err != nil {
		return newErrWrongNum(ErrWrongNum)
	}
	checkOrder, err := s.strg.GetOrderByNum(ctx, order.Number)
	if e, ok := err.(errNoOrder); ok && e.IsErrNoOrder() {
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
		return newErrOrderExists(ErrOrderExists)
	}
	return newErrOrderConflict(ErrOrderConflict)
}

func (s *OrderService) GetUserOrders(ctx context.Context, uid models.UserID) ([]*models.OrderDB, error) {
	orders, err := s.strg.GetOrdersByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
