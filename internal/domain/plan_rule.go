package domain

type PlanRule struct {
	Base
	PlanID    int64  `json:"plan_id"`
	RuleType  string `json:"rule_type"`
	RuleValue string `json:"rule_value"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Enabled   bool   `json:"enabled"`
}
