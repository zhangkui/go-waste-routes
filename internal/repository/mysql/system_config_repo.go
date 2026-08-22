package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type SystemConfigRepo struct{ db *DB }

func NewSystemConfigRepo(db *DB) *SystemConfigRepo { return &SystemConfigRepo{db: db} }

func (r *SystemConfigRepo) Create(ctx context.Context, config domain.SystemConfig) (domain.SystemConfig, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO system_configs(config_key,config_value,value_type,module,description,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, config.ConfigKey, config.ConfigValue, config.ValueType, config.Module, config.Description, time.Now(), time.Now())
	if err != nil {
		return domain.SystemConfig{}, err
	}
	id, _ := result.LastInsertId()
	config.ID = id
	return config, nil
}

func (r *SystemConfigRepo) Get(ctx context.Context, id int64) (domain.SystemConfig, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,config_key,config_value,value_type,module,description,created_at,updated_at,deleted_at FROM system_configs WHERE id=? AND deleted_at IS NULL`, id)
	var config domain.SystemConfig
	var deletedAt sql.NullTime
	if err := row.Scan(&config.ID, &config.ConfigKey, &config.ConfigValue, &config.ValueType, &config.Module, &config.Description, &config.CreatedAt, &config.UpdatedAt, &deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.SystemConfig{}, false, nil
		}
		return domain.SystemConfig{}, false, err
	}
	config.DeletedAt = scanNullTime(deletedAt)
	return config, true, nil
}

func (r *SystemConfigRepo) List(ctx context.Context, page, pageSize int) ([]domain.SystemConfig, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM system_configs WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,config_key,config_value,value_type,module,description,created_at,updated_at,deleted_at FROM system_configs WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.SystemConfig, 0)
	for rows.Next() {
		var item domain.SystemConfig
		var deletedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.ConfigKey, &item.ConfigValue, &item.ValueType, &item.Module, &item.Description, &item.CreatedAt, &item.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}
		item.DeletedAt = scanNullTime(deletedAt)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *SystemConfigRepo) Update(ctx context.Context, id int64, config domain.SystemConfig) (domain.SystemConfig, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE system_configs SET config_key=?,config_value=?,value_type=?,module=?,description=?,updated_at=? WHERE id=? AND deleted_at IS NULL`, config.ConfigKey, config.ConfigValue, config.ValueType, config.Module, config.Description, time.Now(), id)
	if err != nil {
		return domain.SystemConfig{}, err
	}
	config.ID = id
	return config, nil
}

func (r *SystemConfigRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE system_configs SET deleted_at=? WHERE id=?`, time.Now(), id)
	return err
}

