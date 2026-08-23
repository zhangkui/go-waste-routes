package verify

import (
	"context"
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug044_CancelledArrivalHasNoSideEffect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	stop := &domain.TaskStop{Status: "pending"}
	err := service.NewTaskService().ArriveWithContext(ctx, stop, 121.5, 31.2, time.Now())
	if err == nil {
		t.Fatal("cancelled arrival must return context error")
	}
	if stop.Status != "pending" || stop.ArrivedAt != nil {
		t.Fatalf("cancelled arrival mutated task stop: %#v", stop)
	}
}
