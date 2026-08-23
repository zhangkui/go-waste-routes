package verify

import (
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"testing"
)

func TestBug038_BatchTaskSummaryReleasesPerItem(t *testing.T) {
	inUse, max := 0, 0
	acquire := func(domain.Task) func() {
		inUse++
		if inUse > max {
			max = inUse
		}
		return func() { inUse-- }
	}
	tasks := []domain.Task{{Status: domain.TaskCompleted}, {Status: domain.TaskCompleted}, {Status: domain.TaskCompleted}}
	service.NewTaskService().SummarizeBatch(tasks, acquire)
	if max != 1 {
		t.Fatalf("maximum simultaneous batch resources=%d,want 1", max)
	}
}
