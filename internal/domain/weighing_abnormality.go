package domain

import "time"

type WeighingAbnormality struct {
	Base
	WeighingRecordID int64      `json:"weighing_record_id"`
	Type             string     `json:"type"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	ReviewerID       *int64     `json:"reviewer_id,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	ReviewResult     string     `json:"review_result"`
	RejectReason     string     `json:"reject_reason"`
}
