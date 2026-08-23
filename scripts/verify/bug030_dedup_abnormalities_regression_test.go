package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug030_DedupKeepsDistinctAbnormalities(t *testing.T) {
	engine := service.NewWeighingEngine()
	record := domain.WeighingRecord{Base: domain.Base{ID: 30}, Manual: true, GrossWeight: 12, TareWeight: 2, NetWeight: 10}
	bridge := domain.Weighbridge{MaxCapacityTons: 5, Status: "normal"}
	issues := engine.DetectAbnormalities(record, nil, bridge)
	if len(issues) < 2 {
		t.Fatalf("distinct abnormalities were collapsed: %#v", issues)
	}
}
