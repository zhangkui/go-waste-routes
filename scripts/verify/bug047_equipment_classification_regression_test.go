package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug047_EquipmentAbnormalityRequiresAbnormalStatus(t *testing.T) {
	record := domain.WeighingRecord{Base: domain.Base{ID: 47}, Status: "reviewing"}
	result := service.NewWeighingEngine().Resolve(&record, []domain.WeighingAbnormality{{WeighingRecordID: 47, Type: "scale_range"}})
	if result != "abnormal" || record.Status != "abnormal" {
		t.Fatalf("equipment anomaly was not classified as abnormal: result=%s status=%s", result, record.Status)
	}
}
