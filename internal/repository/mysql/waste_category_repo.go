package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type WasteCategoryRepo struct{ db *DB }

func NewWasteCategoryRepo(db *DB) *WasteCategoryRepo { return &WasteCategoryRepo{db: db} }

func (r *WasteCategoryRepo) Create(ctx context.Context, category domain.WasteCategory) (domain.WasteCategory, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO waste_categories(code,name,description,default_rate,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, category.Code, category.Name, category.Description, category.DefaultRate, category.Status, time.Now(), time.Now())
	if err != nil {
		return domain.WasteCategory{}, err
	}
	id, _ := result.LastInsertId()
	category.ID = id
	return category, nil
}

func (r *WasteCategoryRepo) Get(ctx context.Context, id int64) (domain.WasteCategory, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,code,name,description,default_rate,status,created_at,updated_at,deleted_at FROM waste_categories WHERE id=? AND deleted_at IS NULL`, id)
	var category domain.WasteCategory
	var deletedAt sql.NullTime
	if err := row.Scan(&category.ID, &category.Code, &category.Name, &category.Description, &category.DefaultRate, &category.Status, &category.CreatedAt, &category.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.WasteCategory{}, false, nil
		}
		return domain.WasteCategory{}, false, err
	}
	category.DeletedAt = scanNullTime(deletedAt)
	return category, true, nil
}

func (r *WasteCategoryRepo) List(ctx context.Context, page, pageSize int) ([]domain.WasteCategory, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM waste_categories WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,code,name,description,default_rate,status,created_at,updated_at,deleted_at FROM waste_categories WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.WasteCategory, 0)
	for rows.Next() {
		var item domain.WasteCategory
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.DefaultRate, &item.Status, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *WasteCategoryRepo) Update(ctx context.Context, id int64, category domain.WasteCategory) (domain.WasteCategory, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE waste_categories SET code=?,name=?,description=?,default_rate=?,status=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, category.Code, category.Name, category.Description, category.DefaultRate, category.Status, time.Now(), id)
	if err != nil {
		return domain.WasteCategory{}, err
	}
	category.ID = id
	return category, nil
}

func (r *WasteCategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE waste_categories SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

