package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug020_SameStatusTransitionIsIdempotent(t *testing.T) {
	engine := service.NewWorkflowEngine()
	if err := engine.ValidateTransition(service.WorkflowTask, domain.TaskRunning, domain.TaskRunning); err != nil {
		t.Fatalf("workflow engine rejected same-state transition: %v", err)
	}
	if err := service.ValidateStatusTransition(domain.TaskRunning, domain.TaskRunning, engine.FlowFor(service.WorkflowTask)); err != nil {
		t.Fatalf("shared validator rejected same-state transition: %v", err)
	}
}
