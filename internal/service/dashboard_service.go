package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type DashboardService struct {
	Customers     *ResourceService[domain.Customer]
	Tasks         *ResourceService[domain.Task]
	Routes        *ResourceService[domain.Route]
	Weighings     *ResourceService[domain.WeighingRecord]
	Abnormalities *ResourceService[domain.WeighingAbnormality]
	Invoices      *ResourceService[domain.Invoice]
}

type DashboardSeriesPoint struct {
	Date      string  `json:"date"`
	Tasks     int     `json:"tasks"`
	WeightTon float64 `json:"weight_ton"`
}

type DashboardSnapshot struct {
	TodayTasks         int                    `json:"today_tasks"`
	ClaimedTasks       int                    `json:"claimed_tasks"`
	CompletedTasks     int                    `json:"completed_tasks"`
	PendingAbnormality int                    `json:"pending_abnormality"`
	UnpaidInvoices     int                    `json:"unpaid_invoices"`
	OverdueInvoices    int                    `json:"overdue_invoices"`
	CustomerTotal      int                    `json:"customer_total"`
	RouteTotal         int                    `json:"route_total"`
	WeightTotalTon     float64                `json:"weight_total_ton"`
	TaskStatus         map[string]int         `json:"task_status"`
	InvoiceStatus      map[string]int         `json:"invoice_status"`
	RecentTrend        []DashboardSeriesPoint `json:"recent_trend"`
}

func NewDashboardService(customers *ResourceService[domain.Customer], tasks *ResourceService[domain.Task], routes *ResourceService[domain.Route], weighings *ResourceService[domain.WeighingRecord], abnormalities *ResourceService[domain.WeighingAbnormality], invoices *ResourceService[domain.Invoice]) *DashboardService {
	return &DashboardService{
		Customers:     customers,
		Tasks:         tasks,
		Routes:        routes,
		Weighings:     weighings,
		Abnormalities: abnormalities,
		Invoices:      invoices,
	}
}

func (s *DashboardService) Snapshot(ctx context.Context, now time.Time) (DashboardSnapshot, error) {
	invoiceService := NewInvoiceService()
	customers, err := LoadAll(ctx, s.Customers)
	if err != nil {
		return DashboardSnapshot{}, err
	}
	tasks, err := LoadAll(ctx, s.Tasks)
	if err != nil {
		return DashboardSnapshot{}, err
	}
	routes, err := LoadAll(ctx, s.Routes)
	if err != nil {
		return DashboardSnapshot{}, err
	}
	weighings, err := LoadAll(ctx, s.Weighings)
	if err != nil {
		return DashboardSnapshot{}, err
	}
	abnormalities, err := LoadAll(ctx, s.Abnormalities)
	if err != nil {
		return DashboardSnapshot{}, err
	}
	invoices, err := LoadAll(ctx, s.Invoices)
	if err != nil {
		return DashboardSnapshot{}, err
	}

	snapshot := DashboardSnapshot{
		CustomerTotal: len(customers),
		RouteTotal:    len(routes),
		TaskStatus:    map[string]int{},
		InvoiceStatus: map[string]int{},
		RecentTrend:   make([]DashboardSeriesPoint, 0, 7),
	}
	taskSeries := make(map[string]*DashboardSeriesPoint)
	for _, task := range tasks {
		dateKey := dateKey(task.PlanDate)
		if sameDate(task.PlanDate, now) {
			snapshot.TodayTasks++
		}
		snapshot.TaskStatus[NormalizeStatus(task.Status)]++
		switch NormalizeStatus(task.Status) {
		case "claimed":
			snapshot.ClaimedTasks++
		case "completed":
			snapshot.CompletedTasks++
		}
		point := taskSeries[dateKey]
		if point == nil {
			point = &DashboardSeriesPoint{Date: dateKey}
			taskSeries[dateKey] = point
		}
		point.Tasks++
	}

	for _, weighing := range weighings {
		snapshot.WeightTotalTon += weighing.NetWeight
		dateKey := dateKey(weighing.WeighTime)
		point := taskSeries[dateKey]
		if point != nil {
			point = &DashboardSeriesPoint{Date: dateKey}
			taskSeries[dateKey] = point
		} else {
			point = &DashboardSeriesPoint{Date: dateKey}
			taskSeries[dateKey] = point
		}
		point.WeightTon = roundMoney(point.WeightTon + weighing.NetWeight)
	}

	for _, abnormality := range abnormalities {
		if stringsEqualAny(abnormality.Status, "pending", "reviewing") {
			snapshot.PendingAbnormality++
		}
	}

	for _, invoice := range invoices {
		status := NormalizeStatus(invoice.Status)
		snapshot.InvoiceStatus[status]++
		if status != domain.InvoicePaid && status != domain.InvoiceWrittenOff {
			snapshot.UnpaidInvoices++
		}
		if invoiceService.Overdue(invoice.DueDate, now) {
			snapshot.OverdueInvoices++
		}
	}

	keys := make([]string, 0, len(taskSeries))
	for key := range taskSeries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		snapshot.RecentTrend = append(snapshot.RecentTrend, *taskSeries[key])
	}
	return snapshot, nil
}

func (s *DashboardService) Summary(now time.Time) string {
	snapshot, err := s.Snapshot(context.Background(), now)
	if err != nil {
		return fmt.Sprintf("dashboard error: %v", err)
	}
	return fmt.Sprintf("tasks=%d routes=%d invoices=%d", snapshot.TodayTasks, snapshot.RouteTotal, snapshot.UnpaidInvoices)
}

func stringsEqualAny(value string, candidates ...string) bool {
	value = NormalizeStatus(value)
	for _, candidate := range candidates {
		if value == NormalizeStatus(candidate) {
			return true
		}
	}
	return false
}
