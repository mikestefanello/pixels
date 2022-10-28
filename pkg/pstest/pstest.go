package pstest

import (
	"reflect"
	"unsafe"

	"cloud.google.com/go/pubsub"
)

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
	return nil
}

func (a *AckStatus) OnNackWithResult() *pubsub.AckResult {
	a.OnNack()
	return nil
}

func NewMessageWithAckStatus() (*pubsub.Message, *AckStatus) {
	message := pubsub.Message{}

	// Get a reflectable value of message
	messageValue := reflect.ValueOf(message)

	// The value above is unaddressable. So construct a new and addressable message and set it with the value of the unaddressable
	addressableValue := reflect.New(messageValue.Type()).Elem()
	addressableValue.Set(messageValue)

	// Get message's doneFunc field
	doneFuncField := addressableValue.FieldByName("ackh")

	// Get the address of the field
	doneFuncFieldAddress := doneFuncField.UnsafeAddr()

	// Create a pointer based on the address
	doneFuncFieldPointer := unsafe.Pointer(doneFuncFieldAddress)

	// Create a new, exported field element that points to the original
	accessibleDoneFuncField := reflect.NewAt(doneFuncField.Type(), doneFuncFieldPointer).Elem()

	// Set the field with the alternative doneFunc
	status := &AckStatus{}
	accessibleDoneFuncField.Set(reflect.ValueOf(status))

	// Get the modified message to return
	message = addressableValue.Interface().(pubsub.Message)

	return &message, status
}
