package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug001_HolidayAdvancePreviewMatchesScheduler(t *testing.T) {
	config := service.NewConfigService()
	holiday := domain.HolidayConfig{
		HolidayDate: time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC),
		Enabled:     true,
		IsWorkday:   false,
	}
	holidayConfigs := []domain.HolidayConfig{holiday}

	planService := service.NewPlanService(func(day time.Time) bool {
		return !config.IsWorkday(holidayConfigs, day)
	})
	scheduler := service.NewScheduler(config)

	plan := domain.CollectionPlan{
		Base:            domain.Base{ID: 42},
		HolidayStrategy:  "advance",
		EffectiveFrom:    time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC),
		EstimatedWeightTons: 1.5,
	}
	rules := []domain.PlanRule{{
		PlanID:   42,
		RuleType: "daily",
		RuleValue: "advance",
		Enabled:  true,
	}}

	planPreview := planService.Preview(plan, 1)
	if len(planPreview) != 1 {
		t.Fatalf("plan preview length = %d, want 1", len(planPreview))
	}

	schedulePreview := scheduler.Preview(plan, domain.Customer{Base: domain.Base{ID: 7}}, rules, holidayConfigs, plan.EffectiveFrom, plan.EffectiveFrom.AddDate(0, 0, 4))
	if len(schedulePreview) == 0 {
		t.Fatal("scheduler preview is empty")
	}

	if !planPreview[0].Equal(schedulePreview[0].Occurrence.Date) {
		t.Fatalf("preview date = %s, schedule date = %s", planPreview[0].Format("2006-01-02"), schedulePreview[0].Occurrence.Date.Format("2006-01-02"))
	}
}
