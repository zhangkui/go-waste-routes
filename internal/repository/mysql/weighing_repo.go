package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type WeighingRepo struct{ db *DB }

func NewWeighingRepo(db *DB) *WeighingRepo { return &WeighingRepo{db: db} }

func (r *WeighingRepo) Create(ctx context.Context, record domain.WeighingRecord) (domain.WeighingRecord, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO weighing_records(task_id,customer_id,vehicle_id,scale_id,gross_weight,tare_weight,net_weight,weigh_time,is_manual,manual_reason,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`, record.TaskID, record.CustomerID, record.VehicleID, record.ScaleID, record.GrossWeight, record.TareWeight, record.NetWeight, record.WeighTime, record.Manual, record.ManualReason, record.Status, time.Now(), time.Now())
	if err != nil {
		return domain.WeighingRecord{}, err
	}
	id, _ := result.LastInsertId()
	record.ID = id
	return record, nil
}

func (r *WeighingRepo) Get(ctx context.Context, id int64) (domain.WeighingRecord, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,task_id,customer_id,vehicle_id,scale_id,gross_weight,tare_weight,net_weight,weigh_time,is_manual,manual_reason,status,created_at,updated_at,deleted_at FROM weighing_records WHERE id=? AND deleted_at IS NULL`, id)
	var record domain.WeighingRecord
	var deletedAt sql.NullTime
	if err := row.Scan(&record.ID, &record.TaskID, &record.CustomerID, &record.VehicleID, &record.ScaleID, &record.GrossWeight, &record.TareWeight, &record.NetWeight, &record.WeighTime, &record.Manual, &record.ManualReason, &record.Status, &record.CreatedAt, &record.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.WeighingRecord{}, false, nil
		}
		return domain.WeighingRecord{}, false, err
	}
	record.DeletedAt = scanNullTime(deletedAt)
	return record, true, nil
}

func (r *WeighingRepo) List(ctx context.Context, page, pageSize int) ([]domain.WeighingRecord, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM weighing_records WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,task_id,customer_id,vehicle_id,scale_id,gross_weight,tare_weight,net_weight,weigh_time,is_manual,manual_reason,status,created_at,updated_at,deleted_at FROM weighing_records WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.WeighingRecord, 0)
	for rows.Next() {
		var item domain.WeighingRecord
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.TaskID, &item.CustomerID, &item.VehicleID, &item.ScaleID, &item.GrossWeight, &item.TareWeight, &item.NetWeight, &item.WeighTime, &item.Manual, &item.ManualReason, &item.Status, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *WeighingRepo) Update(ctx context.Context, id int64, record domain.WeighingRecord) (domain.WeighingRecord, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE weighing_records SET task_id=?,customer_id=?,vehicle_id=?,scale_id=?,gross_weight=?,tare_weight=?,net_weight=?,weigh_time=?,is_manual=?,manual_reason=?,status=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, record.TaskID, record.CustomerID, record.VehicleID, record.ScaleID, record.GrossWeight, record.TareWeight, record.NetWeight, record.WeighTime, record.Manual, record.ManualReason, record.Status, time.Now(), id)
	if err != nil {
		return domain.WeighingRecord{}, err
	}
	record.ID = id
	return record, nil
}

func (r *WeighingRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE weighing_records SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}
