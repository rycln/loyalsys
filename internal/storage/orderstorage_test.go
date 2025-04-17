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

	strg := NewOrderStorage(db)

	testOrder := &models.Order{
		Number: testOrderNum,
		UserID: testUserID,
	}

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO orders").WithArgs(testOrder.Number, testOrder.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := strg.AddOrder(context.Background(), testOrder)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO orders").WithArgs(testOrder.Number, testOrder.UserID).WillReturnError(errTest)

		err := strg.AddOrder(context.Background(), testOrder)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderStorage_GetOrderByNum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	strg := NewOrderStorage(db)

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

		orderDB, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.NoError(t, err)
		assert.Equal(t, testOrder, orderDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no order error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, number, user_id, status, accrual, created_at FROM orders").WithArgs(testOrder.Number).WillReturnError(sql.ErrNoRows)

		_, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.ErrorIs(t, err, ErrNoOrder)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, number, user_id, status, accrual, created_at FROM orders").WithArgs(testOrder.Number).WillReturnError(errTest)

		_, err := strg.GetOrderByNum(context.Background(), testOrder.Number)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderStorage_GetOrdersByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	strg := NewOrderStorage(db)

	testOrder := &models.OrderDB{
		Number:    testOrderNum,
		Status:    "some status",
		Accrual:   0,
		CreatedAt: testCreatedAt.String(),
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"number", "status", "accrual", "created_at"}).
			AddRow(testOrder.Number, testOrder.Status, testOrder.Accrual, testOrder.CreatedAt)
		mock.ExpectQuery("SELECT number, status, accrual, created_at FROM orders").WithArgs(testUserID).WillReturnRows(rows)

		orderDB, err := strg.GetOrdersByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testOrder, orderDB[0])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT number, status, accrual, created_at FROM orders").WithArgs(testUserID).WillReturnError(errTest)

		_, err := strg.GetOrdersByUserID(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty response", func(t *testing.T) {
		rows := mock.NewRows([]string{"number", "status", "accrual", "created_at"})
		mock.ExpectQuery("SELECT number, status, accrual, created_at FROM orders").WillReturnRows(rows)

		_, err := strg.GetOrdersByUserID(context.Background(), testUserID)
		assert.ErrorIs(t, err, ErrNoOrder)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderStorage_GetInconclusiveOrderNums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	strg := NewOrderStorage(db)

	orderNums := []string{
		"123",
		"456",
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"number"}).AddRow(orderNums[0]).AddRow(orderNums[1])
		mock.ExpectQuery("SELECT number FROM orders").WillReturnRows(rows)

		nums, err := strg.GetInconclusiveOrderNums(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, orderNums, nums)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT number FROM orders").WillReturnError(errTest)

		_, err := strg.GetInconclusiveOrderNums(context.Background())
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderStorage_UpdateOrdersBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	strg := NewOrderStorage(db)

	testOrders := []*models.OrderDB{
		{
			Number:    "123",
			Status:    "some status",
			Accrual:   10,
			CreatedAt: time.Now().String(),
		},
		{
			Number:    "456",
			Status:    "some status",
			Accrual:   10,
			CreatedAt: time.Now().String(),
		},
	}

	t.Run("valid test", func(t *testing.T) {
		mock.ExpectBegin()
		mockStmt := mock.ExpectPrepare("UPDATE orders")
		for _, testOrder := range testOrders {
			mockStmt.ExpectExec().WithArgs(testOrder.Status, testOrder.Accrual, testOrder.Number).WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		err := strg.UpdateOrdersBatch(context.Background(), testOrders)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("begin error", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errTest)

		err := strg.UpdateOrdersBatch(context.Background(), testOrders)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("prepare error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectPrepare("UPDATE orders").WillReturnError(errTest)
		mock.ExpectRollback()

		err := strg.UpdateOrdersBatch(context.Background(), testOrders)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("exec error", func(t *testing.T) {
		mock.ExpectBegin()
		mockStmt := mock.ExpectPrepare("UPDATE orders")

		mockStmt.ExpectExec().WithArgs(testOrders[0].Status, testOrders[0].Accrual, testOrders[0].Number).WillReturnError(errTest)

		mock.ExpectRollback()

		err := strg.UpdateOrdersBatch(context.Background(), testOrders)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
