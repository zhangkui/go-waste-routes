package verify

import (
	"context"
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
)

func TestBug035_MissingResourceIsClassifiedAsNotFound(t *testing.T) {
	store := memory.New[domain.Vehicle]()
	_, err := service.NewResourceService(store).Get(context.Background(), 404)
	if err == nil || !service.ErrNotFound(err) {
		t.Fatalf("missing resource error was not classified as not found: %v", err)
	}
}
