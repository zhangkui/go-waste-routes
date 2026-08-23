package service

import (
	"fmt"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type TaskService struct{}

func NewTaskService() *TaskService { return &TaskService{} }

// BuildStopsFromRoute prepares task stops for execution.
func (s *TaskService) BuildStopsFromRoute(stops []domain.RouteStop) []domain.TaskStop {
	ordered := append([]domain.RouteStop(nil), stops...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	taskStops := make([]domain.TaskStop, 0, len(ordered))
	for _, stop := range ordered {
		taskStops = append(taskStops, domain.TaskStop{
			CustomerID: stop.CustomerID,
			Sequence:   stop.Sequence,
			Status:     domain.TaskPending,
		})
	}
	return taskStops
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
