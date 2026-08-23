package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug050_StatusHistoryKeepsSameMomentOrder(t *testing.T) {
	moment := time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC)
	incoming := []domain.StatusHistory{
		{EntityID: 17, ToStatus: "claimed", CreatedAt: moment},
		{EntityID: 17, ToStatus: "running", CreatedAt: moment},
		{EntityID: 17, ToStatus: "completed", CreatedAt: moment},
	}
	merged := service.MergeStatusHistory(nil, incoming)
	for index, want := range []string{"claimed", "running", "completed"} {
		if merged[index].ToStatus != want {
			t.Fatalf("same-moment history[%d]=%q,want %q", index, merged[index].ToStatus, want)
		}
	}
}
