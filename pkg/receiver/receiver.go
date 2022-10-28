package receiver

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/mikestefanello/pixels/pkg/compress"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/rs/zerolog/log"
)

type pubsubReceiver struct {
	service      event.Service
	decompressor compress.Decompressor
}

func NewPubsubReceiver(service event.Service) *pubsubReceiver {
	return &pubsubReceiver{
		service:      service,
		decompressor: compress.NewZlibDecompressor(),
	}
}

func (r *pubsubReceiver) Receive(ctx context.Context, msg *pubsub.Message) {
	// Create a logger for the message
	logger := log.With().
		Str("message_id", msg.ID).
		Logger()

	// Decompress
	if err := r.decompressor.DecompressMessage(msg); err != nil {
		logger.Error().
			Err(err).
			Msg("could not decompress pubsub message")

		msg.Nack()
		return
	}

	// Unmarshal
	var e event.Event
	if err := json.Unmarshal(msg.Data, &e); err != nil {
		log.Error().
			Err(err).
			Msg("could not unmarshal pubsub message")

		msg.Nack()
		return
	}

	// Store the message
	err := r.service.Insert(ctx, &e)

	// Handle the error
	switch err.(type) {
	case validator.ValidationErrors:
		log.Debug().
			Err(err).
			Msg("invalid pubsub message data received")
		msg.Ack()

	case nil:
		msg.Ack()

	default:
		log.Error().
			Err(err).
			Msg("unable to store pubsub message")
		msg.Nack()
	}
}
