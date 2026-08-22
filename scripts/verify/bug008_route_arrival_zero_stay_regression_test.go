package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug008_RouteArrivalRespectsZeroMinuteStops(t *testing.T) {
	routeService := service.NewRouteService()
	start := time.Date(2026, 8, 22, 8, 0, 0, 0, time.UTC)
	stops := []domain.RouteStop{
		{Base: domain.Base{ID: 2}, Sequence: 2, StayMinutes: 15},
		{Base: domain.Base{ID: 1}, Sequence: 1, StayMinutes: 0},
	}

	ordered := routeService.EstimateArrivalTimes(start, stops)
	if len(ordered) != 2 {
		t.Fatalf("stops length = %d, want 2", len(ordered))
	}
	if ordered[0].ID != 1 || ordered[1].ID != 2 {
		t.Fatalf("stops not sorted by sequence: %#v", ordered)
	}
	if ordered[1].EstimatedArrival != "08:00" {
		t.Fatalf("second stop arrival = %s, want 08:00", ordered[1].EstimatedArrival)
	}
}
