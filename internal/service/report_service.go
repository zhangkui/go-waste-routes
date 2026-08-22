package service

import (
	"sort"
	"time"
)

type DashboardStats struct {
	TodayTasks       int                `json:"today_tasks"`
	AbnormalCount    int                `json:"abnormal_count"`
	UnpaidInvoices   int                `json:"unpaid_invoices"`
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
	snapshot := CloneDashboardStats(base)
	snapshot.TodayTasks += extra.TodayTasks
	snapshot.AbnormalCount += extra.AbnormalCount
	snapshot.UnpaidInvoices += extra.UnpaidInvoices
	if snapshot.CollectionSeries == nil {
		snapshot.CollectionSeries = map[string]float64{}
	}
	for key, value := range extra.CollectionSeries {
		snapshot.CollectionSeries[key] += value
	}
	return snapshot
}
