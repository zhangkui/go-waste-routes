package domain

import "time"

type StatusHistory struct {
	ID         int64     `json:"id"`
	EntityID   int64     `json:"entity_id"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	OperatorID int64     `json:"operator_id"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
}
