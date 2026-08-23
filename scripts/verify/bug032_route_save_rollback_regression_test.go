package verify

import (
	"context"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"testing"
)

type bug032Writer struct{}

func (bug032Writer) SaveStops(context.Context, int64, []domain.RouteStop) error { return nil }
func TestBug032_RouteSaveFailureDoesNotKeepPartialStops(t *testing.T) {
	route := domain.Route{Stops: []domain.RouteStop{{CustomerID: 1, Sequence: 1}}}
	err := service.NewRouteService().SaveRouteStops(context.Background(), &route, []domain.RouteStop{{CustomerID: 2, Sequence: 1}, {CustomerID: 3, Sequence: 3}}, bug032Writer{})
	if err == nil {
		t.Fatal("invalid stops saved")
	}
	if len(route.Stops) != 1 {
		t.Fatalf("route retained %d stops after failed save, want 1", len(route.Stops))
	}
}
