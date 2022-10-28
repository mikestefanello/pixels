package pstest

import (
	"reflect"
	"unsafe"

	"cloud.google.com/go/pubsub"
)

// This _hack_ is required to forcefully inject our own ack/nack handler on a pubsub message,
// which is required for testing. This will be something shared outside of the application so others can
// leverage this.

type AckStatus struct {
	acked  bool
	nacked bool
}

func (a *AckStatus) IsAcked() bool {
	return a.acked
}

func (a *AckStatus) IsNacked() bool {
	return a.nacked
}

func (a *AckStatus) OnAck() {
	a.acked = true
}

func (a *AckStatus) OnNack() {
	a.nacked = true
}

func (a *AckStatus) OnAckWithResult() *pubsub.AckResult {
	a.OnNack()
	return &pubsub.AckResult{}
}

func (a *AckStatus) OnNackWithResult() *pubsub.AckResult {
	a.OnNack()
	return &pubsub.AckResult{}
}

func NewMessageWithAckStatus() (*pubsub.Message, *AckStatus) {
	message := pubsub.Message{}

	// Get a reflectable value of message
	messageValue := reflect.ValueOf(message)

	// The value above is unaddressable. So construct a new and addressable message and set it with the value of the unaddressable
	addressableValue := reflect.New(messageValue.Type()).Elem()
	addressableValue.Set(messageValue)

	// Get message's private ackh field
	ackhField := addressableValue.FieldByName("ackh")

	// Get the address of the field
	ackhFieldAddress := ackhField.UnsafeAddr()

	// Create a pointer based on the address
	ackhFieldPointer := unsafe.Pointer(ackhFieldAddress)

	// Create a new, exported field element that points to the original
	accessibleAckhField := reflect.NewAt(ackhField.Type(), ackhFieldPointer).Elem()

	// Set the field with the alternative ackh
	status := &AckStatus{}
	accessibleAckhField.Set(reflect.ValueOf(status))

	// Get the modified message to return
	message = addressableValue.Interface().(pubsub.Message)

	return &message, status
}
