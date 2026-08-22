package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type BillingCycle string

const (
	BillingCycleMonth   BillingCycle = "month"
	BillingCycleQuarter BillingCycle = "quarter"
)

type RateCard struct {
	BaseFee            float64
	PerTripFee         float64
	PerTonFee          float64
	HolidayFee         float64
	DistanceFee        float64
	OverLimitFee       float64
	OverLimitThreshold float64
	MinimumCharge      float64
}

type ChargeItem struct {
	Code        string
	Description string
	Amount      float64
	TaskID      int64
	CustomerID  int64
	CategoryID  int64
}

type BillingEngine struct {
	config *ConfigService
}

func NewBillingEngine(config *ConfigService) *BillingEngine {
	return &BillingEngine{config: config}
}

func (e *BillingEngine) BuildInvoiceNumber(customerID int64, periodStart, periodEnd time.Time) string {
	return fmt.Sprintf("INV-%d-%s-%s", customerID, periodStart.Format("20060102"), periodEnd.Format("20060102"))
}

func (e *BillingEngine) Period(reference time.Time, cycle BillingCycle) (time.Time, time.Time) {
	reference = truncateDate(reference)
	switch cycle {
	case BillingCycleQuarter:
		quarter := (int(reference.Month())-1)/3 + 1
		firstMonth := time.Month((quarter-1)*3 + 1)
		start := time.Date(reference.Year(), firstMonth, 1, 0, 0, 0, 0, reference.Location())
		end := start.AddDate(0, 3, -1)
		return start, truncateDate(end)
	default:
		start := time.Date(reference.Year(), reference.Month(), 1, 0, 0, 0, 0, reference.Location())
		end := start.AddDate(0, 1, -1)
		return start, truncateDate(end)
	}
}

func (e *BillingEngine) PreviousPeriod(reference time.Time, cycle BillingCycle) (time.Time, time.Time) {
	start, end := e.Period(reference.AddDate(0, -1, 0), cycle)
	if cycle == BillingCycleQuarter {
		return start, end
	}
	return start, end
}

func (e *BillingEngine) BuildCharges(customer domain.Customer, tasks []domain.Task, categories map[int64]domain.WasteCategory, weights map[int64]float64, rates map[int64]float64, rateCard RateCard, mode string) []ChargeItem {
	mode = NormalizeBillingMode(mode)
	charges := make([]ChargeItem, 0)
	tripCount := 0
	for _, task := range tasks {
		tripCount++
		charges = append(charges, ChargeItem{
			Code:        "base_trip",
			Description: fmt.Sprintf("Task %s base charge", task.TaskNumber),
			Amount:      roundMoney(rateCard.PerTripFee),
			TaskID:      task.ID,
			CustomerID:  customer.ID,
		})
		netWeight := weights[task.ID]
		if netWeight < 0 {
			netWeight = 0
		}
		categoryID := task.PlanID
		if task.RouteID != nil {
			categoryID = *task.RouteID
		}
		unitPrice := rateCard.PerTonFee
		if override, ok := rates[categoryID]; ok && override > 0 {
			unitPrice = override
		}
		if category, ok := categories[categoryID]; ok {
			charges = append(charges, ChargeItem{
				Code:        "waste_charge",
				Description: category.Name,
				Amount:      roundMoney(netWeight * unitPrice),
				TaskID:      task.ID,
				CustomerID:  customer.ID,
				CategoryID:  category.ID,
			})
		} else {
			charges = append(charges, ChargeItem{
				Code:        "waste_charge",
				Description: "unknown category",
				Amount:      roundMoney(netWeight * unitPrice),
				TaskID:      task.ID,
				CustomerID:  customer.ID,
			})
		}
	}
	if mode == "mixed" && rateCard.BaseFee > 0 {
		charges = append(charges, ChargeItem{
			Code:        "base_service",
			Description: "Base service fee",
			Amount:      roundMoney(rateCard.BaseFee),
			CustomerID:  customer.ID,
		})
	}
	if tripCount > 0 && rateCard.MinimumCharge > 0 {
		total := e.SumCharges(charges)
		if total < rateCard.MinimumCharge {
			charges = append(charges, ChargeItem{
				Code:        "minimum_charge",
				Description: "Minimum billing floor",
				Amount:      roundMoney(rateCard.MinimumCharge - total),
				CustomerID:  customer.ID,
			})
		}
	}
	return charges
}

func (e *BillingEngine) AddSurcharges(customer domain.Customer, charges []ChargeItem, surcharge map[string]float64) []ChargeItem {
	keys := make([]string, 0, len(surcharge))
	for key := range surcharge {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		amount := surcharge[key]
		if amount <= 0 {
			continue
		}
		charges = append(charges, ChargeItem{
			Code:        key,
			Description: strings.ReplaceAll(key, "_", " "),
			Amount:      roundMoney(amount),
			CustomerID:  customer.ID,
		})
	}
	return charges
}

func (e *BillingEngine) SumCharges(charges []ChargeItem) float64 {
	total := 0.0
	for _, charge := range charges {
		total += charge.Amount
	}
	return roundMoney(total)
}

func (e *BillingEngine) ToInvoice(customer domain.Customer, invoice domain.Invoice, charges []ChargeItem, periodStart, periodEnd time.Time, mode string) domain.Invoice {
	invoice.InvoiceNumber = e.BuildInvoiceNumber(customer.ID, periodStart, periodEnd)
	invoice.CustomerID = customer.ID
	invoice.PeriodStart = periodStart
	invoice.PeriodEnd = periodEnd
	invoice.BillingMode = NormalizeBillingMode(mode)
	invoice.Subtotal = e.SumCharges(charges)
	invoice.TotalAmount = roundMoney(invoice.Subtotal + invoice.Surcharge)
	invoice.DueDate = periodEnd.AddDate(0, 0, 30)
	if invoice.Status == "" {
		invoice.Status = domain.InvoicePending
	}
	invoice.Items = e.ToItems(charges, invoice.ID)
	return invoice
}

func (e *BillingEngine) ToItems(charges []ChargeItem, invoiceID int64) []domain.InvoiceItem {
	items := make([]domain.InvoiceItem, 0, len(charges))
	for _, charge := range charges {
		items = append(items, domain.InvoiceItem{
			InvoiceID:   invoiceID,
			TaskID:      charge.TaskID,
			ItemType:    charge.Code,
			Description: charge.Description,
			Quantity:    1,
			UnitPrice:   charge.Amount,
			Amount:      charge.Amount,
		})
	}
	return items
}

func (e *BillingEngine) Coverage(taskCount int, days int) float64 {
	if days <= 0 {
		return 0
	}
	return roundMoney(float64(taskCount) / float64(days))
}

func (e *BillingEngine) Overdue(invoice domain.Invoice, now time.Time) bool {
	if invoice.DueDate.IsZero() {
		return false
	}
	return now.After(invoice.DueDate)
}

func (e *BillingEngine) ReminderMessage(invoice domain.Invoice) string {
	return fmt.Sprintf("invoice %s overdue by %s", invoice.InvoiceNumber, time.Since(invoice.DueDate).Round(time.Hour))
}

func (e *BillingEngine) Balance(invoice domain.Invoice) float64 {
	return roundMoney(invoice.TotalAmount - invoice.PaidAmount)
}

func (e *BillingEngine) ApplyPayment(invoice *domain.Invoice, amount float64) {
	invoice.PaidAmount = roundMoney(invoice.PaidAmount + amount)
	switch {
	case invoice.PaidAmount <= 0:
		invoice.Status = domain.InvoicePending
	case invoice.PaidAmount < invoice.TotalAmount:
		invoice.Status = domain.InvoicePartial
	default:
		invoice.Status = domain.InvoicePaid
	}
}

func (e *BillingEngine) WriteOff(invoice *domain.Invoice) {
	invoice.Status = domain.InvoiceWrittenOff
}

func (e *BillingEngine) MergeCharges(base []ChargeItem, extra []ChargeItem) []ChargeItem {
	merged := append([]ChargeItem(nil), base...)
	merged = append(merged, extra...)
	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].CustomerID == merged[j].CustomerID {
			if merged[i].TaskID == merged[j].TaskID {
				return merged[i].Code < merged[j].Code
			}
			return merged[i].TaskID < merged[j].TaskID
		}
		return merged[i].CustomerID < merged[j].CustomerID
	})
	return merged
}

func (e *BillingEngine) SplitByTask(charges []ChargeItem) map[int64][]ChargeItem {
	result := make(map[int64][]ChargeItem)
	for _, charge := range charges {
		result[charge.TaskID] = append(result[charge.TaskID], charge)
	}
	return result
}

func (e *BillingEngine) SplitByCode(charges []ChargeItem) map[string][]ChargeItem {
	result := make(map[string][]ChargeItem)
	for _, charge := range charges {
		result[charge.Code] = append(result[charge.Code], charge)
	}
	return result
}

func (e *BillingEngine) ComputeSurcharge(taskCount int, routeDistance float64, holidays int, extras map[string]bool) map[string]float64 {
	surcharge := make(map[string]float64)
	if extras["frequency"] && taskCount > 1 {
		surcharge["extra_frequency"] = float64(taskCount-1) * 20
	}
	if extras["distance"] && routeDistance > 0 {
		surcharge["distance_fee"] = roundMoney(routeDistance * 1.5)
	}
	if extras["holiday"] && holidays > 0 {
		surcharge["holiday_service"] = float64(holidays) * 50
	}
	if extras["special"] {
		surcharge["special_waste"] = 100
	}
	return surcharge
}

func (e *BillingEngine) GenerateItemsFromTasks(tasks []domain.Task, customer domain.Customer, categories map[int64]domain.WasteCategory, rates map[int64]float64, weights map[int64]float64, rateCard RateCard, mode string) []domain.InvoiceItem {
	charges := e.BuildCharges(customer, tasks, categories, weights, rates, rateCard, mode)
	items := e.ToItems(charges, 0)
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].TaskID == items[j].TaskID {
			return items[i].ItemType < items[j].ItemType
		}
		return items[i].TaskID < items[j].TaskID
	})
	return items
}

func (e *BillingEngine) RateForCategory(categoryID int64, rates map[int64]float64, fallback float64) float64 {
	if rate, ok := rates[categoryID]; ok && rate > 0 {
		return rate
	}
	return fallback
}

func (e *BillingEngine) ApplyFloor(total float64, floor float64) float64 {
	if floor <= 0 {
		return roundMoney(total)
	}
	if total < floor {
		return roundMoney(floor)
	}
	return roundMoney(total)
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func NormalizeBillingMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "per_trip", "trip", "trip_fee":
		return "per_trip"
	case "weight", "per_weight", "ton", "tons":
		return "weight"
	case "mixed", "base_plus_weight":
		return "mixed"
	default:
		return "weight"
	}
}

