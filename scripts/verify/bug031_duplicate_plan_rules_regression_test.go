package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug031_DuplicatePlanRulesProduceOneTaskPerDay(t *testing.T) {
	day := time.Date(2026, 8, 23, 0, 0, 0, 0, time.UTC)
	plan := domain.CollectionPlan{Base: domain.Base{ID: 31}}
	rules := []domain.PlanRule{{PlanID: 31, RuleType: "daily", Enabled: true}, {PlanID: 31, RuleType: "daily", Enabled: true}}
	items := service.NewScheduler(nil).BuildOccurrences(plan, rules, nil, day, day)
	if len(items) != 1 {
		t.Fatalf("duplicate rules produced %d occurrences, want 1", len(items))
	}
}
