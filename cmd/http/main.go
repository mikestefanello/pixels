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

	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Initialize the logger
	log.Logger = log.With().Str("application", cfg.App).Logger()

	// Create the event repository based on the current environment
	var repository event.Repository
	switch cfg.Environment {
	case config.EnvProduction:
		client, err := pubsub.NewClient(ctx, cfg.Project)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		repository = repo.NewEventPubsubRepository(client, cfg.Topic)

	case config.EnvLocal:
		repository = repo.NewEventMemoryRepository()
	}

	// Create the event service
	service := event.NewService(repository)

	// Create an HTTP server
	httphandler := handler.NewEventHTTPHandler(service)
	srv := echo.New()
	srv.GET("/event", httphandler.New)

	// Start the server
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
	if err := srv.Start(addr); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("shutting down the server")
	}
}
