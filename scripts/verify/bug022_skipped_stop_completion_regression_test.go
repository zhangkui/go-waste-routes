package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
	handlerpkg "go-waste-routes/internal/transport/http/handler"
)

func TestBug022_SkippedStopCannotHideCompletionError(t *testing.T) {
	t.Run("workflow rejects a pending stop", func(t *testing.T) {
		engine := service.NewWorkflowEngine()
		if err := engine.ValidateTaskStopsForCompletion([]domain.TaskStop{{Status: "pending"}}); err == nil {
			t.Fatal("pending stop must block task completion")
		}
	})

	t.Run("service propagates stop validation error", func(t *testing.T) {
		task := domain.Task{Status: domain.TaskRunning}
		err := service.NewTaskService().CompleteWithStops(&task, []domain.TaskStop{{Status: "pending"}}, task.PlanDate, 12, 3)
		if err == nil { t.Fatal("completion must return the outstanding-stop error") }
		if task.Status != domain.TaskRunning { t.Fatalf("task changed to %s after failed completion", task.Status) }
	})

	t.Run("http completion keeps task running when a skipped stop is not closed", func(t *testing.T) {
		tasks := service.NewResourceService(memory.New[domain.Task]())
		stops := service.NewResourceService(memory.New[domain.TaskStop]())
		created, err := tasks.Create(context.Background(), domain.Task{Status: domain.TaskRunning})
		if err != nil { t.Fatalf("create task: %v", err) }
		if _, err := stops.Create(context.Background(), domain.TaskStop{TaskID: created.ID, Status: "skipped"}); err != nil { t.Fatalf("create stop: %v", err) }
		body, _ := json.Marshal(map[string]float64{"mileage_km": 12, "fuel_liters": 3})
		req := httptest.NewRequest(http.MethodPut, "/tasks/1/complete", bytes.NewReader(body))
		rctx := chi.NewRouteContext(); rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		recorder := httptest.NewRecorder()
		handler := &handlerpkg.OpsHandler{Tasks: tasks, TaskStops: stops, TaskActions: service.NewTaskService()}
		handler.CompleteTask(recorder, req)
		if recorder.Code != http.StatusBadRequest { t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String()) }
		stored, err := tasks.Get(context.Background(), created.ID)
		if err != nil { t.Fatalf("reload task: %v", err) }
		if stored.Status != domain.TaskRunning { t.Fatalf("task was completed despite invalid skipped stop: %s", stored.Status) }
	})
}
