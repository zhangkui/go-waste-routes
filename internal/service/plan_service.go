package service

import (
	"fmt"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type PlanService struct {
	holidayChecker func(time.Time) bool
}

func NewPlanService(holidayChecker func(time.Time) bool) *PlanService {
	return &PlanService{holidayChecker: holidayChecker}
}

func (s *PlanService) Preview(plan domain.CollectionPlan, days int) []time.Time {
	if days <= 0 {
		days = 7
	}
	start := plan.EffectiveFrom
	if start.IsZero() {
		start = time.Now()
	}
	results := make([]time.Time, 0, days)
	current := start
	for len(results) < days {
		if s.shouldRunOn(plan, current) {
			results = append(results, current)
		}
		current = current.AddDate(0, 0, 1)
	}
	return results
}

func (s *PlanService) shouldRunOn(plan domain.CollectionPlan, day time.Time) bool {
	if s.holidayChecker != nil && s.holidayChecker(day) {
		switch NormalizeStatus(plan.HolidayStrategy) {
		case "advance", "pre":
			return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
		case "delay", "post":
			return false
		default:
			return false
		}
	}
	return true
}

func (s *PlanService) VersionUp(plan domain.CollectionPlan) domain.CollectionPlan {
	plan.Version++
	return plan
}

func (s *PlanService) BuildRules(plan domain.CollectionPlan) []domain.PlanRule {
	rules := append([]domain.PlanRule(nil), plan.Rules...)
	sort.Slice(rules, func(i, j int) bool {
		if rules[i].RuleType == rules[j].RuleType {
			return rules[i].RuleValue < rules[j].RuleValue
		}
		return rules[i].RuleType < rules[j].RuleType
	})
	return rules
}

func (s *PlanService) Activate(plan *domain.CollectionPlan) error {
	return ValidateStatusTransition(plan.Status, domain.PlanActive, BuildStatusFlow())
}

func (s *PlanService) Disable(plan *domain.CollectionPlan) error {
	return ValidateStatusTransition(plan.Status, domain.PlanDisabled, BuildStatusFlow())
}

func (s *PlanService) Archive(plan *domain.CollectionPlan) error {
	return ValidateStatusTransition(plan.Status, domain.PlanArchived, BuildStatusFlow())
}

func (s *PlanService) EstimateMonthly(plan domain.CollectionPlan, frequency int) float64 {
	base := plan.EstimatedWeightTons
	if frequency <= 0 {
		frequency = 1
	}
	return base * float64(frequency) * 30
}

func (s *PlanService) Summarize(plan domain.CollectionPlan) string {
	return fmt.Sprintf("%s-v%d", plan.Name, plan.Version)
}

