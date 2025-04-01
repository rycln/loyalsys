package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
	"go.uber.org/zap"
)

type regServicer interface {
	CreateUser(context.Context, *models.User) (models.UserID, error)
}

type RegisterHandler struct {
	regService regServicer
	cfg        *config.Cfg
}

func NewRegisterHandler(regService regServicer, cfg *config.Cfg) func(*fiber.Ctx) error {
	h := &RegisterHandler{
		regService: regService,
		cfg:        cfg,
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
		logger.Log.Info("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	uid, err := h.regService.CreateUser(c.Context(), &user)
	if errors.Is(err, storage.ErrConflict) {
		return c.SendStatus(fiber.StatusConflict)
	}
	if err != nil {
		logger.Log.Info("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	jwt, err := auth.NewJWTString(uid, h.cfg.Key)
	if err != nil {
		logger.Log.Info("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	c.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return c.SendStatus(fiber.StatusOK)
}
