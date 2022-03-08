package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ali-a-a/gophermq/internal/app/gophermq/broker"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Handler represents broker handler.
type Handler struct {
	mq broker.Broker
}

// NewHandler creates new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Publish(ctx echo.Context) error {
	req := &PublishReq{}

	if err := ctx.Bind(req); err != nil {
		logrus.Warnf("failed to bind request: %s", err.Error())

		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "request's body is invalid"})
	}

	_, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("failed to marshal request: %s", err.Error())

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	if err = h.mq.Publish(req.Subject, []byte(req.Data)); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "publish error"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (h *Handler) PublishAsync(ctx echo.Context) error {
	req := &PublishReq{}

	if err := ctx.Bind(req); err != nil {
		logrus.Warnf("failed to bind request: %s", err.Error())

		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "request's body is invalid"})
	}

	_, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("failed to marshal request: %s", err.Error())

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	if err = h.mq.Publish(req.Subject, []byte(req.Data)); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "publish error"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
