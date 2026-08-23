package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug023_EmptyRouteDoesNotCreateTask(t *testing.T) {
	svc := service.NewTaskService()
	for _, route := range []domain.Route{
		{Base: domain.Base{ID: 4}, PlanID: 9},
		{Base: domain.Base{ID: 5}, PlanID: 9, Stops: []domain.RouteStop{{CustomerID: 0, Sequence: 1}}},
	} {
		task, err := svc.InitializeFromRoute(route, time.Now())
		if err == nil || task != nil {
			t.Fatalf("invalid route created task=%#v err=%v", task, err)
		}
	}
}
