package verify

import (
	"errors"
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug049_RouteAbortReleasesVehicleState(t *testing.T) {
	route := &domain.Route{Status: domain.RoutePending}
	vehicle := &domain.Vehicle{Status: "available"}
	err := service.NewRouteService().ExecuteRoute(route, vehicle, func() error { return errors.New("blocked access") })
	if err == nil {
		t.Fatal("route execution error must be returned")
	}
	if route.Status != domain.RouteAbnormal || vehicle.Status != "available" {
		t.Fatalf("route abort did not release vehicle: route=%s vehicle=%s", route.Status, vehicle.Status)
	}
}
