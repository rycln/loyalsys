package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type regServicer interface {
	CreateUser(context.Context, *models.User) (models.UserID, error)
}

type regJWT interface {
	NewJWTString(models.UserID) (string, error)
}

type RegisterHandler struct {
	regService regServicer
	jwt        regJWT
}

func NewRegisterHandler(regService regServicer, jwt regJWT) func(*fiber.Ctx) error {
	h := &RegisterHandler{
		regService: regService,
		jwt:        jwt,
	}
	return h.handle
}

type errLoginConflict interface {
	error
	IsErrLoginConflict() bool
}

func (h *RegisterHandler) handle(c *fiber.Ctx) error {
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

	uid, err := h.regService.CreateUser(c.Context(), &user)
	if e, ok := err.(errLoginConflict); ok && e.IsErrLoginConflict() {
		return c.SendStatus(fiber.StatusConflict)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	jwt, err := h.jwt.NewJWTString(uid)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	c.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	return c.SendStatus(fiber.StatusOK)
}
