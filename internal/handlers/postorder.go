package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/loyalsys/internal/logger"
	"github.com/rycln/loyalsys/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

type postOrderServicer interface {
	SaveOrder(context.Context, *models.Order) error
}

type postOrderJWT interface {
	ParseIDFromAuthHeader(string) (models.UserID, error)
}

type PostOrderHandler struct {
	postOrderService postOrderServicer
	jwt              postOrderJWT
}

func NewPostOrderHandler(postOrderService postOrderServicer, jwt postOrderJWT) func(*fiber.Ctx) error {
	h := &PostOrderHandler{
		postOrderService: postOrderService,
		jwt:              jwt,
	}
	return h.handle
}

type errOrderExists interface {
	error
	IsErrOrderExists() bool
}

type errWrongNum interface {
	error
	IsErrWrongNum() bool
}

type errOrderConflict interface {
	error
	IsErrOrderConflict() bool
}

func (h *PostOrderHandler) handle(c *fiber.Ctx) error {
	uid, err := h.jwt.ParseIDFromAuthHeader(c.Get("Authorization"))
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	order := &models.Order{
		Number: string(c.Body()),
		UserID: uid,
	}
	err = h.postOrderService.SaveOrder(c.Context(), order)
	if e, ok := err.(errOrderExists); ok && e.IsErrOrderExists() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusOK)
	}
	if e, ok := err.(errWrongNum); ok && e.IsErrWrongNum() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if e, ok := err.(errOrderConflict); ok && e.IsErrOrderConflict() {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusConflict)
	}
	if err != nil {
		logger.Log.Debug("path:"+c.Path(), zap.Error(err))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
