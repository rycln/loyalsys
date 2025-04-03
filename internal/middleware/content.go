package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ContentTypeChecker(allowedTypes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentType := c.Get("Content-Type")

		mimeType := strings.Split(contentType, ";")[0]
		mimeType = strings.TrimSpace(mimeType)

		for _, allowedType := range allowedTypes {
			if mimeType == allowedType {
				return c.Next()
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}
