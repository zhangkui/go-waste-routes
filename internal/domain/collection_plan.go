package domain

import "time"

type CollectionPlan struct {
	Base
	Name                string     `json:"name"`
	CustomerID          int64      `json:"customer_id"`
	Version             int        `json:"version"`
	Status              string     `json:"status"`
	HolidayStrategy     string     `json:"holiday_strategy"`
	EffectiveFrom       time.Time  `json:"effective_from"`
	EffectiveTo         *time.Time `json:"effective_to,omitempty"`
	EstimatedWeightTons float64    `json:"estimated_weight_tons"`
	Rules               []PlanRule `json:"rules,omitempty"`
}

const (
	PlanDraft    = "draft"
	PlanActive   = "active"
	PlanDisabled = "disabled"
	PlanArchived = "archived"
)
