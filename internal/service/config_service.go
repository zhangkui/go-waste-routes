package service

import (
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type ConfigService struct{}

func NewConfigService() *ConfigService { return &ConfigService{} }

func (s *ConfigService) IsHoliday(configs []domain.HolidayConfig, day time.Time) bool {
	for _, config := range configs {
		if sameDate(config.HolidayDate, day) {
			return config.Enabled && !config.IsWorkday
		}
	}
	return false
}

func (s *ConfigService) IsWorkday(configs []domain.HolidayConfig, day time.Time) bool {
	for _, config := range configs {
		if sameDate(config.HolidayDate, day) {
			return config.Enabled && config.IsWorkday
		}
	}
	return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
}

func (s *ConfigService) BillingMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "per_trip", "weight", "mixed":
		return value
	default:
		return "weight"
	}
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

