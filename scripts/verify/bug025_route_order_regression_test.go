package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
	"go-waste-routes/internal/transport/http/handler"
)

func TestBug025_RoutePreviewAndExecutionRetainDuplicateStops(t *testing.T) {
	stops := []domain.RouteStop{
		{CustomerID: 10, Sequence: 1},
		{CustomerID: 20, Sequence: 1},
		{CustomerID: 30, Sequence: 2},
	}
	presentation := handler.NewRoutePresentationHandler(service.NewRouteService(), service.NewTaskService()).Build(stops)
	if len(presentation.Preview) != len(presentation.Execution) {
		t.Fatalf("route preview and execution contain different stop counts: preview=%d execution=%d", len(presentation.Preview), len(presentation.Execution))
	}
	for index, previewStop := range presentation.Preview {
		executionStop := presentation.Execution[index]
		if previewStop.CustomerID != executionStop.CustomerID || previewStop.Sequence != executionStop.Sequence {
			t.Fatalf("route order differs at index %d: preview=%+v execution=%+v", index, previewStop, executionStop)
		}
	}
}