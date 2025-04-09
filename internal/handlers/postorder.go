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

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

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
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	uid, err := auth.ParseIDFromJWT(token)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	order := &models.Order{
		Number: string(c.Body()),
		UserID: uid,
	}
	err = h.postOrderService.SaveOrder(c.Context(), order)
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
