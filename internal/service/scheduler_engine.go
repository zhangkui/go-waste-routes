package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type TimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ScheduleOccurrence struct {
	Date        time.Time `json:"date"`
	RuleType    string    `json:"rule_type"`
	RuleValue   string    `json:"rule_value"`
	Shifted     bool      `json:"shifted"`
	ShiftReason string    `json:"shift_reason"`
	Windows     []TimeWindow `json:"windows,omitempty"`
}

type SchedulePreview struct {
	TaskNumber string             `json:"task_number"`
	PlanID     int64              `json:"plan_id"`
	CustomerID int64              `json:"customer_id"`
	Occurrence ScheduleOccurrence `json:"occurrence"`
}

type Scheduler struct {
	config *ConfigService
}

func NewScheduler(config *ConfigService) *Scheduler {
	return &Scheduler{config: config}
}

func (s *Scheduler) BuildOccurrences(plan domain.CollectionPlan, rules []domain.PlanRule, holidays []domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	if end.Before(start) {
		return nil
	}
	holidayMap := buildHolidayMap(holidays)
	occurrences := make([]ScheduleOccurrence, 0)
	for _, rule := range rules {
		if !rule.Enabled || (rule.PlanID != 0 && rule.PlanID != plan.ID) {
			continue
		}
		expansions := s.expandRule(plan, rule, holidayMap, start, end)
		occurrences = append(occurrences, expansions...)
	}
	occurrences = uniqueOccurrences(occurrences)
	sort.SliceStable(occurrences, func(i, j int) bool {
		if occurrences[i].Date.Equal(occurrences[j].Date) {
			return occurrences[i].RuleType < occurrences[j].RuleType
		}
		return occurrences[i].Date.Before(occurrences[j].Date)
	})
	return occurrences
}

func (s *Scheduler) expandRule(plan domain.CollectionPlan, rule domain.PlanRule, holidayMap map[string]domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	switch strings.ToLower(strings.TrimSpace(rule.RuleType)) {
	case "daily":
		return s.expandDaily(rule, holidayMap, start, end)
	case "weekly":
		return s.expandWeekly(rule, holidayMap, start, end)
	case "monthly":
		return s.expandMonthly(rule, holidayMap, start, end)
	default:
		return s.expandFallback(plan, rule, holidayMap, start, end)
	}
}

func (s *Scheduler) expandDaily(rule domain.PlanRule, holidayMap map[string]domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	occurrences := make([]ScheduleOccurrence, 0)
	for day := truncateDate(start); !day.After(end); day = day.AddDate(0, 0, 1) {
		date, shifted, reason, ok := s.adjustForHoliday(day, rule.RuleValue, holidayMap)
		if !ok {
			continue
		}
		occurrences = append(occurrences, ScheduleOccurrence{
			Date:        date,
			RuleType:    "daily",
			RuleValue:   rule.RuleValue,
			Shifted:     shifted,
			ShiftReason: reason,
			Windows:     parseWindows(rule.StartTime, rule.EndTime),
		})
	}
	return occurrences
}

func (s *Scheduler) expandWeekly(rule domain.PlanRule, holidayMap map[string]domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	weekdays := parseWeekdays(rule.RuleValue)
	if len(weekdays) == 0 {
		weekdays = []time.Weekday{truncateDate(start).Weekday()}
	}
	allowed := make(map[time.Weekday]struct{}, len(weekdays))
	for _, weekday := range weekdays {
		allowed[weekday] = struct{}{}
	}
	occurrences := make([]ScheduleOccurrence, 0)
	for day := truncateDate(start); !day.After(end); day = day.AddDate(0, 0, 1) {
		if _, ok := allowed[day.Weekday()]; !ok {
			continue
		}
		date, shifted, reason, ok := s.adjustForHoliday(day, rule.RuleValue, holidayMap)
		if !ok {
			continue
		}
		occurrences = append(occurrences, ScheduleOccurrence{
			Date:        date,
			RuleType:    "weekly",
			RuleValue:   rule.RuleValue,
			Shifted:     shifted,
			ShiftReason: reason,
			Windows:     parseWindows(rule.StartTime, rule.EndTime),
		})
	}
	return occurrences
}

func (s *Scheduler) expandMonthly(rule domain.PlanRule, holidayMap map[string]domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	days := parseMonthDays(rule.RuleValue)
	if len(days) == 0 {
		days = []int{truncateDate(start).Day()}
	}
	allowed := make(map[int]struct{}, len(days))
	for _, day := range days {
		allowed[day] = struct{}{}
	}
	occurrences := make([]ScheduleOccurrence, 0)
	for day := truncateDate(start); !day.After(end); day = day.AddDate(0, 0, 1) {
		if _, ok := allowed[day.Day()]; !ok {
			continue
		}
		date, shifted, reason, ok := s.adjustForHoliday(day, rule.RuleValue, holidayMap)
		if !ok {
			continue
		}
		occurrences = append(occurrences, ScheduleOccurrence{
			Date:        date,
			RuleType:    "monthly",
			RuleValue:   rule.RuleValue,
			Shifted:     shifted,
			ShiftReason: reason,
			Windows:     parseWindows(rule.StartTime, rule.EndTime),
		})
	}
	return occurrences
}

func (s *Scheduler) expandFallback(plan domain.CollectionPlan, rule domain.PlanRule, holidayMap map[string]domain.HolidayConfig, start, end time.Time) []ScheduleOccurrence {
	if strings.TrimSpace(rule.RuleValue) == "" {
		return s.expandDaily(rule, holidayMap, start, end)
	}
	if strings.ContainsAny(rule.RuleValue, "一二三四五六日天周mon") {
		return s.expandWeekly(rule, holidayMap, start, end)
	}
	return s.expandMonthly(rule, holidayMap, start, end)
}

func (s *Scheduler) Preview(plan domain.CollectionPlan, customer domain.Customer, rules []domain.PlanRule, holidays []domain.HolidayConfig, start, end time.Time) []SchedulePreview {
	occurrences := s.BuildOccurrences(plan, rules, holidays, start, end)
	previews := make([]SchedulePreview, 0, len(occurrences))
	for index, occurrence := range occurrences {
		previews = append(previews, SchedulePreview{
			TaskNumber: GenerateTaskNumber(occurrence.Date, index+1),
			PlanID:     plan.ID,
			CustomerID: customer.ID,
			Occurrence: occurrence,
		})
	}
	return previews
}

func (s *Scheduler) ServiceDays(customer domain.Customer) map[time.Weekday]struct{} {
	weekdays := parseWeekdays(customer.ServiceWeekdays)
	if len(weekdays) == 0 {
		return map[time.Weekday]struct{}{}
	}
	allowed := make(map[time.Weekday]struct{}, len(weekdays))
	for _, weekday := range weekdays {
		allowed[weekday] = struct{}{}
	}
	return allowed
}

func (s *Scheduler) ServiceWindows(customer domain.Customer) []TimeWindow {
	return parseWindows(customer.ServiceWindows, "")
}

func (s *Scheduler) adjustForHoliday(day time.Time, strategy string, holidayMap map[string]domain.HolidayConfig) (time.Time, bool, string, bool) {
	normalized := truncateDate(day)
	if config, ok := holidayMap[dateKey(normalized)]; ok && config.Enabled {
		if config.IsWorkday {
			return normalized, false, "holiday_override_workday", true
		}
		switch strings.ToLower(strings.TrimSpace(strategy)) {
		case "skip":
			return time.Time{}, false, "holiday_skip", false
		case "advance", "previous_workday":
			for candidate := normalized.AddDate(0, 0, -1); !candidate.Before(normalized.AddDate(0, 0, -14)); candidate = candidate.AddDate(0, 0, -1) {
				if s.isWorkday(candidate, holidayMap) {
					return candidate, true, "holiday_advance", true
				}
			}
			return time.Time{}, false, "holiday_advance_failed", false
		default:
			for candidate := normalized.AddDate(0, 0, 1); !candidate.After(normalized.AddDate(0, 0, 14)); candidate = candidate.AddDate(0, 0, 1) {
				if s.isWorkday(candidate, holidayMap) {
					return candidate, true, "holiday_postpone", true
				}
			}
			return time.Time{}, false, "holiday_postpone_failed", false
		}
	}
	if !s.isWorkday(normalized, holidayMap) {
		switch strings.ToLower(strings.TrimSpace(strategy)) {
		case "skip":
			return time.Time{}, false, "weekend_skip", false
		case "advance", "previous_workday":
			for candidate := normalized.AddDate(0, 0, -1); !candidate.Before(normalized.AddDate(0, 0, -14)); candidate = candidate.AddDate(0, 0, -1) {
				if s.isWorkday(candidate, holidayMap) {
					return candidate, true, "weekend_advance", true
				}
			}
			return time.Time{}, false, "weekend_advance_failed", false
		default:
			for candidate := normalized.AddDate(0, 0, 1); !candidate.After(normalized.AddDate(0, 0, 14)); candidate = candidate.AddDate(0, 0, 1) {
				if s.isWorkday(candidate, holidayMap) {
					return candidate, true, "weekend_postpone", true
				}
			}
			return time.Time{}, false, "weekend_postpone_failed", false
		}
	}
	return normalized, false, "", true
}

func (s *Scheduler) isWorkday(day time.Time, holidayMap map[string]domain.HolidayConfig) bool {
	if s.config != nil {
		if s.config.IsWorkday(configAsHolidays(holidayMap), day) {
			return true
		}
		return false
	}
	if config, ok := holidayMap[dateKey(day)]; ok && config.Enabled {
		return config.IsWorkday
	}
	return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
}

func configAsHolidays(holidayMap map[string]domain.HolidayConfig) []domain.HolidayConfig {
	values := make([]domain.HolidayConfig, 0, len(holidayMap))
	for _, value := range holidayMap {
		values = append(values, value)
	}
	return values
}

func buildHolidayMap(holidays []domain.HolidayConfig) map[string]domain.HolidayConfig {
	result := make(map[string]domain.HolidayConfig, len(holidays))
	for _, holiday := range holidays {
		result[dateKey(holiday.HolidayDate)] = holiday
	}
	return result
}

func uniqueOccurrences(items []ScheduleOccurrence) []ScheduleOccurrence {
	seen := make(map[string]ScheduleOccurrence, len(items))
	for _, item := range items {
		key := dateKey(item.Date) + "|" + strings.ToLower(item.RuleType) + "|" + strings.TrimSpace(item.RuleValue)
		if existing, ok := seen[key]; ok {
			if existing.Shifted && !item.Shifted {
				seen[key] = item
			}
			continue
		}
		seen[key] = item
	}
	ordered := make([]ScheduleOccurrence, 0, len(seen))
	for _, item := range seen {
		ordered = append(ordered, item)
	}
	return ordered
}

func parseWeekdays(value string) []time.Weekday {
	tokens := splitTokens(value)
	weekdays := make([]time.Weekday, 0, len(tokens))
	for _, token := range tokens {
		switch strings.ToLower(token) {
		case "0", "sun", "sunday", "日", "周日", "周天":
			weekdays = append(weekdays, time.Sunday)
		case "1", "mon", "monday", "周一":
			weekdays = append(weekdays, time.Monday)
		case "2", "tue", "tuesday", "周二":
			weekdays = append(weekdays, time.Tuesday)
		case "3", "wed", "wednesday", "周三":
			weekdays = append(weekdays, time.Wednesday)
		case "4", "thu", "thursday", "周四":
			weekdays = append(weekdays, time.Thursday)
		case "5", "fri", "friday", "周五":
			weekdays = append(weekdays, time.Friday)
		case "6", "sat", "saturday", "周六":
			weekdays = append(weekdays, time.Saturday)
		default:
			if number, err := strconv.Atoi(token); err == nil {
				switch number {
				case 0:
					weekdays = append(weekdays, time.Sunday)
				case 1:
					weekdays = append(weekdays, time.Monday)
				case 2:
					weekdays = append(weekdays, time.Tuesday)
				case 3:
					weekdays = append(weekdays, time.Wednesday)
				case 4:
					weekdays = append(weekdays, time.Thursday)
				case 5:
					weekdays = append(weekdays, time.Friday)
				case 6:
					weekdays = append(weekdays, time.Saturday)
				}
			}
		}
	}
	return dedupeWeekdays(weekdays)
}

func parseMonthDays(value string) []int {
	tokens := splitTokens(value)
	days := make([]int, 0, len(tokens))
	for _, token := range tokens {
		day, err := strconv.Atoi(token)
		if err != nil || day < 1 || day > 31 {
			continue
		}
		days = append(days, day)
	}
	sort.Ints(days)
	return dedupeInts(days)
}

func parseWindows(values ...string) []TimeWindow {
	windows := make([]TimeWindow, 0)
	for _, value := range values {
		for _, token := range splitTokens(value) {
			window := normalizeWindow(token)
			if window.Start == "" || window.End == "" {
				continue
			}
			windows = append(windows, window)
		}
	}
	return windows
}

func normalizeWindow(value string) TimeWindow {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return TimeWindow{}
	}
	start := normalizeClock(parts[0])
	end := normalizeClock(parts[1])
	if start == "" || end == "" {
		return TimeWindow{}
	}
	return TimeWindow{Start: start, End: end}
}

func normalizeClock(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) == 5 && strings.Count(value, ":") == 1 {
		return value
	}
	if parsed, err := time.Parse("15:04", value); err == nil {
		return parsed.Format("15:04")
	}
	return ""
}

func splitTokens(value string) []string {
	value = strings.ReplaceAll(value, ";", ",")
	value = strings.ReplaceAll(value, "，", ",")
	value = strings.ReplaceAll(value, " ", ",")
	parts := strings.Split(value, ",")
	tokens := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			tokens = append(tokens, part)
		}
	}
	return tokens
}

func dedupeWeekdays(values []time.Weekday) []time.Weekday {
	seen := make(map[time.Weekday]struct{}, len(values))
	unique := make([]time.Weekday, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	sort.Slice(unique, func(i, j int) bool { return unique[i] < unique[j] })
	return unique
}

func dedupeInts(values []int) []int {
	seen := make(map[int]struct{}, len(values))
	unique := make([]int, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	sort.Ints(unique)
	return unique
}

func GenerateTaskNumber(date time.Time, sequence int) string {
	return fmt.Sprintf("TASK-%s-%04d", truncateDate(date).Format("20060102"), sequence)
}

func GenerateRouteNumber(planID int64, date time.Time, sequence int) string {
	return fmt.Sprintf("ROUTE-%d-%s-%03d", planID, truncateDate(date).Format("20060102"), sequence)
}

func truncateDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func dateKey(value time.Time) string {
	return truncateDate(value).Format("2006-01-02")
}

