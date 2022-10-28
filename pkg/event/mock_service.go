package event

import (
	"context"
)

type MockService struct {
	Events []Event
}

func NewMockService() *MockService {
	return &MockService{
		Events: make([]Event, 0),
	}
}

func (m *MockService) Validate(e Event) error {
	return nil
}

func (m *MockService) Insert(ctx context.Context, e *Event) error {
	m.Events = append(m.Events, *e)
	return nil
}
