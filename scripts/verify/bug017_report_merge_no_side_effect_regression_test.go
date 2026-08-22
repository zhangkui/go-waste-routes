package verify

import (
	"testing"

	"go-waste-routes/internal/service"
)

func TestBug017_ReportMergeDoesNotMutateBaseSeries(t *testing.T) {
	reportService := service.NewReportService()
	base := service.DashboardStats{CollectionSeries: map[string]float64{"2026-08-22": 1}}
	extra := service.DashboardStats{CollectionSeries: map[string]float64{"2026-08-22": 2}}

	merged := reportService.MergeStats(base, extra)
	if base.CollectionSeries["2026-08-22"] != 1 {
		t.Fatalf("base series mutated to %.2f", base.CollectionSeries["2026-08-22"])
	}
	if merged.CollectionSeries["2026-08-22"] != 3 {
		t.Fatalf("merged series = %.2f, want 3", merged.CollectionSeries["2026-08-22"])
	}
}
