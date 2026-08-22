package domain

import "time"

type Driver struct {
	Base
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	LicenseNumber string    `json:"license_number"`
	LicenseClass  string    `json:"license_class"`
	HireDate      time.Time `json:"hire_date"`
	Status        string    `json:"status"`
}
