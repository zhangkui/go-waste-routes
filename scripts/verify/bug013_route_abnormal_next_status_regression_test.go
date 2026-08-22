package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug013_RouteAbnormalAdvancesToCancelled(t *testing.T) {
	routeService := service.NewRouteService()
	if got := routeService.NextStatus(domain.RouteAbnormal); got != domain.RouteCancelled {
		t.Fatalf("next status = %s, want %s", got, domain.RouteCancelled)
	}
}
