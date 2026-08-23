package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug039_SummaryFailureDoesNotCompleteTask(t *testing.T) {
	task := &domain.Task{Status: domain.TaskRunning}
	stops := []domain.TaskStop{{Status: "completed"}, {Status: "pending"}}
	err := service.NewTaskService().CompleteAfterSummary(task, stops, time.Now(), 12, 3)
	if err == nil {
		t.Fatal("unfinished stop summary must prevent task completion")
	}
	if task.Status != domain.TaskRunning || task.CompletedAt != nil {
		t.Fatalf("task state changed after rejected summary: %#v", task)
	}
}
