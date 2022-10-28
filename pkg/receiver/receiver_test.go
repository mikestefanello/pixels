package receiver

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/go-playground/validator/v10"
	"github.com/mikestefanello/pixels/pkg/compress"
	"github.com/mikestefanello/pixels/pkg/event"
	"github.com/mikestefanello/pixels/pkg/pstest"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testEvent = event.Event{
	CookieID:           "a",
	Country:            "us",
	Email:              "m@m.com",
	Hotel:              "b",
	ConfirmationNumber: "c",
	ExtraField:         "d",
}

func TestPubsubReceiver_Receive(t *testing.T) {
	srv := event.NewMockService()
	rcv := NewPubsubReceiver(srv)

	// Pass the message to the receiver for processing
	msg, status := generateMessage(testEvent, t)
	rcv.Receive(context.Background(), msg)

	// Ensure the event made it to the service and was acked
	require.Len(t, srv.Events, 1)
	assert.Equal(t, testEvent, srv.Events[0])
	assert.True(t, status.IsAcked())
}

func TestPubsubReceiver_Receive_Failure_InvalidData(t *testing.T) {
	srv := event.NewMockService()
	rcv := NewPubsubReceiver(srv)

	// Pass the message to the receiver for processing
	msg, status := pstest.NewMessageWithAckStatus()
	msg.Data = []byte{1, 2, 3, 4, 5}
	rcv.Receive(context.Background(), msg)

	// No event should be added since this was invalid but we ack because it not processable
	require.Len(t, srv.Events, 0)
	assert.True(t, status.IsAcked())
}

func TestPubsubReceiver_Receive_Failure_Validation(t *testing.T) {
	srv := event.NewMockService()
	rcv := NewPubsubReceiver(srv)

	// Tell the mock service to fail validation
	srv.Errors.Validate = validator.ValidationErrors{}

	// Pass the message to the receiver for processing
	msg, status := generateMessage(testEvent, t)
	rcv.Receive(context.Background(), msg)

	// No event should be added since this was invalid but we ack because it not processable
	require.Len(t, srv.Events, 0)
	assert.True(t, status.IsAcked())
}

func TestPubsubReceiver_Receive_Failure_ServiceError(t *testing.T) {
	srv := event.NewMockService()
	rcv := NewPubsubReceiver(srv)

	// Tell the mock service to fail insertion
	srv.Errors.Insert = errors.New("fail")

	// Pass the message to the receiver for processing
	msg, status := generateMessage(testEvent, t)
	rcv.Receive(context.Background(), msg)

	// No event should be added since this operation failed and we nack to retry
	require.Len(t, srv.Events, 1)
	assert.True(t, status.IsNacked())
}

func generateMessage(e event.Event, t *testing.T) (*pubsub.Message, *pstest.AckStatus) {
	data, err := json.Marshal(e)
	require.NoError(t, err)

	msg, status := pstest.NewMessageWithAckStatus()
	msg.ID = ulid.Make().String()
	msg.Attributes = make(map[string]string)
	msg.Data = data

	cmp := compress.NewZlibCompressor()
	err = cmp.CompressMessage(msg)
	require.NoError(t, err)

	return msg, status
}
