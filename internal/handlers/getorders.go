package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rycln/loyalsys/internal/auth"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"github.com/rycln/loyalsys/internal/storage"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type getOrderServicer interface {
	GetUserOrders(context.Context, models.UserID) ([]*models.OrderDB, error)
}

type GetOrderHandler struct {
	getOrderService getOrderServicer
}

func NewGetOrderHandler(getOrderService getOrderServicer) func(*fiber.Ctx) error {
	h := &GetOrderHandler{
		getOrderService: getOrderService,
	}
	return h.handle
}

func (h *GetOrderHandler) handle(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	uid, err := auth.ParseIDFromJWT(token)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	orders, err := h.getOrderService.GetUserOrders(c.Context(), uid)
	if errors.Is(err, storage.ErrNoOrder) {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	resBody, err := json.Marshal(&orders)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(resBody)
}
