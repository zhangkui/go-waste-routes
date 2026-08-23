package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug041_TaskClaimRequiresRouteAssociation(t *testing.T) {
	task := &domain.Task{Status: domain.TaskPending}
	err := service.NewTaskService().ClaimForRoute(task, nil, time.Now())
	if err == nil {
		t.Fatal("task without dispatchable route must not be claimed")
	}
	if task.Status != domain.TaskPending || task.ClaimedAt != nil {
		t.Fatalf("route-less task was claimed: %#v", task)
	}
}
