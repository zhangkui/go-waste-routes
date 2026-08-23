package service

import (
	"fmt"
	"time"

	"go-waste-routes/internal/domain"
)

type TaskService struct{}

func NewTaskService() *TaskService { return &TaskService{} }

func (s *TaskService) SummarizeBatch(tasks []domain.Task, acquire func(domain.Task) func()) int {
	completed := 0
	for _, task := range tasks {
		completed += s.summarizeOne(task, acquire)
	}
	return completed
}

// summarizeOne processes a single task within its own resource scope. The
// release callback is deferred to the per-task boundary so temporary audit
// and computation contexts are freed as soon as the item finishes instead of
// lingering until the whole batch returns. A panic in one item is recovered
// so that only the failing task is skipped; the remaining tasks still run.
func (s *TaskService) summarizeOne(task domain.Task, acquire func(domain.Task) func()) (count int) {
	defer func() {
		_ = recover()
	}()
	release := acquire(task)
	defer release()
	if NormalizeStatus(task.Status) == domain.TaskCompleted {
		return 1
	}
	return 0
}

func (s *TaskService) StatusFlow() map[string][]string {
	return map[string][]string{
		domain.TaskPending:   {domain.TaskClaimed, domain.TaskAbnormal},
		domain.TaskClaimed:   {domain.TaskRunning, domain.TaskSkipped},
		domain.TaskRunning:   {domain.TaskCompleted, domain.TaskSkipped, domain.TaskAbnormal},
		domain.TaskSkipped:   {domain.TaskAbnormal},
		domain.TaskCompleted: {},
		domain.TaskAbnormal:  {},
	}
}

func (s *TaskService) Claim(task *domain.Task, claimedAt time.Time) error {
	if err := ValidateStatusTransition(task.Status, domain.TaskClaimed, s.StatusFlow()); err != nil {
		return err
	}
	task.Status = domain.TaskClaimed
	task.ClaimedAt = &claimedAt
	return nil
}

func (s *TaskService) Arrive(stop *domain.TaskStop, longitude, latitude float64, arrivedAt time.Time) {
	stop.Status = "arrived"
	stop.Longitude = longitude
	stop.Latitude = latitude
	stop.ArrivedAt = &arrivedAt
}

func (s *TaskService) Skip(stop *domain.TaskStop, reason string, note string) error {
	if reason == "" {
		return fmt.Errorf("skip reason required")
	}
	stop.Status = "skipped"
	stop.SkipReason = reason
	stop.Notes = note
	return nil
}

func (s *TaskService) Complete(task *domain.Task, completedAt time.Time, mileage, fuel float64) error {
	if err := ValidateStatusTransition(task.Status, domain.TaskCompleted, s.StatusFlow()); err != nil {
		return err
	}
	task.Status = domain.TaskCompleted
	task.CompletedAt = &completedAt
	task.MileageKm = mileage
	task.FuelLiters = fuel
	return nil
}

func (s *TaskService) BuildTaskNumber(planDate time.Time, sequence int) string {
	return fmt.Sprintf("TK-%s-%04d", planDate.Format("20060102"), sequence)
}
