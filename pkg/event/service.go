package event

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

type Service interface {
	Validate(Event) error
	Insert(context.Context, *Event) error
}

type service struct {
	repo      Repository
	validator *validator.Validate
}

func NewService(repo Repository) Service {
	return &service{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *service) Validate(e Event) error {
	return s.validator.Struct(e)
}

func (s *service) Insert(ctx context.Context, e *Event) error {
	if err := s.Validate(*e); err != nil {
		return err
	}

	e.ID = ulid.Make().String()
	e.CreatedAt = time.Now()
	return s.repo.Insert(ctx, *e)
}
