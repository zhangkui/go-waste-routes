package domain

type WeighingHistory struct {
	Base
	WeighingRecordID int64   `json:"weighing_record_id"`
	OriginalGross    float64 `json:"original_gross"`
	OriginalTare     float64 `json:"original_tare"`
	CorrectedGross   float64 `json:"corrected_gross"`
	CorrectedTare    float64 `json:"corrected_tare"`
	Reason           string  `json:"reason"`
	OperatorID       int64   `json:"operator_id"`
}
