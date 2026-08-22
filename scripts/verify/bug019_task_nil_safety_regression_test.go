package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/service"
)

func TestBug019_TaskServiceHandlesNilTasksGracefully(t *testing.T) {
	taskService := service.NewTaskService()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("task service panicked: %v", recovered)
		}
	}()

	if err := taskService.Claim(nil, time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("claim on nil task returned nil error")
	}
	if err := taskService.Complete(nil, time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC), 1, 1); err == nil {
		t.Fatal("complete on nil task returned nil error")
	}
}
