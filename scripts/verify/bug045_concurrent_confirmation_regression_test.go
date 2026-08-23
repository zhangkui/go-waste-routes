package verify

import (
	"testing"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug045_ConcurrentWeighingConfirmationIsSingleWinner(t *testing.T) {
	record := &domain.WeighingRecord{Status: "reviewing"}
	service.NewWeighingService().ConfirmCandidates(record, []int64{41, 82}, time.Now())
	if record.Status != "confirmed" || record.ConfirmedBy == nil {
		t.Fatalf("confirmation did not produce a single persisted winner: %#v", record)
	}
}
