package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug046_CorrectionPreservesConfirmedHistory(t *testing.T) {
	record := &domain.WeighingRecord{Base: domain.Base{ID: 71}, GrossWeight: 10, TareWeight: 3, NetWeight: 7, Status: "confirmed"}
	history := service.NewWeighingService().CorrectConfirmed(record, 12, 4, "scale reconciliation", 9)
	if history.OriginalGross != 10 || history.OriginalTare != 3 {
		t.Fatalf("history lost confirmed values: %#v", history)
	}
	if record.Status != "corrected" || record.NetWeight != 8 {
		t.Fatalf("correction was not applied: %#v", record)
	}
}
