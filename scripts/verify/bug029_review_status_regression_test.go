package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug029_ReviewOutcomeSynchronizesWeighingStatus(t *testing.T) {
	weighing := &domain.WeighingRecord{Status: "reviewing"}
	abnormality := &domain.WeighingAbnormality{Status: "pending"}
	service.NewWeighingService().ApplyReview(weighing, abnormality, service.NewWeighingEngine(), service.NewWorkflowEngine(), 9, "approved", time.Now())
	if got, want := weighing.Status, "confirmed"; got != want {
		t.Fatalf("approved review left weighing status %q, want %q", got, want)
	}
}
