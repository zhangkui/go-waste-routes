package verify

import (
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"testing"
)

func TestBug037_DashboardExcludesInactivePendingStops(t *testing.T) {
	tasks := []domain.Task{{Status: domain.TaskPending, Stops: []domain.TaskStop{{CustomerID: 1}, {CustomerID: 2}}}}
	customers := []domain.Customer{{Base: domain.Base{ID: 1}, Status: "enabled"}, {Base: domain.Base{ID: 2}, Status: "disabled"}}
	if got := service.NewTaskService().PendingStopsForActiveCustomers(tasks, customers); got != 1 {
		t.Fatalf("pending active stops=%d,want 1", got)
	}
}
