package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug026_InvalidStayDoesNotMoveArrivalBackward(t *testing.T) {
	stops := []domain.RouteStop{{Sequence: 1, StayMinutes: -6}, {Sequence: 2, StayMinutes: 4}}
	arrivals := service.NewRouteService().EstimateArrivalTimes(time.Date(2026, 8, 23, 8, 0, 0, 0, time.UTC), stops)
	if arrivals[0].EstimatedArrival != "08:00" || arrivals[1].EstimatedArrival != "08:01" {
		t.Fatalf("arrival times moved backward: %#v", arrivals)
	}
}
