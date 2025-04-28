package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
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
		return nil, newErrNoOrder(ErrNoOrder)
	}
	if err != nil {
		return nil, err
	}
	return &orderDB, nil
}

func (s *OrderStorage) GetOrdersByUserID(ctx context.Context, uid models.UserID) ([]*models.OrderDB, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetOrdersByUserID, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*models.OrderDB
	for rows.Next() {
		var order models.OrderDB
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return nil, newErrNoOrder(ErrNoOrder)
	}
	return orders, nil
}

func (s *OrderStorage) GetInconclusiveOrderNums(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetInconclusiveOrderNums)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nums []string
	for rows.Next() {
		var num string
		err = rows.Scan(&num)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return nums, nil
}

func (s *OrderStorage) UpdateOrdersBatch(ctx context.Context, orders []*models.OrderDB) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, sqlUpdateOrdersBatch)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, order := range orders {
		if _, err := stmt.ExecContext(ctx, order.Status, order.Accrual, order.Number); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
