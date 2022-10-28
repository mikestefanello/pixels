package main

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/mikestefanello/pixels/config"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/mikestefanello/pixels/pkg/receiver"
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

	// Start receiving from pubsub
	rec := receiver.NewPubsubReceiver(service)
	err = client.Subscription(cfg.Subscription).
		Receive(ctx, rec.Receive)

	if err != nil {
		log.Error().Err(err).Send()
	}
}
