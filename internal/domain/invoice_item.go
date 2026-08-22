package domain

import "time"

type InvoiceItem struct {
	Base
	InvoiceID       int64     `json:"invoice_id"`
	TaskID          int64     `json:"task_id"`
	ServiceDate     time.Time `json:"service_date"`
	WasteCategoryID int64     `json:"waste_category_id"`
	NetWeight       float64   `json:"net_weight"`
	ItemType        string    `json:"item_type"`
	Description     string    `json:"description"`
	UnitPrice       float64   `json:"unit_price"`
	Quantity        float64   `json:"quantity"`
	Amount          float64   `json:"amount"`
}
