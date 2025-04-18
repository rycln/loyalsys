package storage

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStorage_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewUserStorage(db)

	testUser := &models.UserDB{
		Login:        "test",
		PasswordHash: "hashed_password",
	}

	expectedQuery := regexp.QuoteMeta(sqlAddUser)

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id"}).AddRow(testUserID)
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		uid, err := strg.AddUser(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("conflict error", func(t *testing.T) {
		var pgErr = &pgconn.PgError{
			Code: pgerrcode.IntegrityConstraintViolation,
		}

		rows := mock.NewRows([]string{"id"}).AddRow(testUserID).RowError(0, pgErr)
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		_, err = strg.AddUser(context.Background(), testUser)
		assert.ErrorIs(t, err, ErrLoginConflict)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		rows := mock.NewRows([]string{"id"}).AddRow(testUserID).RowError(0, errTest)
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		_, err = strg.AddUser(context.Background(), testUser)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserStorage_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	strg := NewUserStorage(db)

	testUser := &models.UserDB{
		ID:           testUserID,
		Login:        "test",
		PasswordHash: "hashed_password",
	}

	expectedQuery := regexp.QuoteMeta(sqlGetUserByLogin)

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "login", "password_hash"}).AddRow(testUser.ID, testUser.Login, testUser.PasswordHash)
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login).WillReturnRows(rows)

		userDB, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.NoError(t, err)
		assert.Equal(t, testUser, userDB)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no user error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login).WillReturnError(sql.ErrNoRows)

		_, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.ErrorIs(t, err, ErrNoUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).WithArgs(testUser.Login).WillReturnError(errTest)

		_, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
