package domain

import "time"

type TaskStop struct {
	Base
	TaskID           int64      `json:"task_id"`
	CustomerID       int64      `json:"customer_id"`
	Sequence         int        `json:"sequence"`
	Status           string     `json:"status"`
	ArrivedAt        *time.Time `json:"arrived_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Longitude        float64    `json:"longitude"`
	Latitude         float64    `json:"latitude"`
	ActualWeightTons float64    `json:"actual_weight_tons"`
	SkipReason       string     `json:"skip_reason"`
	Notes            string     `json:"notes"`
	Photos           string     `json:"photos"`
}
