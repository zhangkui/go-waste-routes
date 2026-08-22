package domain

import (
	"strings"
	"time"
)

type TaskStop struct {
	Base
	TaskID           int64      `json:"task_id"`
	CustomerID       int64      `json:"customer_id"`
	Sequence         int        `json:"sequence"`
	Status           string     `json:"status"`
	ArrivedAt        *time.Time `json:"arrived_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Longitude        float64    `json:"longitude"`
	Latitude         float64    `json:"latitude"`
	ActualWeightTons float64    `json:"actual_weight_tons"`
	SkipReason       string     `json:"skip_reason"`
	Notes            string     `json:"notes"`
	Photos           string     `json:"photos"`
}

const (
	SkipReasonContainerNotFull = "container_not_full"
	SkipReasonChannelBlocked   = "channel_blocked"
	SkipReasonCustomerRequest  = "customer_request"
	SkipReasonEquipmentFailure = "equipment_failure"
	SkipReasonRoadClosed       = "road_closed"
	SkipReasonOther            = "other"
)

var AllowedSkipReasons = []string{
	SkipReasonContainerNotFull,
	SkipReasonChannelBlocked,
	SkipReasonCustomerRequest,
	SkipReasonEquipmentFailure,
	SkipReasonRoadClosed,
	SkipReasonOther,
}

// IsAllowedSkipReason reports whether the given value matches one of the
// predefined skip reasons. Comparison is case-insensitive and trims
// surrounding whitespace, mirroring NormalizeStatus semantics used by the
// rest of the service layer.
func IsAllowedSkipReason(reason string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(reason))
	if trimmed == "" {
		return false
	}
	for _, allowed := range AllowedSkipReasons {
		if allowed == trimmed {
			return true
		}
	}
	return false
}
