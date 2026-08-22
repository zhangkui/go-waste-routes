package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type VehicleRepo struct{ db *DB }

func NewVehicleRepo(db *DB) *VehicleRepo { return &VehicleRepo{db: db} }

func (r *VehicleRepo) Create(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO vehicles(plate_number,model,rated_load_tons,current_mileage_km,status,inspection_due_date,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, vehicle.PlateNumber, vehicle.Model, vehicle.RatedLoadTons, vehicle.CurrentMileageKm, vehicle.Status, vehicle.InspectionDueDate, time.Now(), time.Now())
	if err != nil {
		return domain.Vehicle{}, err
	}
	id, _ := result.LastInsertId()
	vehicle.ID = id
	return vehicle, nil
}

func (r *VehicleRepo) Get(ctx context.Context, id int64) (domain.Vehicle, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,plate_number,model,rated_load_tons,current_mileage_km,status,inspection_due_date,created_at,updated_at,deleted_at FROM vehicles WHERE id=? AND deleted_at IS NULL`, id)
	var vehicle domain.Vehicle
	var dueDate sql.NullTime
	var deletedAt sql.NullTime
	if err := row.Scan(&vehicle.ID, &vehicle.PlateNumber, &vehicle.Model, &vehicle.RatedLoadTons, &vehicle.CurrentMileageKm, &vehicle.Status, &dueDate, &vehicle.CreatedAt, &vehicle.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Vehicle{}, false, nil
		}
		return domain.Vehicle{}, false, err
	}
	vehicle.DeletedAt = scanNullTime(deletedAt)
	vehicle.InspectionDueDate = dueDate.Time
	return vehicle, true, nil
}

func (r *VehicleRepo) List(ctx context.Context, page, pageSize int) ([]domain.Vehicle, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM vehicles WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,plate_number,model,rated_load_tons,current_mileage_km,status,inspection_due_date,created_at,updated_at,deleted_at FROM vehicles WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Vehicle, 0)
	for rows.Next() {
		var item domain.Vehicle
		var dueDate sql.NullTime
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.PlateNumber, &item.Model, &item.RatedLoadTons, &item.CurrentMileageKm, &item.Status, &dueDate, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		if dueDate.Valid {
			item.InspectionDueDate = dueDate.Time
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *VehicleRepo) Update(ctx context.Context, id int64, vehicle domain.Vehicle) (domain.Vehicle, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE vehicles SET plate_number=?,model=?,rated_load_tons=?,current_mileage_km=?,status=?,inspection_due_date=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, vehicle.PlateNumber, vehicle.Model, vehicle.RatedLoadTons, vehicle.CurrentMileageKm, vehicle.Status, vehicle.InspectionDueDate, time.Now(), id)
	if err != nil {
		return domain.Vehicle{}, err
	}
	vehicle.ID = id
	return vehicle, nil
}

func (r *VehicleRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE vehicles SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

