package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug020_SameStatusTransitionIsIdempotent(t *testing.T) {
	engine := service.NewWorkflowEngine()
	taskService := service.NewTaskService()

	claimedAt := time.Date(2026, 8, 22, 9, 10, 0, 0, time.UTC)
	runningAt := claimedAt.Add(8 * time.Minute)
	retryAt := runningAt.Add(2 * time.Minute)

	task := domain.Task{Base: domain.Base{ID: 2048}, Status: domain.TaskPending}
	if err := taskService.Claim(&task, claimedAt); err != nil {
		t.Fatalf("claim task failed: %v", err)
	}
	if task.Status != domain.TaskClaimed {
		t.Fatalf("task status after claim = %s, want %s", task.Status, domain.TaskClaimed)
	}

	if err := engine.ApplyTaskStatus(&task, domain.TaskRunning, runningAt); err != nil {
		t.Fatalf("move task into running failed: %v", err)
	}
	if task.Status != domain.TaskRunning {
		t.Fatalf("task status after run = %s, want %s", task.Status, domain.TaskRunning)
	}

	if err := engine.ValidateTransition(service.WorkflowTask, domain.TaskRunning, domain.TaskRunning); err != nil {
		t.Fatalf("workflow engine rejected same-state transition: %v", err)
	}
	if err := service.ValidateStatusTransition(domain.TaskRunning, domain.TaskRunning, engine.FlowFor(service.WorkflowTask)); err != nil {
		t.Fatalf("shared validator rejected same-state transition: %v", err)
	}

	if err := engine.ApplyTaskStatus(&task, domain.TaskRunning, retryAt); err != nil {
		t.Fatalf("retrying same task state should be idempotent: %v", err)
	}
	if task.Status != domain.TaskRunning {
		t.Fatalf("task status after retry = %s, want %s", task.Status, domain.TaskRunning)
	}
	if task.ClaimedAt == nil || !task.ClaimedAt.Equal(claimedAt) {
		t.Fatalf("claimed time changed after retry: %v", task.ClaimedAt)
	}
}
