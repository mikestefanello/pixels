package handler

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/go-playground/validator/v10"
	eventv1 "github.com/mikestefanello/pixels/gen/protos/event/v1"
	"github.com/mikestefanello/pixels/gen/protos/event/v1/eventv1connect"
	"github.com/mikestefanello/pixels/pkg/event"
)

type grpcEventHandler struct {
	service event.Service
}

func NewEventGRPCHandler(service event.Service) eventv1connect.EventServiceHandler {
	return &grpcEventHandler{
		service: service,
	}
}

func (h *grpcEventHandler) New(ctx context.Context,
	req *connect.Request[eventv1.NewEventRequest]) (*connect.Response[eventv1.NewEventResponse], error) {
	// Create an event from the request
	e := &event.Event{
		CookieID:           req.Msg.CookieID,
		Country:            req.Msg.Country,
		Email:              req.Msg.Email,
		Hotel:              req.Msg.Hotel,
		ConfirmationNumber: req.Msg.ConfirmationNumber,
	}

	resp := &connect.Response[eventv1.NewEventResponse]{}

	// Attempt to save the event
	err := h.service.Insert(ctx, e)

	// Handle the error
	switch err.(type) {
	case nil:
		return resp, nil
	case validator.ValidationErrors:
		return resp, connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return resp, connect.NewError(connect.CodeInternal, err)
	}
}
