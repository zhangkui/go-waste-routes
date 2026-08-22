package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug012_ApprovedAbnormalityClearsRejectReason(t *testing.T) {
	reviewService := service.NewReviewService(service.NewWorkflowEngine(), service.NewBillingEngine(service.NewConfigService()), service.NewWeighingEngine())
	abnormality := domain.WeighingAbnormality{Base: domain.Base{ID: 66}, Status: "pending"}

	if err := reviewService.ReviewAbnormality(&abnormality, service.ReviewDecision{
		Approved:   true,
		Reason:     "looks fine",
		ReviewerID: 7,
		ReviewedAt: time.Date(2026, 8, 22, 11, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("review abnormality failed: %v", err)
	}
	if abnormality.Status != "confirmed" || abnormality.ReviewResult != "confirmed" {
		t.Fatalf("approved abnormality not confirmed: %#v", abnormality)
	}
	if abnormality.RejectReason != "" {
		t.Fatalf("approved abnormality kept reject reason: %q", abnormality.RejectReason)
	}
}
