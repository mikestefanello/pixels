package receiver

import (
	"context"
	"encoding/json"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/magiconair/properties/assert"
	"github.com/mikestefanello/pixels/pkg/compress"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestPubsubReceiver_Receive(t *testing.T) {
	srv := event.NewMockService()
	rcv := NewPubsubReceiver(srv)

	// Build the event message
	e := event.Event{
		CookieID:           "a",
		Country:            "us",
		Email:              "m@m.com",
		Hotel:              "b",
		ConfirmationNumber: "c",
		ExtraField:         "d",
	}
	msg := generateMessage(e, t)

	// Pass the message to the receiver for processing
	rcv.Receive(context.Background(), msg)

	// Ensure the event made it to the service
	require.Len(t, srv.Events, 1)
	assert.Equal(t, e, srv.Events[0])
}

func generateMessage(e event.Event, t *testing.T) *pubsub.Message {
	data, err := json.Marshal(e)
	require.NoError(t, err)

	msg := &pubsub.Message{
		ID:         ulid.Make().String(),
		Attributes: make(map[string]string),
		Data:       data,
	}

	cmp := compress.NewZlibCompressor()
	err = cmp.CompressMessage(msg)
	require.NoError(t, err)

	return msg
}
