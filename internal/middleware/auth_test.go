package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoTokenChecker(t *testing.T) {
	app := fiber.New()
	app.Get("/", NoTokenChecker(), SendStausOK)

	t.Run("valid test", func(t *testing.T) {
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)
		request.Header.Set("Authorization", "some token")

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, fiber.StatusOK)
	})

	t.Run("wrong content type", func(t *testing.T) {
		request := httptest.NewRequest(fiber.MethodGet, "/", nil)

		res, err := app.Test(request, -1)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, fiber.StatusUnauthorized)
	})
}
