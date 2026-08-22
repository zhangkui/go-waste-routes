package domain

import "time"

type Invoice struct {
	Base
	InvoiceNumber string        `json:"invoice_number"`
	CustomerID    int64         `json:"customer_id"`
	PeriodStart   time.Time     `json:"period_start"`
	PeriodEnd     time.Time     `json:"period_end"`
	BillingMode   string        `json:"billing_mode"`
	Subtotal      float64       `json:"subtotal"`
	Surcharge     float64       `json:"surcharge"`
	TotalAmount   float64       `json:"total_amount"`
	PaidAmount    float64       `json:"paid_amount"`
	Status        string        `json:"status"`
	DueDate       time.Time     `json:"due_date"`
	Items         []InvoiceItem `json:"items,omitempty"`
}

const (
	InvoicePending    = "pending"
	InvoiceConfirmed  = "confirmed"
	InvoicePartial    = "partial"
	InvoicePaid       = "paid"
	InvoiceWrittenOff = "written_off"
)
