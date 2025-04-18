package storage

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	testWithdrawalID    = 1
	testWithdrawalOrder = "12345"
	testWithdrawalSum   = float64(10)
)

var testProcessedAt = time.Now().String()

func TestWithdrawalStorage_GetWithdrawalsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	strg := NewWithdrawalStorage(db)

	testWithdrawal := &models.WithdrawalDB{
		ID:          testWithdrawalID,
		Order:       testWithdrawalOrder,
		UserID:      testUserID,
		Sum:         testWithdrawalSum,
		ProcessedAt: testProcessedAt,
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "order", "sum", "processed_at"}).
			AddRow(testWithdrawal.ID, testWithdrawal.Order, testWithdrawal.Sum, testWithdrawal.ProcessedAt)
		mock.ExpectQuery("SELECT id, order, sum, processed_at FROM withdrawals").WithArgs(testUserID).WillReturnRows(rows)

		withdrawalsDB, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testWithdrawal, withdrawalsDB[0])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, order, sum, processed_at FROM withdrawals").WithArgs(testUserID).WillReturnError(errTest)

		_, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty response", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "order", "sum", "processed_at"})
		mock.ExpectQuery("SELECT id, order, sum, processed_at FROM withdrawals").WithArgs(testUserID).WillReturnRows(rows)

		_, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.ErrorIs(t, err, ErrNoWithdrawal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
