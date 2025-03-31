package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
)

type regServicer interface {
	CreateUser(context.Context, *models.User) (models.UserID, error)
}

type RegisterHandler struct {
	regService regServicer
}

func NewRegisterHandler(regService regServicer) func(*fiber.Ctx) error {
	h := &RegisterHandler{
		regService: regService,
	}
	return h.handle
}

func (h *RegisterHandler) handle(c *fiber.Ctx) error {
	if !c.Is("json") {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := h.regService.CreateUser(c.Context(), &user)
	if errors.Is(err, storage.ErrConflict) {
		return c.SendStatus(fiber.StatusConflict)
	}
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

}
