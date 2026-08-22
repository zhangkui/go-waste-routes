package domain

type Route struct {
	Base
	Name                string      `json:"name"`
	PlanID              int64       `json:"plan_id"`
	VehicleID           int64       `json:"vehicle_id"`
	DriverID            int64       `json:"driver_id"`
	Status              string      `json:"status"`
	EstimatedWeightTons float64     `json:"estimated_weight_tons"`
	Stops               []RouteStop `json:"stops,omitempty"`
}

const (
	RoutePending   = "pending"
	RouteRunning   = "running"
	RouteCompleted = "completed"
	RouteCancelled = "cancelled"
	RouteAbnormal  = "abnormal"
)
