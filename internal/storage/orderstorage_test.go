package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	testOrderNum       = "12345"
	testOrderID  int64 = 1
)

var testCreatedAt = time.Now()

func TestOrderStorage_AddOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testOrder := &models.Order{
		Number: testOrderNum,
		UserID: testUserID,
	}

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO orders").WithArgs(testOrder.Number, testOrder.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		strg := NewOrderStorage(db)
		err := strg.AddOrder(context.Background(), testOrder)
		assert.NoError(t, err)
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO orders").WithArgs(testOrder.Number, testOrder.UserID).WillReturnError(errTest)

		strg := NewOrderStorage(db)
		err := strg.AddOrder(context.Background(), testOrder)
		assert.Error(t, err)
	})
}

func TestOrderStorage_GetOrderByNum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testOrder := &models.OrderDB{
		ID:        testOrderID,
		Number:    testOrderNum,
		UserID:    testUserID,
		Status:    "some status",
		Accrual:   0,
		CreatedAt: testCreatedAt.String(),
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "number", "user_id", "status", "accrual", "created_at"}).AddRow(testOrder.ID, testOrder.Number, testOrder.UserID, testOrder.Status, testOrder.Accrual, testOrder.CreatedAt)
		mock.ExpectQuery("SELECT id, number, user_id, status, accrual, created_at FROM orders").WithArgs(testOrder.Number).WillReturnRows(rows)

		strg := NewOrderStorage(db)
		orderDB, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.NoError(t, err)
		assert.Equal(t, testOrder, orderDB)
	})

	t.Run("no order error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, number, user_id, status, accrual, created_at FROM orders").WithArgs(testOrder.Number).WillReturnError(sql.ErrNoRows)

		strg := NewOrderStorage(db)
		_, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.ErrorIs(t, err, ErrNoOrder)
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, number, user_id, status, accrual, created_at FROM orders").WithArgs(testOrder.Number).WillReturnError(errTest)

		strg := NewOrderStorage(db)
		_, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.Error(t, err)
	})
}
