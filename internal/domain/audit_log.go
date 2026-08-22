package domain

import "time"

type AuditLog struct {
	ID           int64     `json:"id"`
	UserID       *int64    `json:"user_id,omitempty"`
	Username     string    `json:"username"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	BeforeValue  string    `json:"before_value"`
	AfterValue   string    `json:"after_value"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	RequestID    string    `json:"request_id"`
	CreatedAt    time.Time `json:"created_at"`
}
