package service

import (
	"fmt"
	"sort"
	"time"

	"go-waste-routes/internal/domain"
)

type InvoiceService struct{}

func NewInvoiceService() *InvoiceService { return &InvoiceService{} }

func (s *InvoiceService) GenerateNumber(customerID int64, periodStart, periodEnd time.Time) string {
	return fmt.Sprintf("INV-%d-%s-%s", customerID, periodStart.Format("20060102"), periodEnd.Format("20060102"))
}

func (s *InvoiceService) ComposeItems(task domain.Task, category domain.WasteCategory, netWeight float64, unitPrice float64) []domain.InvoiceItem {
	amount := netWeight * unitPrice
	return []domain.InvoiceItem{
		{
			TaskID:          task.ID,
			ServiceDate:     task.PlanDate,
			WasteCategoryID: category.ID,
			NetWeight:       netWeight,
			ItemType:        "service",
			Description:     category.Name,
			UnitPrice:       unitPrice,
			Quantity:        netWeight,
			Amount:          amount,
		},
	}
}

func (s *InvoiceService) ApplySurcharges(base float64, items []float64) float64 {
	total := base
	for _, item := range items {
		total += item
	}
	return total
}

func (s *InvoiceService) SortItems(items []domain.InvoiceItem) []domain.InvoiceItem {
	ordered := append([]domain.InvoiceItem(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].ServiceDate.Before(ordered[j].ServiceDate) })
	return ordered
}

func (s *InvoiceService) MarkPaid(invoice *domain.Invoice, amount float64) {
	invoice.PaidAmount += amount
	switch {
	case invoice.PaidAmount <= 0:
		invoice.Status = domain.InvoicePending
	case invoice.PaidAmount < invoice.TotalAmount:
		invoice.Status = domain.InvoicePartial
	default:
		invoice.Status = domain.InvoicePaid
	}
}

func (s *InvoiceService) MarkWrittenOff(invoice *domain.Invoice) {
	invoice.Status = domain.InvoiceWrittenOff
}

func (s *InvoiceService) Overdue(dueDate time.Time, now time.Time) bool {
	return !dueDate.IsZero() && now.After(dueDate.AddDate(0, 0, 30))
}

func (s *InvoiceService) BillSummary(invoice domain.Invoice) string {
	return fmt.Sprintf("%s total=%.2f paid=%.2f", invoice.InvoiceNumber, invoice.TotalAmount, invoice.PaidAmount)
}

