package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug043_CancelledRouteBlocksTaskExecution(t *testing.T) {
	task := &domain.Task{Status: domain.TaskClaimed}
	route := &domain.Route{Status: domain.RouteCancelled}
	err := service.NewTaskService().StartForRoute(task, route, time.Now())
	if err == nil {
		t.Fatal("cancelled route must block task execution")
	}
	if task.Status != domain.TaskClaimed {
		t.Fatalf("task started after route cancellation: %s", task.Status)
	}
}
