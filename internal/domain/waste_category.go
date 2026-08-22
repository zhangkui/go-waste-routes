package domain

type WasteCategory struct {
	Base
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	DefaultRate float64 `json:"default_rate"`
	Status      string  `json:"status"`
}
