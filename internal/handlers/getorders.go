package handlers

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type getOrdersServicer interface {
	GetUserOrders(context.Context, models.UserID) ([]*models.OrderDB, error)
}

type getOrdersJWT interface {
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type GetOrdersHandler struct {
	getOrderService getOrdersServicer
	jwt             getOrdersJWT
}

func NewGetOrdersHandler(getOrderService getOrdersServicer, jwt getOrdersJWT) func(*fiber.Ctx) error {
	h := &GetOrdersHandler{
		getOrderService: getOrderService,
		jwt:             jwt,
	}
	return h.handle
}

type errNoOrder interface {
	error
	IsErrNoOrder() bool
}

func (h *GetOrdersHandler) handle(c *fiber.Ctx) error {
	uid, err := h.jwt.ParseIDFromAuthHeader(c.Get("Authorization"))
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	orders, err := h.getOrderService.GetUserOrders(c.Context(), uid)
	if e, ok := err.(errNoOrder); ok && e.IsErrNoOrder() {
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
