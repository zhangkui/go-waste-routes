package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug003_TaskSkipRejectsUnknownReason(t *testing.T) {
	taskService := service.NewTaskService()
	stop := domain.TaskStop{Base: domain.Base{ID: 88}}

	if err := taskService.Skip(&stop, "external_delay", "need approval"); err == nil {
		t.Fatal("task skip accepted an unsupported reason")
	}
}
