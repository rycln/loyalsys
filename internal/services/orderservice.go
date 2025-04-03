package services

import (
	"context"
	"errors"

	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
)

var (
	ErrOrderExists   = errors.New("order already registered by user")
	ErrOrderConflict = errors.New("order already registered by other user")
)

type orderStorager interface {
	AddOrder(context.Context, *models.Order) error
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
