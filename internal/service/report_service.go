package service

import (
	"fmt"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

// ValidateTaskSummaryStops ensures a completion summary accounts for every stop.
func ValidateTaskSummaryStops(stops []domain.TaskStop) error {
	if len(stops) == 0 {
		return fmt.Errorf("task summary requires stops")
	}
	for _, stop := range stops {
		if NormalizeStatus(stop.Status) != "completed" && NormalizeStatus(stop.Status) != "skipped" {
			return fmt.Errorf("task summary contains unfinished stop")
		}
	}
	return nil
}

type DashboardStats struct {
	TodayTasks       int            `json:"today_tasks"`
	AbnormalCount    int            `json:"abnormal_count"`
	UnpaidInvoices   int            `json:"unpaid_invoices"`
	CollectionSeries map[string]float64 `json:"collection_series"`
}

type ReportService struct{}

func NewReportService() *ReportService { return &ReportService{} }

func (s *ReportService) BuildCollectionSeries(tasks []time.Time, weights []float64) map[string]float64 {
	series := make(map[string]float64)
	for index, date := range tasks {
		key := date.Format("2006-01-02")
		if index < len(weights) {
			series[key] += weights[index]
		}
	}
	return series
}

func (s *ReportService) SortKeys(values map[string]float64) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *ReportService) MergeStats(base DashboardStats, extra DashboardStats) DashboardStats {
	base.TodayTasks += extra.TodayTasks
	base.AbnormalCount += extra.AbnormalCount
	base.UnpaidInvoices += extra.UnpaidInvoices
	if base.CollectionSeries == nil {
		base.CollectionSeries = map[string]float64{}
	}
	for key, value := range extra.CollectionSeries {
		base.CollectionSeries[key] += value
	}
	return base
}

