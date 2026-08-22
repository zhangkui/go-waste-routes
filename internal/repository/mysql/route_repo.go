package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type RouteRepo struct{ db *DB }

func NewRouteRepo(db *DB) *RouteRepo { return &RouteRepo{db: db} }

func (r *RouteRepo) Create(ctx context.Context, route domain.Route) (domain.Route, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO routes(name,plan_id,vehicle_id,driver_id,status,estimated_weight_tons,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, route.Name, route.PlanID, route.VehicleID, route.DriverID, route.Status, route.EstimatedWeightTons, time.Now(), time.Now())
	if err != nil {
		return domain.Route{}, err
	}
	id, _ := result.LastInsertId()
	route.ID = id
	return route, nil
}

func (r *RouteRepo) Get(ctx context.Context, id int64) (domain.Route, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,name,plan_id,vehicle_id,driver_id,status,estimated_weight_tons,created_at,updated_at,deleted_at FROM routes WHERE id=? AND deleted_at IS NULL`, id)
	var route domain.Route
	var deletedAt sql.NullTime
	if err := row.Scan(&route.ID, &route.Name, &route.PlanID, &route.VehicleID, &route.DriverID, &route.Status, &route.EstimatedWeightTons, &route.CreatedAt, &route.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Route{}, false, nil
		}
		return domain.Route{}, false, err
	}
	route.DeletedAt = scanNullTime(deletedAt)
	return route, true, nil
}

func (r *RouteRepo) List(ctx context.Context, page, pageSize int) ([]domain.Route, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM routes WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,name,plan_id,vehicle_id,driver_id,status,estimated_weight_tons,created_at,updated_at,deleted_at FROM routes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Route, 0)
	for rows.Next() {
		var item domain.Route
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.PlanID, &item.VehicleID, &item.DriverID, &item.Status, &item.EstimatedWeightTons, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *RouteRepo) Update(ctx context.Context, id int64, route domain.Route) (domain.Route, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE routes SET name=?,plan_id=?,vehicle_id=?,driver_id=?,status=?,estimated_weight_tons=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, route.Name, route.PlanID, route.VehicleID, route.DriverID, route.Status, route.EstimatedWeightTons, time.Now(), id)
	if err != nil {
		return domain.Route{}, err
	}
	route.ID = id
	return route, nil
}

func (r *RouteRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE routes SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

func (r *RouteRepo) SaveStops(ctx context.Context, routeID int64, stops []domain.RouteStop) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM route_stops WHERE route_id=?`, routeID); err != nil {
			return err
		}
		for _, stop := range stops {
			if _, err := tx.ExecContext(ctx, `INSERT INTO route_stops(route_id,customer_id,sequence_no,estimated_arrival,stay_minutes,estimated_weight_tons,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, routeID, stop.CustomerID, stop.Sequence, stop.EstimatedArrival, stop.StayMinutes, stop.EstimatedWeightTons, time.Now(), time.Now()); err != nil {
				return err
			}
		}
		return nil
	})
}
