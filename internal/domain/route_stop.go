package domain

type RouteStop struct {
	Base
	RouteID             int64   `json:"route_id"`
	CustomerID          int64   `json:"customer_id"`
	Sequence            int     `json:"sequence"`
	EstimatedArrival    string  `json:"estimated_arrival"`
	StayMinutes         int     `json:"stay_minutes"`
	EstimatedWeightTons float64 `json:"estimated_weight_tons"`
}
