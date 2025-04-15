package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/handlers/mocks"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockregServicer(ctrl)

	registerHandler := NewRegisterHandler(mService, testKey)

	app := fiber.New()
	app.Post("/", registerHandler)

	t.Run("valid test", func(t *testing.T) {
		testUser := &models.User{
			Login:    testUserLogin,
			Password: testUserPassword,
		}
		mService.EXPECT().CreateUser(gomock.Any(), testUser).Return(testUserID, nil)

		body, err := json.Marshal(testUser)
		require.NoError(t, err)
		bodyReader := bytes.NewReader(body)
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.NotEmpty(t, res.Header.Get("Authorization"))
	})

	t.Run("wrong json body", func(t *testing.T) {
		body := "wrong json string"
		bodyReader := bytes.NewReader([]byte(body))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("user conflict", func(t *testing.T) {
		testUser := &models.User{
			Login:    testUserLogin,
			Password: testUserPassword,
		}
		mService.EXPECT().CreateUser(gomock.Any(), testUser).Return(models.UserID(0), storage.ErrLoginConflict)

		body, err := json.Marshal(testUser)
		require.NoError(t, err)
		bodyReader := bytes.NewReader(body)
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusConflict, res.StatusCode)
	})

	t.Run("some error", func(t *testing.T) {
		testUser := &models.User{
			Login:    testUserLogin,
			Password: testUserPassword,
		}
		mService.EXPECT().CreateUser(gomock.Any(), testUser).Return(models.UserID(0), errTest)

		body, err := json.Marshal(testUser)
		require.NoError(t, err)
		bodyReader := bytes.NewReader(body)
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("invalid user", func(t *testing.T) {
		testUser := &models.User{
			Login:    "",
			Password: testUserPassword,
		}

		body, err := json.Marshal(testUser)
		require.NoError(t, err)
		bodyReader := bytes.NewReader(body)
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})
}
