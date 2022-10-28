package repo

import (
	"context"
	"sync"

	"github.com/mikestefanello/pixels/pkg/event"
)

type eventMemoryRepo struct {
	storage sync.Map
}

func NewEventMemoryRepository() event.Repository {
	return &eventMemoryRepo{storage: sync.Map{}}
}

func (r *eventMemoryRepo) Insert(ctx context.Context, e event.Event) error {
	r.storage.Store(e.ID, e)
	return nil
}
