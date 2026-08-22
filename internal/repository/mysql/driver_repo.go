package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type DriverRepo struct{ db *DB }

func NewDriverRepo(db *DB) *DriverRepo { return &DriverRepo{db: db} }

func (r *DriverRepo) Create(ctx context.Context, driver domain.Driver) (domain.Driver, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO drivers(name,phone,license_number,license_class,hire_date,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, driver.Name, driver.Phone, driver.LicenseNumber, driver.LicenseClass, driver.HireDate, driver.Status, time.Now(), time.Now())
	if err != nil {
		return domain.Driver{}, err
	}
	id, _ := result.LastInsertId()
	driver.ID = id
	return driver, nil
}

func (r *DriverRepo) Get(ctx context.Context, id int64) (domain.Driver, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,name,phone,license_number,license_class,hire_date,status,created_at,updated_at,deleted_at FROM drivers WHERE id=? AND deleted_at IS NULL`, id)
	var driver domain.Driver
	var hireDate sql.NullTime
	var deletedAt sql.NullTime
	if err := row.Scan(&driver.ID, &driver.Name, &driver.Phone, &driver.LicenseNumber, &driver.LicenseClass, &hireDate, &driver.Status, &driver.CreatedAt, &driver.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Driver{}, false, nil
		}
		return domain.Driver{}, false, err
	}
	driver.DeletedAt = scanNullTime(deletedAt)
	if hireDate.Valid {
		driver.HireDate = hireDate.Time
	}
	return driver, true, nil
}

func (r *DriverRepo) List(ctx context.Context, page, pageSize int) ([]domain.Driver, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM drivers WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,name,phone,license_number,license_class,hire_date,status,created_at,updated_at,deleted_at FROM drivers WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Driver, 0)
	for rows.Next() {
		var item domain.Driver
		var hireDate sql.NullTime
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.LicenseNumber, &item.LicenseClass, &hireDate, &item.Status, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		if hireDate.Valid {
			item.HireDate = hireDate.Time
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *DriverRepo) Update(ctx context.Context, id int64, driver domain.Driver) (domain.Driver, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE drivers SET name=?,phone=?,license_number=?,license_class=?,hire_date=?,status=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, driver.Name, driver.Phone, driver.LicenseNumber, driver.LicenseClass, driver.HireDate, driver.Status, time.Now(), id)
	if err != nil {
		return domain.Driver{}, err
	}
	driver.ID = id
	return driver, nil
}

func (r *DriverRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE drivers SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

