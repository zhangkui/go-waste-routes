package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug009_TaskCompleteRejectsNegativeMetrics(t *testing.T) {
	taskService := service.NewTaskService()
	task := domain.Task{Base: domain.Base{ID: 100}, Status: domain.TaskRunning}

	if err := taskService.Complete(&task, time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC), -1, 12.5); err == nil {
		t.Fatal("task completion accepted negative mileage")
	}
	if err := taskService.Complete(&task, time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC), 12.5, -1); err == nil {
		t.Fatal("task completion accepted negative fuel")
	}
}
