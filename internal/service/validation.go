package service

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	upperRegexp   = regexp.MustCompile(`[A-Z]`)
	digitRegexp   = regexp.MustCompile(`[0-9]`)
	specialRegexp = regexp.MustCompile(`[^A-Za-z0-9]`)
	emailRegexp   = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	plateRegexp   = regexp.MustCompile(`^[A-Z]{1}[A-Z0-9·-]{5,7}$`)
)

func ValidatePasswordPolicy(password string, minLength int, requireUpper, requireDigit, requireSpecial bool) error {
	if len(password) < minLength {
		return fmt.Errorf("password too short")
	}
	if requireUpper && !upperRegexp.MatchString(password) {
		return fmt.Errorf("password must contain upper-case letters")
	}
	if requireDigit && !digitRegexp.MatchString(password) {
		return fmt.Errorf("password must contain digits")
	}
	if requireSpecial && !specialRegexp.MatchString(password) {
		return fmt.Errorf("password must contain special characters")
	}
	return nil
}

func NormalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func ValidateStatusTransition(current string, next string, flow map[string][]string) error {
	current = NormalizeStatus(current)
	next = NormalizeStatus(next)
	allowed := flow[current]
	for _, candidate := range allowed {
		if candidate == next {
			return nil
		}
	}
	return fmt.Errorf("invalid status transition from %s to %s", current, next)
}

func ContainsString(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(value, needle) {
			return true
		}
	}
	return false
}

func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func NonEmpty(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func BuildStatusFlow() map[string][]string {
	return map[string][]string{
		"draft":    {"active", "archived"},
		"active":   {"disabled", "archived"},
		"disabled": {"active", "archived"},
		"pending":  {"running", "cancelled", "abnormal", "completed"},
		"running":  {"completed", "cancelled", "abnormal"},
		"completed": {"written_off"},
		"confirmed": {"partial", "paid", "written_off"},
		"partial":  {"paid", "written_off"},
	}
}

func ValidateEmailAddress(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if !emailRegexp.MatchString(value) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func ValidatePhoneNumber(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	digits := 0
	for _, r := range value {
		if r >= '0' && r <= '9' {
			digits++
		}
	}
	if digits < 7 || digits > 15 {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}

func ValidatePlateNumber(value string) error {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return fmt.Errorf("plate number required")
	}
	if !plateRegexp.MatchString(value) {
		return fmt.Errorf("invalid plate number")
	}
	return nil
}

func ValidateMoneyAmount(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("invalid money amount")
	}
	if value < 0 {
		return fmt.Errorf("money amount cannot be negative")
	}
	return nil
}

func ValidateNonNegativeAmount(name string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative", name)
	}
	return nil
}

func ValidateWeightRange(name string, value, max float64) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative", name)
	}
	if max > 0 && value > max {
		return fmt.Errorf("%s exceeds max %.2f", name, max)
	}
	return nil
}

func ValidatePercentage(name string, value float64) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("%s must be between 0 and 100", name)
	}
	return nil
}

func ParseSeparatedInts(value string) ([]int, error) {
	tokens := splitTokens(value)
	values := make([]int, 0, len(tokens))
	for _, token := range tokens {
		number, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		values = append(values, number)
	}
	return values, nil
}

func ParseSeparatedFloats(value string) ([]float64, error) {
	tokens := splitTokens(value)
	values := make([]float64, 0, len(tokens))
	for _, token := range tokens {
		number, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, number)
	}
	return values, nil
}

func ValidateRouteCapacity(estimatedWeightTons, ratedLoadTons float64) error {
	if estimatedWeightTons < 0 {
		return fmt.Errorf("estimated weight cannot be negative")
	}
	if ratedLoadTons <= 0 {
		return fmt.Errorf("rated load must be positive")
	}
	if estimatedWeightTons > ratedLoadTons*0.80 {
		return fmt.Errorf("route exceeds safety threshold")
	}
	return nil
}

func ValidateServiceWindow(value string) error {
	window := normalizeServiceWindow(value)
	if window == "" {
		return fmt.Errorf("invalid service window")
	}
	return nil
}

func ValidateServiceWindows(values []string) error {
	for _, value := range values {
		if err := ValidateServiceWindow(value); err != nil {
			return err
		}
	}
	return nil
}

func ValidatePlanRuleTypeAndValue(ruleType, ruleValue string) error {
	ruleType = strings.ToLower(strings.TrimSpace(ruleType))
	switch ruleType {
	case "daily":
		return nil
	case "weekly":
		if len(ParseWeekdayTokensOrZero(ruleValue)) == 0 {
			return fmt.Errorf("weekly rule requires weekday value")
		}
	case "monthly":
		if len(ParseMonthDayTokensOrZero(ruleValue)) == 0 {
			return fmt.Errorf("monthly rule requires day value")
		}
	default:
		return fmt.Errorf("unknown rule type %s", ruleType)
	}
	return nil
}

func ValidateDateRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return fmt.Errorf("date range required")
	}
	if end.Before(start) {
		return fmt.Errorf("end must not be before start")
	}
	return nil
}

func ValidateWeightSplit(weights []float64) error {
	total := 0.0
	for _, weight := range weights {
		if weight < 0 {
			return fmt.Errorf("weight split cannot contain negative values")
		}
		total += weight
	}
	if total == 0 {
		return fmt.Errorf("weight split cannot be empty")
	}
	if math.Abs(total-100) > 0.5 {
		return fmt.Errorf("weight split should sum to 100")
	}
	return nil
}

func ValidateStatusPath(entity string, current string, next string, flow map[string][]string) error {
	current = NormalizeStatus(current)
	next = NormalizeStatus(next)
	for _, candidate := range flow[current] {
		if candidate == next {
			return nil
		}
	}
	return fmt.Errorf("invalid %s transition from %s to %s", entity, current, next)
}

func BuildValidationErrors(errs ...error) error {
	filtered := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		filtered = append(filtered, err.Error())
	}
	if len(filtered) == 0 {
		return nil
	}
	sort.Strings(filtered)
	return errors.New(strings.Join(filtered, "; "))
}

func NormalizeList(values []string) []string {
	items := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		items = append(items, normalized)
	}
	sort.Strings(items)
	return items
}

func ParseWeekdayTokensOrZero(value string) []time.Weekday {
	weekdays := make([]time.Weekday, 0)
	for _, token := range splitTokens(value) {
		switch strings.ToLower(token) {
		case "sun", "sunday", "0":
			weekdays = append(weekdays, time.Sunday)
		case "mon", "monday", "1":
			weekdays = append(weekdays, time.Monday)
		case "tue", "tuesday", "2":
			weekdays = append(weekdays, time.Tuesday)
		case "wed", "wednesday", "3":
			weekdays = append(weekdays, time.Wednesday)
		case "thu", "thursday", "4":
			weekdays = append(weekdays, time.Thursday)
		case "fri", "friday", "5":
			weekdays = append(weekdays, time.Friday)
		case "sat", "saturday", "6":
			weekdays = append(weekdays, time.Saturday)
		}
	}
	return weekdays
}

func ParseMonthDayTokensOrZero(value string) []int {
	days := make([]int, 0)
	for _, token := range splitTokens(value) {
		day, err := strconv.Atoi(token)
		if err != nil || day < 1 || day > 31 {
			continue
		}
		days = append(days, day)
	}
	sort.Ints(days)
	return days
}

func normalizeServiceWindow(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Count(value, "-") != 1 {
		return ""
	}
	parts := strings.Split(value, "-")
	start := normalizeClock(parts[0])
	end := normalizeClock(parts[1])
	if start == "" || end == "" {
		return ""
	}
	return start + "-" + end
}
