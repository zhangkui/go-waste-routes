package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type InvoiceRepo struct{ db *DB }

func NewInvoiceRepo(db *DB) *InvoiceRepo { return &InvoiceRepo{db: db} }

func (r *InvoiceRepo) Create(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO invoices(invoice_number,customer_id,period_start,period_end,billing_mode,subtotal,surcharge,total_amount,paid_amount,status,due_date,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`, invoice.InvoiceNumber, invoice.CustomerID, invoice.PeriodStart, invoice.PeriodEnd, invoice.BillingMode, invoice.Subtotal, invoice.Surcharge, invoice.TotalAmount, invoice.PaidAmount, invoice.Status, invoice.DueDate, time.Now(), time.Now())
	if err != nil {
		return domain.Invoice{}, err
	}
	id, _ := result.LastInsertId()
	invoice.ID = id
	return invoice, nil
}

func (r *InvoiceRepo) Get(ctx context.Context, id int64) (domain.Invoice, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,invoice_number,customer_id,period_start,period_end,billing_mode,subtotal,surcharge,total_amount,paid_amount,status,due_date,created_at,updated_at,deleted_at FROM invoices WHERE id=? AND deleted_at IS NULL`, id)
	var invoice domain.Invoice
	var dueDate sql.NullTime
	var deletedAt sql.NullTime
	if err := row.Scan(&invoice.ID, &invoice.InvoiceNumber, &invoice.CustomerID, &invoice.PeriodStart, &invoice.PeriodEnd, &invoice.BillingMode, &invoice.Subtotal, &invoice.Surcharge, &invoice.TotalAmount, &invoice.PaidAmount, &invoice.Status, &dueDate, &invoice.CreatedAt, &invoice.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Invoice{}, false, nil
		}
		return domain.Invoice{}, false, err
	}
	if dueDate.Valid {
		invoice.DueDate = dueDate.Time
	}
	invoice.DeletedAt = scanNullTime(deletedAt)
	return invoice, true, nil
}

func (r *InvoiceRepo) List(ctx context.Context, page, pageSize int) ([]domain.Invoice, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,invoice_number,customer_id,period_start,period_end,billing_mode,subtotal,surcharge,total_amount,paid_amount,status,due_date,created_at,updated_at,deleted_at FROM invoices WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Invoice, 0)
	for rows.Next() {
		var item domain.Invoice
		var dueDate sql.NullTime
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.InvoiceNumber, &item.CustomerID, &item.PeriodStart, &item.PeriodEnd, &item.BillingMode, &item.Subtotal, &item.Surcharge, &item.TotalAmount, &item.PaidAmount, &item.Status, &dueDate, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		if dueDate.Valid {
			item.DueDate = dueDate.Time
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *InvoiceRepo) Update(ctx context.Context, id int64, invoice domain.Invoice) (domain.Invoice, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE invoices SET invoice_number=?,customer_id=?,period_start=?,period_end=?,billing_mode=?,subtotal=?,surcharge=?,total_amount=?,paid_amount=?,status=?,due_date=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, invoice.InvoiceNumber, invoice.CustomerID, invoice.PeriodStart, invoice.PeriodEnd, invoice.BillingMode, invoice.Subtotal, invoice.Surcharge, invoice.TotalAmount, invoice.PaidAmount, invoice.Status, invoice.DueDate, time.Now(), id)
	if err != nil {
		return domain.Invoice{}, err
	}
	invoice.ID = id
	return invoice, nil
}

func (r *InvoiceRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE invoices SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

