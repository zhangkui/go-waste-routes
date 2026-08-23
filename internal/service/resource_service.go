package service

import (
	"context"
	"errors"
	"fmt"

	"go-waste-routes/internal/repository/memory"
)

var ErrResourceNotFound = errors.New("resource not found")

type ResourceService[T any] struct {
	store *memory.Store[T]
}

func NewResourceService[T any](store *memory.Store[T]) *ResourceService[T] {
	return &ResourceService[T]{store: store}
}

func (s *ResourceService[T]) Create(ctx context.Context, value T) (T, error) {
	return s.store.Create(ctx, value)
}

func (s *ResourceService[T]) Get(ctx context.Context, id int64) (T, error) {
	value, ok, err := s.store.Get(ctx, id)
	if err != nil {
		var zero T
		return zero, err
	}
	if !ok {
		var zero T
		return zero, fmt.Errorf("%w: %d", ErrResourceNotFound, id)
	}
	return value, nil
}

func (s *ResourceService[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	return s.store.List(ctx, page, pageSize)
}

func (s *ResourceService[T]) Update(ctx context.Context, id int64, value T) (T, error) {
	return s.store.Update(ctx, id, value)
}

func (s *ResourceService[T]) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func ErrNotFound(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}
