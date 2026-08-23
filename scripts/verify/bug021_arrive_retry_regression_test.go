package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
	handlerpkg "go-waste-routes/internal/transport/http/handler"
)

func TestBug021_ArriveRetryPreservesFirstArrival(t *testing.T) {
	t.Run("workflow treats an arrived stop as an idempotent replay", func(t *testing.T) {
		engine := service.NewWorkflowEngine()
		changed, err := engine.PlanTaskStopArrival(&domain.TaskStop{Status: "arrived"})
		if err != nil {
			t.Fatalf("plan arrived retry: %v", err)
		}
		if changed {
			t.Fatal("arrived retry must not request another state mutation")
		}
	})

	t.Run("concurrent retries preserve the first arrival", func(t *testing.T) {
		taskService := service.NewTaskService()
		first := time.Date(2026, 8, 23, 8, 0, 0, 0, time.UTC)
		stop := domain.TaskStop{Status: "pending"}
		changed, err := taskService.Arrive(&stop, 121.10, 31.20, first)
		if err != nil || !changed {
			t.Fatalf("first arrival failed: changed=%v err=%v", changed, err)
		}

		var wait sync.WaitGroup
		for attempt := 1; attempt <= 8; attempt++ {
			attempt := attempt
			wait.Add(1)
			go func() {
				defer wait.Done()
				_, _ = taskService.Arrive(&stop, 121.10+float64(attempt), 31.20+float64(attempt), first.Add(time.Duration(attempt)*time.Minute))
			}()
		}
		wait.Wait()

		if stop.ArrivedAt == nil || !stop.ArrivedAt.Equal(first) {
			t.Fatalf("arrival time changed after retries: %v", stop.ArrivedAt)
		}
		if stop.Longitude != 121.10 || stop.Latitude != 31.20 {
			t.Fatalf("arrival coordinates changed after retries: %.2f, %.2f", stop.Longitude, stop.Latitude)
		}
	})

	t.Run("http retry does not perform another persistence update", func(t *testing.T) {
		store := memory.New[domain.TaskStop]()
		resources := service.NewResourceService(store)
		first := time.Date(2026, 8, 23, 8, 0, 0, 0, time.UTC)
		created, err := resources.Create(context.Background(), domain.TaskStop{
			Status: "arrived", ArrivedAt: &first, Longitude: 121.10, Latitude: 31.20,
		})
		if err != nil {
			t.Fatalf("create task stop: %v", err)
		}
		originalUpdatedAt := created.UpdatedAt
		time.Sleep(2 * time.Millisecond)

		payload, _ := json.Marshal(map[string]any{
			"longitude":  122.20,
			"latitude":   32.30,
			"arrived_at": first.Add(5 * time.Minute),
		})
		req := httptest.NewRequest(http.MethodPut, "/task-stops/1/arrive", bytes.NewReader(payload))
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
		recorder := httptest.NewRecorder()
		handler := &handlerpkg.OpsHandler{TaskStops: resources, TaskActions: service.NewTaskService()}
		handler.ArriveTaskStop(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("retry status=%d body=%s", recorder.Code, recorder.Body.String())
		}

		stored, err := resources.Get(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("reload task stop: %v", err)
		}
		if !stored.UpdatedAt.Equal(originalUpdatedAt) {
			t.Fatalf("retry unexpectedly persisted task stop: before=%s after=%s", originalUpdatedAt, stored.UpdatedAt)
		}
		if stored.ArrivedAt == nil || !stored.ArrivedAt.Equal(first) || stored.Longitude != 121.10 || stored.Latitude != 31.20 {
			t.Fatalf("stored first arrival was overwritten: %+v", stored)
		}
	})
}
