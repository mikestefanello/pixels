package repo

import (
	"context"
	"sync"

	"github.com/mikestefanello/pixels/pkg/event"
)

type eventMemoryRepo struct {
	store sync.Map
}

func NewEventMemoryRepository() event.Repository {
	return &eventMemoryRepo{store: sync.Map{}}
}

func (r *eventMemoryRepo) Insert(ctx context.Context, e event.Event) error {
	r.store.Store(e.ID, e)
	return nil
}
