package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type TaskRepo struct{ db *DB }

func NewTaskRepo(db *DB) *TaskRepo { return &TaskRepo{db: db} }

func (r *TaskRepo) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO tasks(task_number,plan_id,route_id,driver_id,plan_date,status,claimed_at,completed_at,mileage_km,fuel_liters,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, task.TaskNumber, task.PlanID, task.RouteID, task.DriverID, task.PlanDate, task.Status, task.ClaimedAt, task.CompletedAt, task.MileageKm, task.FuelLiters, time.Now(), time.Now())
	if err != nil {
		return domain.Task{}, err
	}
	id, _ := result.LastInsertId()
	task.ID = id
	return task, nil
}

func (r *TaskRepo) Get(ctx context.Context, id int64) (domain.Task, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,task_number,plan_id,route_id,driver_id,plan_date,status,claimed_at,completed_at,mileage_km,fuel_liters,created_at,updated_at,deleted_at FROM tasks WHERE id=? AND deleted_at IS NULL`, id)
	var task domain.Task
	var routeID sql.NullInt64
	var driverID sql.NullInt64
	var claimedAt sql.NullTime
	var completedAt sql.NullTime
	var deletedAt sql.NullTime
	if err := row.Scan(&task.ID, &task.TaskNumber, &task.PlanID, &routeID, &driverID, &task.PlanDate, &task.Status, &claimedAt, &completedAt, &task.MileageKm, &task.FuelLiters, &task.CreatedAt, &task.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, false, nil
		}
		return domain.Task{}, false, err
	}
	task.RouteID = scanNullInt64(routeID)
	task.DriverID = scanNullInt64(driverID)
	task.ClaimedAt = scanNullTime(claimedAt)
	task.CompletedAt = scanNullTime(completedAt)
	task.DeletedAt = scanNullTime(deletedAt)
	return task, true, nil
}

func (r *TaskRepo) List(ctx context.Context, page, pageSize int) ([]domain.Task, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,task_number,plan_id,route_id,driver_id,plan_date,status,claimed_at,completed_at,mileage_km,fuel_liters,created_at,updated_at,deleted_at FROM tasks WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Task, 0)
	for rows.Next() {
		var item domain.Task
		var routeID sql.NullInt64
		var driverID sql.NullInt64
		var claimedAt sql.NullTime
		var completedAt sql.NullTime
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.TaskNumber, &item.PlanID, &routeID, &driverID, &item.PlanDate, &item.Status, &claimedAt, &completedAt, &item.MileageKm, &item.FuelLiters, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.RouteID = scanNullInt64(routeID)
		item.DriverID = scanNullInt64(driverID)
		item.ClaimedAt = scanNullTime(claimedAt)
		item.CompletedAt = scanNullTime(completedAt)
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *TaskRepo) Update(ctx context.Context, id int64, task domain.Task) (domain.Task, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE tasks SET task_number=?,plan_id=?,route_id=?,driver_id=?,plan_date=?,status=?,claimed_at=?,completed_at=?,mileage_km=?,fuel_liters=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, task.TaskNumber, task.PlanID, task.RouteID, task.DriverID, task.PlanDate, task.Status, task.ClaimedAt, task.CompletedAt, task.MileageKm, task.FuelLiters, time.Now(), id)
	if err != nil {
		return domain.Task{}, err
	}
	task.ID = id
	return task, nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE tasks SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

