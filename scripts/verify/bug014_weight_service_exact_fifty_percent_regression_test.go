package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug014_WeightServiceFlagsExactFiftyPercentDeviation(t *testing.T) {
	weightService := service.NewWeighingService()
	history := []domain.WeighingRecord{{CustomerID: 3, NetWeight: 100, WeighTime: time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)}}
	record := domain.WeighingRecord{Base: domain.Base{ID: 3}, CustomerID: 3, NetWeight: 150}

	abnormality := weightService.DetectAbnormality(record, history, 0)
	if abnormality == nil || abnormality.Type != "deviation_over_50_percent" {
		t.Fatalf("abnormality = %#v, want deviation_over_50_percent", abnormality)
	}
}
