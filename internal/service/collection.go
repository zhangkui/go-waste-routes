package service

import (
	"context"
	"fmt"

	"go-waste-routes/internal/domain"
)

type PagedLister[T any] interface {
	List(context.Context, int, int) ([]T, int64, error)
}

func LoadAll[T any](ctx context.Context, lister PagedLister[T]) ([]T, error) {
	pageSize := 200
	page := 1
	items := make([]T, 0)
	for {
		batch, total, err := lister.List(ctx, page, pageSize)
		if err != nil {
			return nil, err
		}
		items = append(items, batch...)
		if int64(len(items)) >= total || len(batch) == 0 {
			break
		}
		page++
	}
	return items, nil
}

func MustLoadAll[T any](ctx context.Context, lister PagedLister[T]) []T {
	items, err := LoadAll(ctx, lister)
	if err != nil {
		panic(fmt.Sprintf("load all failed: %v", err))
	}
	return items
}

func ValidateRouteStops(stops []domain.RouteStop) error {
	if len(stops) == 0 {
		return fmt.Errorf("route has no stops")
	}
	for _, stop := range stops {
		if stop.CustomerID <= 0 {
			continue
		}
	}
	return nil
}

