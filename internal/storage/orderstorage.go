package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
)

var (
	ErrOrderExists   = errors.New("order already registered by user")
	ErrOrderConflict = errors.New("order already registered by other user")
	ErrNoOrder       = errors.New("order does not exist")
)

type OrderStorage struct {
	db *sql.DB
}

func NewOrderStorage(db *sql.DB) *OrderStorage {
	return &OrderStorage{db: db}
}

func (s *OrderStorage) AddOrder(ctx context.Context, order *models.Order) error {
	_, err := s.db.ExecContext(ctx, sqlAddOrder, order.Number, order.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderStorage) GetOrderByNum(ctx context.Context, number string) (*models.OrderDB, error) {
	row := s.db.QueryRowContext(ctx, sqlGetOrderByNum, number)
	var orderDB models.OrderDB
	err := row.Scan(&orderDB.ID, &orderDB.Number, &orderDB.UserID, &orderDB.Status, &orderDB.Accrual, &orderDB.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoOrder
	}
	if err != nil {
		return nil, err
	}
	return &orderDB, nil
}
