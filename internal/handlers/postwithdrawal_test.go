package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rycln/loyalsys/internal/handlers/mocks"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostWithdrawalHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockpostWithdrawalServicer(ctrl)
	mJWT := mocks.NewMockpostWithdrawalJWT(ctrl)

	postWithdrawalHandler := NewPostWithdrawalHandler(mService, mJWT)

	app := fiber.New()
	app.Post("/", postWithdrawalHandler)

	withdrawal := &models.Withdrawal{
		Order:  validLuhnString,
		UserID: testUserID,
		Sum:    10,
	}
	testWithdrawalsJSON, err := json.Marshal(withdrawal)
	require.NoError(t, err)

	t.Run("valid test", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().WithdrawalProcessing(gomock.Any(), withdrawal).Return(nil)

		bodyReader := bytes.NewReader([]byte(testWithdrawalsJSON))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("jwt error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(models.UserID(0), errTest)

		bodyReader := bytes.NewReader([]byte(testWithdrawalsJSON))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("wrong json body", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)

		bodyReader := bytes.NewReader([]byte("wrong json"))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("not enough currency error", func(t *testing.T) {
		mErr := mocks.NewMockerrNotEnoughCurrency(ctrl)

		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mErr.EXPECT().IsErrNotEnoughCurrency().Return(true)
		mService.EXPECT().WithdrawalProcessing(gomock.Any(), withdrawal).Return(mErr)

		bodyReader := bytes.NewReader([]byte(testWithdrawalsJSON))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusPaymentRequired, res.StatusCode)
	})

	t.Run("luhn validation error", func(t *testing.T) {
		mErr := mocks.NewMockerrWrongOrderNum(ctrl)

		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mErr.EXPECT().IsErrWrongOrderNum().Return(true)
		mService.EXPECT().WithdrawalProcessing(gomock.Any(), withdrawal).Return(mErr)

		bodyReader := bytes.NewReader([]byte(testWithdrawalsJSON))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
	})

	t.Run("some error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().WithdrawalProcessing(gomock.Any(), withdrawal).Return(errTest)

		bodyReader := bytes.NewReader([]byte(testWithdrawalsJSON))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}
