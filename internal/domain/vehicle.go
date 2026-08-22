package domain

import "time"

type Vehicle struct {
	Base
	PlateNumber       string    `json:"plate_number"`
	Model             string    `json:"model"`
	RatedLoadTons     float64   `json:"rated_load_tons"`
	CurrentMileageKm  float64   `json:"current_mileage_km"`
	Status            string    `json:"status"`
	InspectionDueDate time.Time `json:"inspection_due_date"`
}
