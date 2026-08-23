package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"go-waste-routes/internal/transport/http/handler"
)

func TestBug027_InvalidRouteStopsAreRejected(t *testing.T) {
	h := handler.NewRouteCapacityHandler(service.NewRouteService())
	vehicle := domain.Vehicle{RatedLoadTons: 10}

	if err := h.Validate(domain.Route{}, vehicle, nil); err == nil {
		t.Error("empty route was accepted")
	}
	if err := h.Validate(domain.Route{}, vehicle, []domain.RouteStop{{CustomerID: 10, Sequence: 1, EstimatedWeightTons: -1}}); err == nil {
		t.Error("negative estimated weight was accepted")
	}
	if err := h.Validate(domain.Route{}, vehicle, []domain.RouteStop{{CustomerID: 10, Sequence: 1, EstimatedWeightTons: 4}}); err != nil {
		t.Fatalf("valid route rejected: %v", err)
	}
}
