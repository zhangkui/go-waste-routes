package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug011_RejectWeighingClearsConfirmationMetadata(t *testing.T) {
	reviewService := service.NewReviewService(service.NewWorkflowEngine(), service.NewBillingEngine(service.NewConfigService()), service.NewWeighingEngine())
	confirmer := int64(9)
	when := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	record := domain.WeighingRecord{Base: domain.Base{ID: 55}, Status: "confirmed", ConfirmedBy: &confirmer, ConfirmedAt: &when}

	if err := reviewService.RejectWeighing(&record, "recheck required"); err != nil {
		t.Fatalf("reject weighing failed: %v", err)
	}
	if record.Status != "rejected" {
		t.Fatalf("record status = %s, want rejected", record.Status)
	}
	if record.ConfirmedBy != nil || record.ConfirmedAt != nil {
		t.Fatal("reject weighing kept confirmation metadata")
	}
}
