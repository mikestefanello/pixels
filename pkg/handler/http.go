package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/rs/zerolog/log"
)

var trackingPixel = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

type EventHTTPHandler interface {
	New(ctx echo.Context) error
}

type eventHTTPHandler struct {
	service event.Service
}

func NewEventHTTPHandler(service event.Service) EventHTTPHandler {
	return &eventHTTPHandler{
		service: service,
	}
}

func (h *eventHTTPHandler) New(ctx echo.Context) error {
	var e event.Event

	// Get the event from the query parameters
	if err := ctx.Bind(&e); err != nil {
		return err
	}

	defer func() {
		// Store the event
		if err := h.service.Insert(ctx.Request().Context(), &e); err != nil {
			log.Error().Err(err).Msg("unable to save HTTP event")
		}
	}()

	return ctx.Blob(http.StatusOK, "image/gif", trackingPixel)
}
