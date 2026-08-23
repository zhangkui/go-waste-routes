package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug042_DepartureRejectsStaleTaskStops(t *testing.T) {
	route := &domain.Route{Status: domain.RoutePending}
	task := domain.Task{Base: domain.Base{ID: 12}}
	err := service.NewRouteService().StartDeparture(route, task, []domain.TaskStop{{TaskID: 99, Status: "pending"}})
	if err == nil {
		t.Fatal("stale stops must block departure")
	}
	if route.Status != domain.RoutePending {
		t.Fatalf("route started with stale stops: %s", route.Status)
	}
}
