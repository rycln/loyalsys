package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/handlers/mocks"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalanceHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockgetBalanceServicer(ctrl)
	mJWT := mocks.NewMockgetBalanceJWT(ctrl)

	getBalanceHandler := NewGetBalanceHandler(mService, mJWT)

	app := fiber.New()
	app.Get("/", getBalanceHandler)

	t.Run("valid test", func(t *testing.T) {
		testBalance := &models.Balance{
			UserID:    testUserID,
			Current:   10,
			Withdrawn: 20,
		}

		testBalanceJSON, err := json.Marshal(&testBalance)
		require.NoError(t, err)

		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserBalance(gomock.Any(), testUserID).Return(testBalance, nil)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.JSONEq(t, string(testBalanceJSON), string(body))
	})

	t.Run("some error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserBalance(gomock.Any(), testUserID).Return(nil, errTest)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("jwt error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(models.UserID(0), errTest)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}
