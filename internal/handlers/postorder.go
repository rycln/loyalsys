package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/config"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/services"
	"go.uber.org/zap"
)

type postOrderServicer interface {
	SaveOrder(context.Context, *models.Order) error
}

type PostOrderHandler struct {
	postOrderService postOrderServicer
	cfg              *config.Cfg
}

func NewPostOrderHandler(postOrderService postOrderServicer, cfg *config.Cfg) func(*fiber.Ctx) error {
	h := &PostOrderHandler{
		postOrderService: postOrderService,
		cfg:              cfg,
	}
	return h.handle
}

func (h *PostOrderHandler) handle(c *fiber.Ctx) error {
	if !c.Is("text/plain") {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	t := c.Locals("user").(*jwt.Token)
	claims := t.Claims.(auth.JwtClaims)
	uid := claims.GetUserID()

	order := &models.Order{
		Number: string(c.Body()),
		UserID: uid,
	}
	err := h.postOrderService.SaveOrder(c.Context(), order)
	if errors.Is(err, services.ErrOrderExists) {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusOK)
	}
	if errors.Is(err, services.ErrWrongNum) {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if errors.Is(err, services.ErrOrderConflict) {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusConflict)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
