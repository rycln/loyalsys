package middleware

import "github.com/gofiber/fiber/v2"

func SendStausOK(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
