package verify

import (
	"context"
	"errors"
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

type bug028HistoryReader struct{}

func (bug028HistoryReader) HistoryForCustomer(context.Context, int64) ([]domain.WeighingRecord, error) {
	return nil, errors.New("history database unavailable")
}

func TestBug028_HistoryReadErrorIsNotTreatedAsEmpty(t *testing.T) {
	record := domain.WeighingRecord{NetWeight: 12, GrossWeight: 14, TareWeight: 2}
	_, err := service.NewWeighingService().DetectWithHistory(context.Background(), bug028HistoryReader{}, 7, record, 20)
	if err == nil {
		t.Fatal("history read error was treated as empty history")
	}
}
