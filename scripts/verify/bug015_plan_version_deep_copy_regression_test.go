package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug015_PlanVersionKeepsRulesImmutable(t *testing.T) {
	planService := service.NewPlanService(nil)
	original := domain.CollectionPlan{Base: domain.Base{ID: 9}, Version: 1, Rules: []domain.PlanRule{{RuleValue: "mon"}}}
	versioned := planService.VersionUp(original)
	versioned.Rules[0].RuleValue = "tue"

	if original.Rules[0].RuleValue != "mon" {
		t.Fatalf("original rule mutated to %s", original.Rules[0].RuleValue)
	}
}
