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

type getBalanceServicer interface {
	GetUserBalance(context.Context, models.UserID) (*models.Balance, error)
}

type getBalanceJWT interface {
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type GetBalanceHandler struct {
	getBalanceService getBalanceServicer
	jwt               getBalanceJWT
}

func NewGetBalanceHandler(getBalanceService getBalanceServicer, jwt getBalanceJWT) func(*fiber.Ctx) error {
	h := &GetBalanceHandler{
		getBalanceService: getBalanceService,
		jwt:               jwt,
	}
	return h.handle
}

func (h *GetBalanceHandler) handle(c *fiber.Ctx) error {
	uid, err := h.jwt.ParseIDFromAuthHeader(c.Get("Authorization"))
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	balance, err := h.getBalanceService.GetUserBalance(c.Context(), uid)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	resBody, err := json.Marshal(&balance)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(resBody)
}
