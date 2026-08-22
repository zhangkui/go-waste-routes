package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug018_CorrectionHistoryTrimsReason(t *testing.T) {
	weightService := service.NewWeighingService()
	history := weightService.BuildCorrection(domain.WeighingRecord{Base: domain.Base{ID: 7}, GrossWeight: 10, TareWeight: 4}, "  recheck  ", 99)
	if history.Reason != "recheck" {
		t.Fatalf("reason = %q, want recheck", history.Reason)
	}
}
