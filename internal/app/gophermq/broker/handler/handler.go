package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ali-a-a/gophermq/internal/app/gophermq/broker"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	// workerSize represents the number of workers in the pool
	// that be used in async mode publish.
	workerSize = 100
)

// Handler represents broker handler.
type Handler struct {
	mq         broker.Broker
	pubReqChan chan *PublishReq
}

// NewHandler creates new Handler.
func NewHandler(mq broker.Broker) *Handler {
	h := &Handler{
		mq:         mq,
		pubReqChan: make(chan *PublishReq),
	}

	h.startPublishChan(workerSize)

	return h
}

func (h *Handler) startPublishChan(size int) {
	for i := 0; i < size; i++ {
		go h.publishAsync()
	}
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
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (h *Handler) publishAsync() {
	for event := range h.pubReqChan {
		if err := h.mq.Publish(event.Subject, []byte(event.Data)); err != nil {
			continue
		}
	}
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

	h.pubReqChan <- req

	return ctx.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
