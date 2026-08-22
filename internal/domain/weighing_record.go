package domain

import "time"

type WeighingRecord struct {
	Base
	TaskID       int64      `json:"task_id"`
	CustomerID   int64      `json:"customer_id"`
	VehicleID    int64      `json:"vehicle_id"`
	ScaleID      int64      `json:"scale_id"`
	GrossWeight  float64    `json:"gross_weight"`
	TareWeight   float64    `json:"tare_weight"`
	NetWeight    float64    `json:"net_weight"`
	WeighTime    time.Time  `json:"weigh_time"`
	Manual       bool       `json:"manual"`
	ManualReason string     `json:"manual_reason"`
	Status       string     `json:"status"`
	ConfirmedBy  *int64     `json:"confirmed_by,omitempty"`
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty"`
}

type Weighbridge struct {
	Base
	Number          string    `json:"number"`
	Name            string    `json:"name"`
	Location        string    `json:"location"`
	MaxCapacityTons float64   `json:"max_capacity_tons"`
	CalibrationDate time.Time `json:"calibration_date"`
	Status          string    `json:"status"`
}
