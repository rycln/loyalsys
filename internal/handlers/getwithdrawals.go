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

type getWithdrawalsServicer interface {
	GetUserWithdrawals(context.Context, models.UserID) ([]*models.Withdrawal, error)
}

type getWithdrawalsJWT interface {
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type GetWithdrawalsHandler struct {
	getWithdrawalService getWithdrawalsServicer
	jwt                  getWithdrawalsJWT
}

func NewGetWithdrawalsHandler(getWithdrawalService getWithdrawalsServicer, jwt getWithdrawalsJWT) func(*fiber.Ctx) error {
	h := &GetWithdrawalsHandler{
		getWithdrawalService: getWithdrawalService,
		jwt:                  jwt,
	}
	return h.handle
}

type errNoWithdrawal interface {
	error
	IsErrNoWithdrawal() bool
}

func (h *GetWithdrawalsHandler) handle(c *fiber.Ctx) error {
	uid, err := h.jwt.ParseIDFromAuthHeader(c.Get("Authorization"))
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	withdrawals, err := h.getWithdrawalService.GetUserWithdrawals(c.Context(), uid)
	if e, ok := err.(errNoWithdrawal); ok && e.IsErrNoWithdrawal() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	resBody, err := json.Marshal(&withdrawals)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		c.SendStatus(fiber.StatusInternalServerError)
	}
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusOK).Send(resBody)
}
