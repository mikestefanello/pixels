package event

import (
	"context"
)

type Repository interface {
	Insert(context.Context, Event) error
}
