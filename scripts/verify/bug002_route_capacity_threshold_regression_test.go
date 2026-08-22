package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug002_RouteCapacityAllowsEightyFivePercentThreshold(t *testing.T) {
	routeService := service.NewRouteService()
	route := domain.Route{Base: domain.Base{ID: 11}, EstimatedWeightTons: 8.4}
	vehicle := domain.Vehicle{Base: domain.Base{ID: 21}, RatedLoadTons: 10}
	stops := []domain.RouteStop{
		{Base: domain.Base{ID: 31}, Sequence: 1, EstimatedWeightTons: 4.2},
		{Base: domain.Base{ID: 32}, Sequence: 2, EstimatedWeightTons: 4.2},
	}

	if err := routeService.ValidateCapacity(route, vehicle, stops); err != nil {
		t.Fatalf("route service rejected 84%% load: %v", err)
	}
	if err := service.ValidateRouteCapacity(route.EstimatedWeightTons, vehicle.RatedLoadTons); err != nil {
		t.Fatalf("shared validator rejected 84%% load: %v", err)
	}
}
