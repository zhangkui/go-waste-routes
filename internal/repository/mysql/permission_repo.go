package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type PermissionRepo struct{ db *DB }

func NewPermissionRepo(db *DB) *PermissionRepo { return &PermissionRepo{db: db} }

func (r *PermissionRepo) Create(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO permissions(code,name,module,description,created_at,updated_at) VALUES(?,?,?,?,?,?)`, permission.Code, permission.Name, permission.Module, permission.Description, time.Now(), time.Now())
	if err != nil {
		return domain.Permission{}, err
	}
	id, _ := result.LastInsertId()
	permission.ID = id
	return permission, nil
}

func (r *PermissionRepo) Get(ctx context.Context, id int64) (domain.Permission, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,code,name,module,description,created_at,updated_at,deleted_at FROM permissions WHERE id=? AND deleted_at IS NULL`, id)
	var permission domain.Permission
	var deletedAt sql.NullTime
	if err := row.Scan(&permission.ID, &permission.Code, &permission.Name, &permission.Module, &permission.Description, &permission.CreatedAt, &permission.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Permission{}, false, nil
		}
		return domain.Permission{}, false, err
	}
	permission.DeletedAt = scanNullTime(deletedAt)
	return permission, true, nil
}

func (r *PermissionRepo) List(ctx context.Context, page, pageSize int) ([]domain.Permission, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM permissions WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,code,name,module,description,created_at,updated_at,deleted_at FROM permissions WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Permission, 0)
	for rows.Next() {
		var item domain.Permission
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Module, &item.Description, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PermissionRepo) Update(ctx context.Context, id int64, permission domain.Permission) (domain.Permission, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE permissions SET code=?,name=?,module=?,description=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, permission.Code, permission.Name, permission.Module, permission.Description, time.Now(), id)
	if err != nil {
		return domain.Permission{}, err
	}
	permission.ID = id
	return permission, nil
}

func (r *PermissionRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE permissions SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}
