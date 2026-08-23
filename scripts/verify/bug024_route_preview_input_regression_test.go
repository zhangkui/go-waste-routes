package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug024_RoutePreviewDoesNotMutateInput(t *testing.T) {
	stops := []domain.RouteStop{
		{Sequence: 2, CustomerID: 20},
		{Sequence: 1, CustomerID: 10},
	}
	original := append([]domain.RouteStop(nil), stops...)
	ordered := service.NewRouteService().SortStops(stops)
	if len(ordered) != 2 || ordered[0].Sequence != 1 || ordered[1].Sequence != 2 {
		t.Fatalf("preview order = %#v, want sequence [1 2]", ordered)
	}
	if stops[0] != original[0] || stops[1] != original[1] {
		t.Fatalf("input stops mutated: got %#v, want %#v", stops, original)
	}
}