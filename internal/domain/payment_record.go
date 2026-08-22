package domain

import "time"

type PaymentRecord struct {
	Base
	InvoiceID     int64     `json:"invoice_id"`
	PaymentTime   time.Time `json:"payment_time"`
	Amount        float64   `json:"amount"`
	OperatorID    int64     `json:"operator_id"`
	PaymentMethod string    `json:"payment_method"`
	VoucherNumber string    `json:"voucher_number"`
	Notes         string    `json:"notes"`
}
