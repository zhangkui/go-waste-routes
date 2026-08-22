package domain

type Customer struct {
	Base
	Name             string                  `json:"name"`
	ContactName      string                  `json:"contact_name"`
	ContactPhone     string                  `json:"contact_phone"`
	Address          string                  `json:"address"`
	Region           string                  `json:"region"`
	Level            string                  `json:"level"`
	ContractNumber   string                  `json:"contract_number"`
	ServiceWeekdays  string                  `json:"service_weekdays"`
	ServiceMonthDays string                  `json:"service_month_days"`
	TemporaryAllowed bool                    `json:"temporary_allowed"`
	ServiceWindows   string                  `json:"service_windows"`
	Status           string                  `json:"status"`
	WasteCategories  []CustomerWasteCategory `json:"waste_categories,omitempty"`
}

type CustomerWasteCategory struct {
	CustomerID      int64   `json:"customer_id"`
	WasteCategoryID int64   `json:"waste_category_id"`
	Weight          float64 `json:"weight"`
}
