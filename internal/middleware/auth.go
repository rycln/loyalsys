package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func NoTokenChecker() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
