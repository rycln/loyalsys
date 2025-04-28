package handlers

import (
	"bytes"
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

func TestPostOrderHandler_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mService := mocks.NewMockpostOrderServicer(ctrl)
	mJWT := mocks.NewMockpostOrderJWT(ctrl)

	postOrderHandler := NewPostOrderHandler(mService, mJWT)

	app := fiber.New()
	app.Post("/", postOrderHandler)

	t.Run("valid test", func(t *testing.T) {
		order := &models.Order{
			Number: validLuhnString,
			UserID: testUserID,
		}
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(nil)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

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

		mErr := mocks.NewMockerrOrderExists(ctrl)
		mErr.EXPECT().IsErrOrderExists().Return(true)
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(mErr)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

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

		mErr := mocks.NewMockerrWrongNum(ctrl)
		mErr.EXPECT().IsErrWrongNum().Return(true)
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(mErr)

		bodyReader := bytes.NewReader([]byte(order.Number))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

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

		mErr := mocks.NewMockerrOrderConflict(ctrl)
		mErr.EXPECT().IsErrOrderConflict().Return(true)
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(mErr)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

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
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(testUserID, nil)
		mService.EXPECT().SaveOrder(gomock.Any(), order).Return(errTest)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("jwt error", func(t *testing.T) {
		mJWT.EXPECT().ParseIDFromAuthHeader(fmt.Sprintf("Bearer %s", testJWTString)).Return(models.UserID(0), errTest)

		bodyReader := bytes.NewReader([]byte(validLuhnString))
		request := httptest.NewRequest(fiber.MethodPost, "/", bodyReader)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testJWTString))

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}
