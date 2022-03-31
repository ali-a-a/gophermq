package handler

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, broker.ErrMaxPending) {
			return ctx.JSON(http.StatusTooManyRequests, echo.Map{"message": err.Error()})
		}

		if errors.Is(err, broker.ErrSubscriberNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		}

		if errors.Is(err, broker.ErrBadSubject) {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
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

func (h *Handler) Subscribe(ctx echo.Context) error {
	req := &SubscribeReq{}

	if err := ctx.Bind(req); err != nil {
		logrus.Warnf("failed to bind request: %s", err.Error())

		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "request's body is invalid"})
	}

	_, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("failed to marshal request: %s", err.Error())

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	sub, err := h.mq.Subscribe(req.Subject)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to subscribe"})
	}

	res := &SubscribeRes{
		Subject: sub.Subj,
		ID:      sub.ID,
	}

	return ctx.JSON(http.StatusOK, echo.Map{"subject": res.Subject, "id": res.ID})
}

func (h *Handler) Fetch(ctx echo.Context) error {
	req := &FetchReq{}

	if err := ctx.Bind(req); err != nil {
		logrus.Warnf("failed to bind request: %s", err.Error())

		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "request's body is invalid"})
	}

	_, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("failed to marshal request: %s", err.Error())

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	data, err := h.mq.Fetch(req.Subject, req.ID)
	if err != nil {
		if errors.Is(err, broker.ErrBadID) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": "server error"})
	}

	finalData := make([]string, len(data))

	for i := range data {
		finalData[i] = string(data[i])
	}

	res := &FetchRes{
		Subject: req.Subject,
		ID:      req.ID,
		Data:    finalData,
	}

	return ctx.JSON(http.StatusOK, echo.Map{"subject": res.Subject, "id": res.ID, "data": res.Data})
}
