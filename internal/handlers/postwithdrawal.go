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

type postWithdrawalServicer interface {
	WithdrawalProcessing(context.Context, *models.Withdrawal) error
}

type postWithdrawalJWT interface {
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type PostWithdrawalHandler struct {
	postWithdrawalService postWithdrawalServicer
	jwt                   postWithdrawalJWT
}

func NewPostWithdrawalHandler(postWithdrawalService postWithdrawalServicer, jwt postWithdrawalJWT) func(*fiber.Ctx) error {
	h := &PostWithdrawalHandler{
		postWithdrawalService: postWithdrawalService,
		jwt:                   jwt,
	}
	return h.handle
}

type errNotEnoughCurrency interface {
	error
	IsErrNotEnoughCurrency() bool
}

type errWrongOrderNum interface {
	error
	IsErrWrongOrderNum() bool
}

func (h *PostWithdrawalHandler) handle(c *fiber.Ctx) error {
	uid, err := h.jwt.ParseIDFromAuthHeader(c.Get("Authorization"))
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var withdrawal models.Withdrawal
	err = json.Unmarshal(c.Body(), &withdrawal)
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusBadRequest)
	}
	withdrawal.UserID = uid

	err = h.postWithdrawalService.WithdrawalProcessing(c.Context(), &withdrawal)
	if e, ok := err.(errNotEnoughCurrency); ok && e.IsErrNotEnoughCurrency() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusPaymentRequired)
	}
	if e, ok := err.(errWrongOrderNum); ok && e.IsErrWrongOrderNum() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
