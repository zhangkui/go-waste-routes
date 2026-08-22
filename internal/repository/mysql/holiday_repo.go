package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type HolidayRepo struct{ db *DB }

func NewHolidayRepo(db *DB) *HolidayRepo { return &HolidayRepo{db: db} }

func (r *HolidayRepo) Create(ctx context.Context, holiday domain.HolidayConfig) (domain.HolidayConfig, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO holiday_config(holiday_date,name,is_workday,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?)`, holiday.HolidayDate, holiday.Name, holiday.IsWorkday, holiday.Enabled, time.Now(), time.Now())
	if err != nil {
		return domain.HolidayConfig{}, err
	}
	id, _ := result.LastInsertId()
	holiday.ID = id
	return holiday, nil
}

func (r *HolidayRepo) Get(ctx context.Context, id int64) (domain.HolidayConfig, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,holiday_date,name,is_workday,enabled,created_at,updated_at,deleted_at FROM holiday_config WHERE id=? AND deleted_at IS NULL`, id)
	var holiday domain.HolidayConfig
	var deletedAt sql.NullTime
	if err := row.Scan(&holiday.ID, &holiday.HolidayDate, &holiday.Name, &holiday.IsWorkday, &holiday.Enabled, &holiday.CreatedAt, &holiday.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.HolidayConfig{}, false, nil
		}
		return domain.HolidayConfig{}, false, err
	}
	holiday.DeletedAt = scanNullTime(deletedAt)
	return holiday, true, nil
}

func (r *HolidayRepo) List(ctx context.Context, page, pageSize int) ([]domain.HolidayConfig, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM holiday_config WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,holiday_date,name,is_workday,enabled,created_at,updated_at,deleted_at FROM holiday_config WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.HolidayConfig, 0)
	for rows.Next() {
		var item domain.HolidayConfig
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.HolidayDate, &item.Name, &item.IsWorkday, &item.Enabled, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *HolidayRepo) Update(ctx context.Context, id int64, holiday domain.HolidayConfig) (domain.HolidayConfig, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE holiday_config SET holiday_date=?,name=?,is_workday=?,enabled=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, holiday.HolidayDate, holiday.Name, holiday.IsWorkday, holiday.Enabled, time.Now(), id)
	if err != nil {
		return domain.HolidayConfig{}, err
	}
	holiday.ID = id
	return holiday, nil
}

func (r *HolidayRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE holiday_config SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

