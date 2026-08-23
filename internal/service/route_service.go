package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type RouteService struct{}

type RouteStopWriter interface {
	SaveStops(context.Context, int64, []domain.RouteStop) error
}

func NewRouteService() *RouteService { return &RouteService{} }

func (s *RouteService) SaveRouteStops(ctx context.Context, route *domain.Route, stops []domain.RouteStop, writer RouteStopWriter) error {
	if route == nil || writer == nil {
		return errors.New("route and writer required")
	}
	route.Stops = append(route.Stops, stops...)
	for index, stop := range stops {
		if stop.CustomerID <= 0 || stop.Sequence != index+1 {
			return errors.New("invalid route stop sequence")
		}
	}
	return writer.SaveStops(ctx, route.ID, stops)
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
