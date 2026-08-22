package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug007_ResolveMarksCleanWeighingAsConfirmed(t *testing.T) {
	engine := service.NewWeighingEngine()
	record := domain.WeighingRecord{Base: domain.Base{ID: 77}, Status: "draft"}

	got := engine.Resolve(&record, nil)
	if got != "confirmed" {
		t.Fatalf("resolve result = %s, want confirmed", got)
	}
	if record.Status != "confirmed" {
		t.Fatalf("record status = %s, want confirmed", record.Status)
	}
}
