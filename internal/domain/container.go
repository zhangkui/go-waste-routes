package domain

type Container struct {
	Base
	CustomerID int64   `json:"customer_id"`
	Type       string  `json:"type"`
	Number     string  `json:"number"`
	CapacityL  int     `json:"capacity_l"`
	Quantity   int     `json:"quantity"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Location   string  `json:"location"`
	Status     string  `json:"status"`
}
