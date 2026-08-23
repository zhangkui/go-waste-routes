package service

import (
	"errors"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type RouteService struct{}

func NewRouteService() *RouteService { return &RouteService{} }

func (s *RouteService) ValidateCapacity(route domain.Route, vehicle domain.Vehicle, stops []domain.RouteStop) error {
	total := 0.0
	for _, stop := range stops {
		total += stop.EstimatedWeightTons
	}
	if total > vehicle.RatedLoadTons*0.80 {
		return errors.New("route capacity exceeded")
	}
	return nil
}

func (s *RouteService) SortStops(stops []domain.RouteStop) []domain.RouteStop {
	ordered := append([]domain.RouteStop(nil), stops...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	return ordered
}

func (s *RouteService) EstimateArrivalTimes(start time.Time, stops []domain.RouteStop) []domain.RouteStop {
	ordered := s.SortStops(stops)
	current := start
	for index := range ordered {
		ordered[index].EstimatedArrival = current.Format("15:04")
		current = current.Add(time.Duration(max(ordered[index].StayMinutes, 1)) * time.Minute)
	}
	return ordered
}

func (s *RouteService) NextStatus(status string) string {
	switch NormalizeStatus(status) {
	case domain.RoutePending:
		return domain.RouteRunning
	case domain.RouteRunning:
		return domain.RouteCompleted
	default:
		return status
	}
}

func (s *RouteService) ExecuteRoute(route *domain.Route, vehicle *domain.Vehicle, execute func() error) error {
	s.acquireRoute(route, vehicle)
	if err := execute(); err != nil {
		s.finishAbnormal(route, vehicle)
		return err
	}
	s.finishComplete(route, vehicle)
	return nil
}

func (s *RouteService) CancelRoute(route *domain.Route, vehicle *domain.Vehicle) {
	s.finishCancelled(route, vehicle)
}

func (s *RouteService) acquireRoute(route *domain.Route, vehicle *domain.Vehicle) {
	SetVehicleDispatchState(vehicle, "running")
	route.Status = domain.RouteRunning
}

func (s *RouteService) finishComplete(route *domain.Route, vehicle *domain.Vehicle) {
	route.Status = domain.RouteCompleted
	releaseVehicleDispatch(vehicle)
}

func (s *RouteService) finishCancelled(route *domain.Route, vehicle *domain.Vehicle) {
	route.Status = domain.RouteCancelled
	releaseVehicleDispatch(vehicle)
}

func (s *RouteService) finishAbnormal(route *domain.Route, vehicle *domain.Vehicle) {
	route.Status = domain.RouteAbnormal
	releaseVehicleDispatch(vehicle)
}

// releaseVehicleDispatch restores the vehicle to the idle dispatch state so it
// can be rescheduled. It is the single cleanup point shared by the complete,
// cancel, and abnormal exit paths so resource release stays consistent.
func releaseVehicleDispatch(vehicle *domain.Vehicle) {
	if vehicle == nil {
		return
	}
	SetVehicleDispatchState(vehicle, "available")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
