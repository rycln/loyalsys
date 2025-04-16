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

type loginServicer interface {
	UserAuth(context.Context, *models.User) (models.UserID, error)
}

type loginJWT interface {
	NewJWTString(models.UserID) (string, error)
}

type LoginHandler struct {
	loginService loginServicer
	jwt          loginJWT
}

func NewLoginHandler(loginService loginServicer, jwt loginJWT) func(*fiber.Ctx) error {
	h := &LoginHandler{
		loginService: loginService,
		jwt:          jwt,
	}
	return h.handle
}

type errNoUser interface {
	error
	IsErrNoUser() bool
}

type errWrongPassword interface {
	error
	IsErrWrongPassword() bool
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
	if e, ok := err.(errNoUser); ok && e.IsErrNoUser() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if e, ok := err.(errWrongPassword); ok && e.IsErrWrongPassword() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnauthorized)
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
