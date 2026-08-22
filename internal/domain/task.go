package domain

import "time"

type Task struct {
	Base
	TaskNumber  string     `json:"task_number"`
	PlanID      int64      `json:"plan_id"`
	RouteID     *int64     `json:"route_id,omitempty"`
	DriverID    *int64     `json:"driver_id,omitempty"`
	PlanDate    time.Time  `json:"plan_date"`
	Status      string     `json:"status"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	MileageKm   float64    `json:"mileage_km"`
	FuelLiters  float64    `json:"fuel_liters"`
	Stops       []TaskStop `json:"stops,omitempty"`
}

const (
	TaskPending   = "pending"
	TaskClaimed   = "claimed"
	TaskRunning   = "running"
	TaskCompleted = "completed"
	TaskSkipped   = "skipped"
	TaskAbnormal  = "abnormal"
)
