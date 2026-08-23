package verify

import (
	"testing"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/service"
)

func TestBug040_ConcurrentMileageUpdatesAccumulate(t *testing.T) {
	tasks := []domain.Task{{MileageKm: 3}, {MileageKm: 5}, {MileageKm: 7}}
	vehicle := &domain.Vehicle{CurrentMileageKm: 100}
	deltas := service.NewTaskService().MileageDeltas(tasks)
	service.NewRouteService().AccumulateMileage(vehicle, deltas)
	if vehicle.CurrentMileageKm != 115 {
		t.Fatalf("mileage=%v,want 115", vehicle.CurrentMileageKm)
	}
}
