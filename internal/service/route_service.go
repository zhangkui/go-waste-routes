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
		current = current.Add(time.Duration(max(ordered[index].StayMinutes, 0)) * time.Minute)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

