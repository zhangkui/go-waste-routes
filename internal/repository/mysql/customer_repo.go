package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type CustomerRepo struct{ db *DB }

func NewCustomerRepo(db *DB) *CustomerRepo { return &CustomerRepo{db: db} }

func (r *CustomerRepo) Create(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO customers(name,contact_name,contact_phone,address,region,level,contract_number,service_weekdays,service_month_days,temporary_allowed,service_windows,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, customer.Name, customer.ContactName, customer.ContactPhone, customer.Address, customer.Region, customer.Level, customer.ContractNumber, customer.ServiceWeekdays, customer.ServiceMonthDays, customer.TemporaryAllowed, customer.ServiceWindows, customer.Status, time.Now(), time.Now())
	if err != nil {
		return domain.Customer{}, err
	}
	id, _ := result.LastInsertId()
	customer.ID = id
	return customer, nil
}

func (r *CustomerRepo) Get(ctx context.Context, id int64) (domain.Customer, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,name,contact_name,contact_phone,address,region,level,contract_number,service_weekdays,service_month_days,temporary_allowed,service_windows,status,created_at,updated_at,deleted_at FROM customers WHERE id=? AND deleted_at IS NULL`, id)
	var customer domain.Customer
	var deletedAt sql.NullTime
	if err := row.Scan(&customer.ID, &customer.Name, &customer.ContactName, &customer.ContactPhone, &customer.Address, &customer.Region, &customer.Level, &customer.ContractNumber, &customer.ServiceWeekdays, &customer.ServiceMonthDays, &customer.TemporaryAllowed, &customer.ServiceWindows, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Customer{}, false, nil
		}
		return domain.Customer{}, false, err
	}
	customer.DeletedAt = scanNullTime(deletedAt)
	return customer, true, nil
}

func (r *CustomerRepo) List(ctx context.Context, page, pageSize int) ([]domain.Customer, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,name,contact_name,contact_phone,address,region,level,contract_number,service_weekdays,service_month_days,temporary_allowed,service_windows,status,created_at,updated_at,deleted_at FROM customers WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Customer, 0)
	for rows.Next() {
		var item domain.Customer
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.ContactName, &item.ContactPhone, &item.Address, &item.Region, &item.Level, &item.ContractNumber, &item.ServiceWeekdays, &item.ServiceMonthDays, &item.TemporaryAllowed, &item.ServiceWindows, &item.Status, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *CustomerRepo) Update(ctx context.Context, id int64, customer domain.Customer) (domain.Customer, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE customers SET name=?,contact_name=?,contact_phone=?,address=?,region=?,level=?,contract_number=?,service_weekdays=?,service_month_days=?,temporary_allowed=?,service_windows=?,status=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, customer.Name, customer.ContactName, customer.ContactPhone, customer.Address, customer.Region, customer.Level, customer.ContractNumber, customer.ServiceWeekdays, customer.ServiceMonthDays, customer.TemporaryAllowed, customer.ServiceWindows, customer.Status, time.Now(), id)
	if err != nil {
		return domain.Customer{}, err
	}
	customer.ID = id
	return customer, nil
}

func (r *CustomerRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE customers SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

func (r *CustomerRepo) SetStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE customers SET status=?,updated_at=? WHERE id=?`, status, time.Now(), id)
	return err
}
