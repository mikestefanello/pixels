package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/mikestefanello/pixels/config"
	"github.com/mikestefanello/pixels/gen/protos/event/v1/eventv1connect"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/mikestefanello/pixels/pkg/handler"
	"github.com/mikestefanello/pixels/pkg/repo"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	ctx := context.Background()

	// Initialize the logger
	log.Logger = log.With().Str("application", "pixels").Logger()

	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Connect to pubsub
	client, err := pubsub.NewClient(ctx, cfg.Project)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Create the event service
	repository := repo.NewEventPubsubRepo(client, cfg.Topic)
	service := event.NewService(repository)

	// Start a Connect server
	err = http.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port),
		buildHandler(service),
	)

	if err != nil {
		log.Error().Err(err).Msg("server terminated")
	}
}

func buildHandler(service event.Service) http.Handler {
	grpchandler := handler.NewEventGRPCHandler(service)
	path, hdlr := eventv1connect.NewEventServiceHandler(grpchandler)

	mux := http.NewServeMux()
	mux.Handle(path, hdlr)
	mux.Handle(grpcreflect.NewHandlerV1(
		grpcreflect.NewStaticReflector(eventv1connect.EventServiceName),
		connect.WithCompressMinBytes(1024),
	))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(eventv1connect.EventServiceName),
		connect.WithCompressMinBytes(1024),
	))
	mux.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(eventv1connect.EventServiceName),
		connect.WithCompressMinBytes(1024),
	))

	return h2c.NewHandler(mux, &http2.Server{})
}
