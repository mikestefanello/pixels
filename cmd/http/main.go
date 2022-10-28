package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pixels/config"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/mikestefanello/pixels/pkg/handler"
	"github.com/mikestefanello/pixels/pkg/repo"
	"github.com/rs/zerolog/log"
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

	// Create an HTTP server
	httphandler := handler.NewEventHTTPHandler(service)
	srv := echo.New()
	srv.GET("/event", httphandler.NewEvent)

	// Start the server
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
	if err := srv.Start(addr); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("shutting down the server")
	}
}
