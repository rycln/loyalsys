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

func TestGetWithdrawalsHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockgetWithdrawalsServicer(ctrl)
	mJWT := mocks.NewMockgetWithdrawalsJWT(ctrl)

	getWithdrawalsHandler := NewGetWithdrawalsHandler(mService, mJWT)

	app := fiber.New()
	app.Get("/", getWithdrawalsHandler)

	t.Run("valid test", func(t *testing.T) {
		testWithdrawals := []*models.Withdrawal{
			{
				ID:          1,
				Order:       "123",
				UserID:      testUserID,
				Sum:         10,
				ProcessedAt: time.Now().String(),
			},
			{
				ID:          2,
				Order:       "456",
				UserID:      testUserID,
				Sum:         5,
				ProcessedAt: time.Now().String(),
			},
		}

		testWithdrawalsJSON, err := json.Marshal(&testWithdrawals)
		require.NoError(t, err)

		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserWithdrawals(gomock.Any(), testUserID).Return(testWithdrawals, nil)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.JSONEq(t, string(testWithdrawalsJSON), string(body))
	})

	t.Run("no withdrawal error", func(t *testing.T) {
		mErr := mocks.NewMockerrNoWithdrawal(ctrl)
		mErr.EXPECT().IsErrNoWithdrawal().Return(true)
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserWithdrawals(gomock.Any(), testUserID).Return(nil, mErr)

		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusNoContent, res.StatusCode)
	})

	t.Run("some error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().GetUserWithdrawals(gomock.Any(), testUserID).Return(nil, errTest)

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
