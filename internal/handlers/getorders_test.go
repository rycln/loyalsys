package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/handlers/mocks"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrderHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockgetOrdersServicer(ctrl)
	mJWT := mocks.NewMockgetOrdersJWT(ctrl)

	getOrdersHandler := NewGetOrdersHandler(mService, mJWT)

	app := fiber.New()
	app.Get("/", getOrdersHandler)

	t.Run("valid test", func(t *testing.T) {
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

		testOrdersJSON, err := json.Marshal(&testOrders)
		require.NoError(t, err)

		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserOrders(gomock.Any(), testUserID).Return(testOrders, nil)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.JSONEq(t, string(testOrdersJSON), string(body))
	})

	t.Run("no order error", func(t *testing.T) {
		mErr := mocks.NewMockerrNoOrder(ctrl)
		mErr.EXPECT().IsErrNoOrder().Return(true)
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserOrders(gomock.Any(), testUserID).Return(nil, mErr)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusNoContent, res.StatusCode)
	})

	t.Run("some error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserOrders(gomock.Any(), testUserID).Return(nil, errTest)

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
