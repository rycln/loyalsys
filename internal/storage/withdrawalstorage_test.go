package storage

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testWithdrawalID    = 1
	testWithdrawalOrder = "12345"
	testWithdrawalSum   = float64(10)
)

var testProcessedAt = time.Now().String()

func TestWithdrawalStorage_GetWithdrawalsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewWithdrawalStorage(db)

	testWithdrawal := &models.WithdrawalDB{
		ID:          testWithdrawalID,
		Order:       testWithdrawalOrder,
		UserID:      testUserID,
		Sum:         testWithdrawalSum,
		ProcessedAt: testProcessedAt,
	}

	expectedQuery := regexp.QuoteMeta(sqlGetWithdrawalsByUserID)

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "order", "sum", "processed_at"}).
			AddRow(testWithdrawal.ID, testWithdrawal.Order, testWithdrawal.Sum, testWithdrawal.ProcessedAt)
		mock.ExpectQuery(expectedQuery).WithArgs(testUserID).WillReturnRows(rows)

		withdrawalsDB, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testWithdrawal, withdrawalsDB[0])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WithArgs(testUserID).WillReturnError(errTest)

		_, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty response", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "order", "sum", "processed_at"})
		mock.ExpectQuery(expectedQuery).WithArgs(testUserID).WillReturnRows(rows)

		_, err := strg.GetWithdrawalsByUserID(context.Background(), testUserID)
		assert.ErrorIs(t, err, ErrNoWithdrawal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
