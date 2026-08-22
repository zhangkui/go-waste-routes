package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type RoleRepo struct{ db *DB }

func NewRoleRepo(db *DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) Create(ctx context.Context, role domain.Role) (domain.Role, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO roles(code,name,description,built_in,created_at,updated_at) VALUES(?,?,?,?,?,?)`, role.Code, role.Name, role.Description, role.BuiltIn, time.Now(), time.Now())
	if err != nil {
		return domain.Role{}, err
	}
	id, _ := result.LastInsertId()
	role.ID = id
	return role, nil
}

func (r *RoleRepo) Get(ctx context.Context, id int64) (domain.Role, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,code,name,description,built_in,created_at,updated_at,deleted_at FROM roles WHERE id=? AND deleted_at IS NULL`, id)
	var role domain.Role
	var deletedAt sql.NullTime
	if err := row.Scan(&role.ID, &role.Code, &role.Name, &role.Description, &role.BuiltIn, &role.CreatedAt, &role.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Role{}, false, nil
		}
		return domain.Role{}, false, err
	}
	role.DeletedAt = scanNullTime(deletedAt)
	return role, true, nil
}

func (r *RoleRepo) List(ctx context.Context, page, pageSize int) ([]domain.Role, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,code,name,description,built_in,created_at,updated_at,deleted_at FROM roles WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.Role, 0)
	for rows.Next() {
		var item domain.Role
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.BuiltIn, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *RoleRepo) Update(ctx context.Context, id int64, role domain.Role) (domain.Role, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE roles SET code=?,name=?,description=?,built_in=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, role.Code, role.Name, role.Description, role.BuiltIn, time.Now(), id)
	if err != nil {
		return domain.Role{}, err
	}
	role.ID = id
	return role, nil
}

func (r *RoleRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE roles SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

func (r *RoleRepo) FindByCode(ctx context.Context, code string) (domain.Role, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,code,name,description,built_in,created_at,updated_at,deleted_at FROM roles WHERE code=? AND deleted_at IS NULL`, code)
	var role domain.Role
	var deletedAt sql.NullTime
	if err := row.Scan(&role.ID, &role.Code, &role.Name, &role.Description, &role.BuiltIn, &role.CreatedAt, &role.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Role{}, false, nil
		}
		return domain.Role{}, false, err
	}
	role.DeletedAt = scanNullTime(deletedAt)
	return role, true, nil
}

func (r *RoleRepo) ReplacePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id=?`, roleID); err != nil {
			return err
		}
		for _, permissionID := range permissionIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions(role_id,permission_id) VALUES(?,?)`, roleID, permissionID); err != nil {
				return err
			}
		}
		return nil
	})
}
