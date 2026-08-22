package repository

import "context"

type CRUD[T any] interface {
	Create(context.Context, T) (T, error)
	Get(context.Context, int64) (T, bool, error)
	List(context.Context, int, int) ([]T, int64, error)
	Update(context.Context, int64, T) (T, error)
	Delete(context.Context, int64) error
}
