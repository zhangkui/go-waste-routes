package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type PlanRepo struct{ db *DB }

func NewPlanRepo(db *DB) *PlanRepo { return &PlanRepo{db: db} }

func (r *PlanRepo) Create(ctx context.Context, plan domain.CollectionPlan) (domain.CollectionPlan, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO collection_plans(name,customer_id,version,status,holiday_strategy,effective_from,effective_to,estimated_weight_tons,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, plan.Name, plan.CustomerID, plan.Version, plan.Status, plan.HolidayStrategy, plan.EffectiveFrom, plan.EffectiveTo, plan.EstimatedWeightTons, time.Now(), time.Now())
	if err != nil {
		return domain.CollectionPlan{}, err
	}
	id, _ := result.LastInsertId()
	plan.ID = id
	return plan, nil
}

func (r *PlanRepo) Get(ctx context.Context, id int64) (domain.CollectionPlan, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,name,customer_id,version,status,holiday_strategy,effective_from,effective_to,estimated_weight_tons,created_at,updated_at,deleted_at FROM collection_plans WHERE id=? AND deleted_at IS NULL`, id)
	var plan domain.CollectionPlan
	var effectiveTo sql.NullTime
	var deletedAt sql.NullTime
	if err := row.Scan(&plan.ID, &plan.Name, &plan.CustomerID, &plan.Version, &plan.Status, &plan.HolidayStrategy, &plan.EffectiveFrom, &effectiveTo, &plan.EstimatedWeightTons, &plan.CreatedAt, &plan.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.CollectionPlan{}, false, nil
		}
		return domain.CollectionPlan{}, false, err
	}
	plan.EffectiveTo = scanNullTime(effectiveTo)
	plan.DeletedAt = scanNullTime(deletedAt)
	return plan, true, nil
}

func (r *PlanRepo) List(ctx context.Context, page, pageSize int) ([]domain.CollectionPlan, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM collection_plans WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,name,customer_id,version,status,holiday_strategy,effective_from,effective_to,estimated_weight_tons,created_at,updated_at,deleted_at FROM collection_plans WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.CollectionPlan, 0)
	for rows.Next() {
		var item domain.CollectionPlan
		var effectiveTo sql.NullTime
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Name, &item.CustomerID, &item.Version, &item.Status, &item.HolidayStrategy, &item.EffectiveFrom, &effectiveTo, &item.EstimatedWeightTons, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.EffectiveTo = scanNullTime(effectiveTo)
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PlanRepo) Update(ctx context.Context, id int64, plan domain.CollectionPlan) (domain.CollectionPlan, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE collection_plans SET name=?,customer_id=?,version=?,status=?,holiday_strategy=?,effective_from=?,effective_to=?,estimated_weight_tons=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, plan.Name, plan.CustomerID, plan.Version, plan.Status, plan.HolidayStrategy, plan.EffectiveFrom, plan.EffectiveTo, plan.EstimatedWeightTons, time.Now(), id)
	if err != nil {
		return domain.CollectionPlan{}, err
	}
	plan.ID = id
	return plan, nil
}

func (r *PlanRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE collection_plans SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

func (r *PlanRepo) ActivateVersion(ctx context.Context, customerID int64, version int) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE collection_plans SET status='active',updated_at=? WHERE customer_id=? AND version=?`, time.Now(), customerID, version)
	return err
}
