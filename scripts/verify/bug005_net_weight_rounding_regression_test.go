package verify

import (
	"math"
	"testing"

	"go-waste-routes/internal/service"
)

func TestBug005_NetWeightIsRoundedToTwoDecimals(t *testing.T) {
	got := service.NewWeighingService().CalculateNetWeight(10.015, 9.01)
	if math.Abs(got-1.01) > 1e-9 {
		t.Fatalf("net weight = %.3f, want 1.01", got)
	}
}
