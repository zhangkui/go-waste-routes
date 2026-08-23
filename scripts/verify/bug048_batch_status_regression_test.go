package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug048_BatchStatusUpdateKeepsAllCustomerStops(t *testing.T) {
	stops := []domain.TaskStop{{CustomerID: 8, Status: "pending"}, {CustomerID: 8, Status: "pending"}, {CustomerID: 9, Status: "pending"}}
	updated := service.NewTaskService().BatchCompleteCustomerStops(stops)
	if pending := service.NewReportService().PendingStops(updated); pending != 0 {
		t.Fatalf("batch update left %d pending stops", pending)
	}
}
