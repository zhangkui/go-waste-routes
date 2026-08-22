package service

import (
	"fmt"
	"time"

	"go-waste-routes/internal/domain"
)

type WeighingService struct{}

func NewWeighingService() *WeighingService { return &WeighingService{} }

func (s *WeighingService) CalculateNetWeight(gross, tare float64) float64 {
	return gross - tare
}

func (s *WeighingService) DetectAbnormality(record domain.WeighingRecord, history []domain.WeighingRecord, capacity float64) *domain.WeighingAbnormality {
	if record.NetWeight < 0 {
		return s.buildAbnormal(record.ID, "negative_net_weight", "鍑€閲嶅皬浜庨浂")
	}
	if capacity > 0 && (record.GrossWeight > capacity || record.TareWeight > capacity) {
		return s.buildAbnormal(record.ID, "exceeded_scale_capacity", "瓒呰繃鍦扮閲忕▼")
	}
	if record.Manual {
		return s.buildAbnormal(record.ID, "manual_entry", "鎵嬪伐褰曞叆")
	}
	if len(history) == 0 {
		return nil
	}
	avg := s.averageNetWeight(history)
	if avg > 0 && abs(record.NetWeight-avg)/avg > 0.5 {
		return s.buildAbnormal(record.ID, "deviation_over_50_percent", "鍋忕鍘嗗彶鍧囧€艰繃澶?)
	}
	return nil
}

func (s *WeighingService) Confirm(record *domain.WeighingRecord, confirmer int64, confirmedAt time.Time) {
	record.Status = "confirmed"
	record.ConfirmedBy = &confirmer
	record.ConfirmedAt = &confirmedAt
}

func (s *WeighingService) BuildCorrection(record domain.WeighingRecord, reason string, operatorID int64) domain.WeighingHistory {
	return domain.WeighingHistory{
		WeighingRecordID: record.ID,
		OriginalGross:    record.GrossWeight,
		OriginalTare:     record.TareWeight,
		CorrectedGross:   record.GrossWeight,
		CorrectedTare:    record.TareWeight,
		Reason:           reason,
		OperatorID:       operatorID,
	}
}

func (s *WeighingService) averageNetWeight(history []domain.WeighingRecord) float64 {
	total := 0.0
	for _, item := range history {
		total += item.NetWeight
	}
	return total / float64(len(history))
}

func (s *WeighingService) buildAbnormal(recordID int64, code, description string) *domain.WeighingAbnormality {
	return &domain.WeighingAbnormality{WeighingRecordID: recordID, Type: code, Description: description, Status: "pending"}
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func (s *WeighingService) Summary(record domain.WeighingRecord) string {
	return fmt.Sprintf("gross=%.2f tare=%.2f net=%.2f", record.GrossWeight, record.TareWeight, record.NetWeight)
}


