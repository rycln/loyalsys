package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserStorage_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testUser := &models.UserDB{
		Login:        "test",
		PasswordHash: "hashed_password",
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id"}).AddRow(testUserID)
		mock.ExpectQuery("INSERT INTO users").WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		strg := NewUserStorage(db)
		uid, err := strg.AddUser(context.Background(), testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, uid)
	})

	t.Run("conflict error", func(t *testing.T) {
		var pgErr = &pgconn.PgError{
			Code: pgerrcode.IntegrityConstraintViolation,
		}

		rows := mock.NewRows([]string{"id"}).AddRow(testUserID).RowError(0, pgErr)
		mock.ExpectQuery("INSERT INTO users").WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		strg := NewUserStorage(db)
		_, err = strg.AddUser(context.Background(), testUser)
		assert.ErrorIs(t, err, ErrLoginConflict)
	})

	t.Run("some error", func(t *testing.T) {
		rows := mock.NewRows([]string{"id"}).AddRow(testUserID).RowError(0, errTest)
		mock.ExpectQuery("INSERT INTO users").WithArgs(testUser.Login, testUser.PasswordHash).WillReturnRows(rows)

		strg := NewUserStorage(db)
		_, err = strg.AddUser(context.Background(), testUser)
		assert.Error(t, err)
	})
}

func TestUserStorage_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testUser := &models.UserDB{
		ID:           testUserID,
		Login:        "test",
		PasswordHash: "hashed_password",
	}

	t.Run("valid test", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "login", "password_hash"}).AddRow(testUser.ID, testUser.Login, testUser.PasswordHash)
		mock.ExpectQuery("SELECT id, login, password_hash FROM users").WithArgs(testUser.Login).WillReturnRows(rows)

		strg := NewUserStorage(db)
		userDB, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.NoError(t, err)
		assert.Equal(t, testUser, userDB)
	})

	t.Run("no user error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, login, password_hash FROM users").WithArgs(testUser.Login).WillReturnError(sql.ErrNoRows)

		strg := NewUserStorage(db)
		_, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.ErrorIs(t, err, ErrNoUser)
	})

	t.Run("some error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, login, password_hash FROM users").WithArgs(testUser.Login).WillReturnError(errTest)

		strg := NewUserStorage(db)
		_, err := strg.GetUserByLogin(context.Background(), testUser.Login)
		assert.Error(t, err)
	})
}
