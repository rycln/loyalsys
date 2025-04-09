package middleware

import (
	"strings"

	"slices"

	"github.com/gofiber/fiber/v2"
)

func ContentTypeChecker(allowedTypes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")

		mimeType := strings.Split(contentType, ";")[0]
		mimeType = strings.TrimSpace(mimeType)

		if slices.Contains(allowedTypes, mimeType) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}
