package event

import (
	"context"
)

type MockService struct {
	Events []Event
	Errors struct {
		Validate error
		Insert   error
	}
}

func NewMockService() *MockService {
	return &MockService{
		Events: make([]Event, 0),
	}
}

func (m *MockService) Validate(e Event) error {
	return m.Errors.Validate
}

func (m *MockService) Insert(ctx context.Context, e *Event) error {
	if m.Errors.Insert != nil {
		return m.Errors.Insert
	}
	m.Events = append(m.Events, *e)
	return nil
}
