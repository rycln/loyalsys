package handlers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/handlers/mocks"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validLuhnString = "4512812345678909"

func TestPostOrderHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockpostOrderServicer(ctrl)
	mJWT := mocks.NewMockpostOrderJWT(ctrl)

	postOrderHandler := NewPostOrderHandler(mService, mJWT)

	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(testKey)},
	}))
	app.Post("/", postOrderHandler)

	claims := jwt.MapClaims{
		"uid": testUserID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testKey))
	require.NoError(t, err)

	t.Run("valid test", func(t *testing.T) {
		order := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(nil)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusAccepted, res.StatusCode)
	})

	t.Run("order exists", func(t *testing.T) {
		order := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(services.ErrOrderExists)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("wrong order number", func(t *testing.T) {
		order := &models.Order{
			Number: "12345",
			UserID: testUserID,
		}
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(services.ErrWrongNum)

		bodyReader := bytes.NewReader([]byte(order.Number))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
	})

	t.Run("order conflict", func(t *testing.T) {
		order := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(services.ErrOrderConflict)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusConflict, res.StatusCode)
	})

	t.Run("some error", func(t *testing.T) {
		order := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(errTest)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}
