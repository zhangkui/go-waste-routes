package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/service"
)

func TestBug004_TaskNumberUsesTaskPrefix(t *testing.T) {
	got := service.NewTaskService().BuildTaskNumber(time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC), 1)
	if got != "TASK-20260822-0001" {
		t.Fatalf("task number = %s, want TASK-20260822-0001", got)
	}
}
