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

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type loginServicer interface {
	UserAuth(context.Context, *models.User) (models.UserID, error)
}

type LoginHandler struct {
	loginService loginServicer
	cfg          *config.Cfg
}

func NewLoginHandler(loginService loginServicer, cfg *config.Cfg) func(*fiber.Ctx) error {
	h := &LoginHandler{
		loginService: loginService,
		cfg:          cfg,
	}
	return h.handle
}

func (h *LoginHandler) handle(c *fiber.Ctx) error {
	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = user.Validate()
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}

	uid, err := h.loginService.UserAuth(c.Context(), &user)
	if errors.Is(err, storage.ErrNoUser) || errors.Is(err, auth.ErrWrongPassword) {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	jwt, err := auth.NewJWTString(uid, h.cfg.Key)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	c.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return c.SendStatus(fiber.StatusOK)
}
