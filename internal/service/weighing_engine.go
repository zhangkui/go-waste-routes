package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type WeightBaseline struct {
	SampleSize int
	GrossMean  float64
	TareMean   float64
	NetMean    float64
}

type WeightDeviation struct {
	Field      string
	Baseline   float64
	Actual     float64
	Percent    float64
	Exceeded   bool
}

type WeighingEngine struct{}

func NewWeighingEngine() *WeighingEngine { return &WeighingEngine{} }

func (e *WeighingEngine) NetWeight(grossWeight, tareWeight float64) float64 {
	return roundMoney(grossWeight - tareWeight)
}

func (e *WeighingEngine) BuildBaseline(records []domain.WeighingRecord) WeightBaseline {
	if len(records) == 0 {
		return WeightBaseline{}
	}
	var grossTotal float64
	var tareTotal float64
	var netTotal float64
	for _, record := range records {
		grossTotal += record.GrossWeight
		tareTotal += record.TareWeight
		netTotal += record.NetWeight
	}
	sampleSize := len(records)
	return WeightBaseline{
		SampleSize: sampleSize,
		GrossMean:  roundMoney(grossTotal / float64(sampleSize)),
		TareMean:   roundMoney(tareTotal / float64(sampleSize)),
		NetMean:    roundMoney(netTotal / float64(sampleSize)),
	}
}

func (e *WeighingEngine) BaselineForCustomer(records []domain.WeighingRecord, customerID int64) WeightBaseline {
	filtered := make([]domain.WeighingRecord, 0)
	for _, record := range records {
		if record.CustomerID == customerID {
			filtered = append(filtered, record)
		}
	}
	return e.BuildBaseline(filtered)
}

func (e *WeighingEngine) RecentAverage(records []domain.WeighingRecord, limit int) float64 {
	if len(records) == 0 || limit <= 0 {
		return 0
	}
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].WeighTime.After(records[j].WeighTime)
	})
	if limit > len(records) {
		limit = len(records)
	}
	var total float64
	for index := 0; index < limit; index++ {
		total += records[index].NetWeight
	}
	return roundMoney(total / float64(limit))
}

func (e *WeighingEngine) DeviationPercent(actual, baseline float64) float64 {
	if baseline == 0 {
		if actual == 0 {
			return 0
		}
		return 100
	}
	return roundMoney(math.Abs(actual-baseline) / math.Abs(baseline) * 100)
}

func (e *WeighingEngine) DetectAbnormalities(record domain.WeighingRecord, history []domain.WeighingRecord, bridge domain.Weighbridge) []domain.WeighingAbnormality {
	issues := make([]domain.WeighingAbnormality, 0)
	if record.Manual {
		issues = append(issues, e.abnormality(record.ID, "manual_entry", "manual weighing entry requires review", "pending"))
	}
	if record.NetWeight < 0 {
		issues = append(issues, e.abnormality(record.ID, "negative_net", "net weight is negative", "pending"))
	}
	if bridge.MaxCapacityTons > 0 && (record.GrossWeight > bridge.MaxCapacityTons || record.TareWeight > bridge.MaxCapacityTons) {
		issues = append(issues, e.abnormality(record.ID, "scale_range", "weight exceeds scale capacity", "pending"))
	}
	baseline := e.BaselineForCustomer(history, record.CustomerID)
	if baseline.SampleSize > 0 {
		deviation := e.DeviationPercent(record.NetWeight, baseline.NetMean)
		if deviation >= 50 {
			issues = append(issues, e.abnormality(record.ID, "history_deviation", fmt.Sprintf("net weight deviates %.2f%% from historical average", deviation), "pending"))
		}
	}
	if bridge.Status != "" && strings.EqualFold(bridge.Status, "fault") {
		issues = append(issues, e.abnormality(record.ID, "scale_fault", "scale is faulted", "pending"))
	}
	return dedupeAbnormalities(issues)
}

func (e *WeighingEngine) abnormality(recordID int64, kind, description, status string) domain.WeighingAbnormality {
	return domain.WeighingAbnormality{
		WeighingRecordID: recordID,
		Type:             kind,
		Description:      description,
		Status:           status,
		ReviewResult:     "pending",
	}
}

func (e *WeighingEngine) CreateCorrectionHistory(original, corrected domain.WeighingRecord, operatorID int64, reason string) domain.WeighingHistory {
	return domain.WeighingHistory{
		WeighingRecordID: original.ID,
		OriginalGross:    original.GrossWeight,
		OriginalTare:     original.TareWeight,
		CorrectedGross:   corrected.GrossWeight,
		CorrectedTare:    corrected.TareWeight,
		Reason:           strings.TrimSpace(reason),
		OperatorID:       operatorID,
	}
}

func (e *WeighingEngine) Confirm(record *domain.WeighingRecord, confirmer int64, confirmedAt time.Time) {
	if record == nil {
		return
	}
	record.Status = "confirmed"
	record.ConfirmedBy = &confirmer
	record.ConfirmedAt = &confirmedAt
}

func (e *WeighingEngine) Reject(record *domain.WeighingRecord) {
	if record == nil {
		return
	}
	record.Status = "rejected"
}

func (e *WeighingEngine) Review(abnormality *domain.WeighingAbnormality, reviewerID int64, result, rejectReason string, reviewedAt time.Time) {
	if abnormality == nil {
		return
	}
	abnormality.ReviewerID = &reviewerID
	abnormality.ReviewedAt = &reviewedAt
	abnormality.ReviewResult = NormalizeStatus(result)
	abnormality.Status = NormalizeStatus(result)
	if strings.EqualFold(abnormality.ReviewResult, "confirmed") {
		// Clearing the reject reason on confirmation prevents a stale
		// reason from a prior rejection lingering after the abnormality
		// is approved. Reviewer and review time remain intact.
		abnormality.RejectReason = ""
	} else {
		abnormality.RejectReason = strings.TrimSpace(rejectReason)
	}
}

func (e *WeighingEngine) Resolve(record *domain.WeighingRecord, abnormalities []domain.WeighingAbnormality) string {
	if len(abnormalities) == 0 {
		return "confirmed"
	}
	for _, abnormality := range abnormalities {
		if strings.EqualFold(abnormality.Type, "negative_net") {
			record.Status = "abnormal"
			return "abnormal"
		}
	}
	record.Status = "reviewing"
	return "reviewing"
}

func (e *WeighingEngine) IsManual(record domain.WeighingRecord) bool {
	return record.Manual || strings.TrimSpace(record.ManualReason) != ""
}

func (e *WeighingEngine) IsOutOfRange(record domain.WeighingRecord, bridge domain.Weighbridge) bool {
	if bridge.MaxCapacityTons <= 0 {
		return false
	}
	return record.GrossWeight > bridge.MaxCapacityTons || record.TareWeight > bridge.MaxCapacityTons
}

func (e *WeighingEngine) Snapshot(record domain.WeighingRecord, baseline WeightBaseline) map[string]any {
	return map[string]any{
		"record_id":      record.ID,
		"customer_id":    record.CustomerID,
		"gross_weight":   record.GrossWeight,
		"tare_weight":    record.TareWeight,
		"net_weight":     record.NetWeight,
		"baseline_net":   baseline.NetMean,
		"deviation_pct":  e.DeviationPercent(record.NetWeight, baseline.NetMean),
		"is_manual":      record.Manual,
		"manual_reason":  strings.TrimSpace(record.ManualReason),
		"status":         record.Status,
	}
}

func (e *WeighingEngine) BucketByCustomer(records []domain.WeighingRecord) map[int64][]domain.WeighingRecord {
	result := make(map[int64][]domain.WeighingRecord)
	for _, record := range records {
		result[record.CustomerID] = append(result[record.CustomerID], record)
	}
	return result
}

func (e *WeighingEngine) CustomerSummary(records []domain.WeighingRecord, customerID int64) WeightBaseline {
	filtered := make([]domain.WeighingRecord, 0)
	for _, record := range records {
		if record.CustomerID == customerID {
			filtered = append(filtered, record)
		}
	}
	return e.BuildBaseline(filtered)
}

func (e *WeighingEngine) Deviations(records []domain.WeighingRecord, baseline WeightBaseline) []WeightDeviation {
	deviations := make([]WeightDeviation, 0, len(records))
	for _, record := range records {
		deviations = append(deviations, WeightDeviation{
			Field:    "net_weight",
			Baseline: baseline.NetMean,
			Actual:   record.NetWeight,
			Percent:  e.DeviationPercent(record.NetWeight, baseline.NetMean),
			Exceeded: e.DeviationPercent(record.NetWeight, baseline.NetMean) >= 50,
		})
	}
	return deviations
}

func (e *WeighingEngine) SortByNewest(records []domain.WeighingRecord) {
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].WeighTime.Equal(records[j].WeighTime) {
			return records[i].ID > records[j].ID
		}
		return records[i].WeighTime.After(records[j].WeighTime)
	})
}

func dedupeAbnormalities(items []domain.WeighingAbnormality) []domain.WeighingAbnormality {
	seen := make(map[string]struct{}, len(items))
	result := make([]domain.WeighingAbnormality, 0, len(items))
	for _, item := range items {
		key := fmt.Sprintf("%d:%s", item.WeighingRecordID, item.Type)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].WeighingRecordID == result[j].WeighingRecordID {
			return result[i].Type < result[j].Type
		}
		return result[i].WeighingRecordID < result[j].WeighingRecordID
	})
	return result
}

