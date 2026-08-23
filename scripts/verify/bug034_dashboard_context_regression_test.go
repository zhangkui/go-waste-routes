package verify

import (
	"context"
	"errors"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
	"testing"
	"time"
)

func TestBug034_DashboardStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	newCustomer := service.NewResourceService(memory.New[domain.Customer]())
	newTask := service.NewResourceService(memory.New[domain.Task]())
	newRoute := service.NewResourceService(memory.New[domain.Route]())
	newWeight := service.NewResourceService(memory.New[domain.WeighingRecord]())
	newAbnormal := service.NewResourceService(memory.New[domain.WeighingAbnormality]())
	newInvoice := service.NewResourceService(memory.New[domain.Invoice]())
	_, err := service.NewDashboardService(newCustomer, newTask, newRoute, newWeight, newAbnormal, newInvoice).Snapshot(ctx, time.Now())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled dashboard error = %v, want context canceled", err)
	}
}
