package storage

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceStorage_GetBalanceByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewBalanceStorage(db)

	testBalance := &models.Balance{
		UserID:    testUserID,
		Current:   10,
		Withdrawn: 20,
	}

	expectedQuery := regexp.QuoteMeta(sqlGetBalanceByUserID)

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"current", "withdrawn"}).AddRow(testBalance.Current, testBalance.Withdrawn)
		mock.ExpectQuery(expectedQuery).WithArgs(testBalance.UserID).WillReturnRows(rows)

		balance, err := strg.GetBalanceByUserID(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testBalance, balance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WithArgs(testBalance.UserID).WillReturnError(errTest)

		_, err := strg.GetBalanceByUserID(context.Background(), testUserID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
