package domain

import "time"

type HolidayConfig struct {
	Base
	HolidayDate time.Time `json:"holiday_date"`
	Name        string    `json:"name"`
	IsWorkday   bool      `json:"is_workday"`
	Enabled     bool      `json:"enabled"`
}
