package repo

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/mikestefanello/pixels/pkg/compress"
	"github.com/mikestefanello/pixels/pkg/event"
)

type eventPubsubRepo struct {
	client     *pubsub.Client
	topic      string
	compressor compress.Compressor
}

func NewEventPubsubRepo(client *pubsub.Client, topic string) event.Repository {
	return &eventPubsubRepo{
		client:     client,
		topic:      topic,
		compressor: compress.NewZlibCompressor(),
	}
}

func (r *eventPubsubRepo) Insert(ctx context.Context, e event.Event) error {
	// Marshal to JSON
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// Compress the data
	data, err = r.compressor.Compress(data)
	if err != nil {
		return err
	}

	// Publish to the topic
	result := r.client.
		Topic(r.topic).
		Publish(ctx, &pubsub.Message{
			Data: data,
		})

	if _, err := result.Get(ctx); err != nil {
		return err
	}

	return nil
}
