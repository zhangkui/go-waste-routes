package verify

import (
	"context"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
	"testing"
	"time"
)

func TestBug036_DashboardTrendKeepsSameDayCategories(t *testing.T) {
	day := time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC)
	customers := service.NewResourceService(memory.New[domain.Customer]())
	tasks := service.NewResourceService(memory.New[domain.Task]())
	routes := service.NewResourceService(memory.New[domain.Route]())
	weights := service.NewResourceService(memory.New[domain.WeighingRecord]())
	abnormal := service.NewResourceService(memory.New[domain.WeighingAbnormality]())
	invoices := service.NewResourceService(memory.New[domain.Invoice]())
	tasks.Create(context.Background(), domain.Task{PlanDate: day, Status: domain.TaskPending})
	weights.Create(context.Background(), domain.WeighingRecord{WeighTime: day, NetWeight: 3})
	snapshot, err := service.NewDashboardService(customers, tasks, routes, weights, abnormal, invoices).Snapshot(context.Background(), day)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.RecentTrend) != 1 || snapshot.RecentTrend[0].Tasks != 1 || snapshot.RecentTrend[0].WeightTon != 3 {
		t.Fatalf("same-day trend lost a category: %#v", snapshot.RecentTrend)
	}
}
