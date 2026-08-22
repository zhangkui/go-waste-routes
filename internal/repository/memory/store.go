package memory

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"time"
)

type Store[T any] struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]T
}

func New[T any]() *Store[T] {
	return &Store[T]{items: make(map[int64]T)}
}

func (s *Store[T]) Create(_ context.Context, value T) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	setID(&value, s.nextID)
	setTimestamps(&value, true)
	s.items[s.nextID] = value
	return value, nil
}

func (s *Store[T]) Get(_ context.Context, id int64) (T, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.items[id]
	return value, ok, nil
}

func (s *Store[T]) List(_ context.Context, page, pageSize int) ([]T, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]int64, 0, len(s.items))
	for id := range s.items {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	end := start + pageSize
	if end > len(ids) {
		end = len(ids)
	}
	items := make([]T, 0, end-start)
	for _, id := range ids[start:end] {
		items = append(items, s.items[id])
	}
	return items, int64(len(ids)), nil
}

func (s *Store[T]) Update(_ context.Context, id int64, value T) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return value, context.Canceled
	}
	setID(&value, id)
	setTimestamps(&value, false)
	s.items[id] = value
	return value, nil
}

func (s *Store[T]) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, id)
	return nil
}

func (s *Store[T]) All() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]T, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items
}

func setID[T any](value *T, id int64) {
	rv := reflect.ValueOf(value).Elem()
	base := rv.FieldByName("Base")
	if base.IsValid() {
		field := base.FieldByName("ID")
		if field.IsValid() && field.CanSet() {
			field.SetInt(id)
		}
	}
}

func setTimestamps[T any](value *T, create bool) {
	rv := reflect.ValueOf(value).Elem()
	base := rv.FieldByName("Base")
	if !base.IsValid() {
		return
	}
	now := time.Now().UTC()
	if create {
		field := base.FieldByName("CreatedAt")
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(now))
		}
	}
	updatedField := base.FieldByName("UpdatedAt")
	if updatedField.IsValid() && updatedField.CanSet() {
		updatedField.Set(reflect.ValueOf(now))
	}
}
