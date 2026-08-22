package mysql

import (
	"context"
	"database/sql"
	"time"

	"go-waste-routes/internal/domain"
)

type AuditRepo struct{ db *DB }

func NewAuditRepo(db *DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Create(ctx context.Context, log domain.AuditLog) (domain.AuditLog, error) {
	result, err := r.db.Conn.ExecContext(ctx, `INSERT INTO audit_logs(user_id,username,action,resource_type,resource_id,before_value,after_value,ip_address,user_agent,request_id,created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, log.UserID, log.Username, log.Action, log.ResourceType, log.ResourceID, log.BeforeValue, log.AfterValue, log.IPAddress, log.UserAgent, log.RequestID, time.Now())
	if err != nil {
		return domain.AuditLog{}, err
	}
	id, _ := result.LastInsertId()
	log.ID = id
	log.CreatedAt = time.Now()
	return log, nil
}

func (r *AuditRepo) Get(ctx context.Context, id int64) (domain.AuditLog, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,user_id,username,action,resource_type,resource_id,before_value,after_value,ip_address,user_agent,request_id,created_at FROM audit_logs WHERE id=?`, id)
	var log domain.AuditLog
	var userID sql.NullInt64
	if err := row.Scan(&log.ID, &userID, &log.Username, &log.Action, &log.ResourceType, &log.ResourceID, &log.BeforeValue, &log.AfterValue, &log.IPAddress, &log.UserAgent, &log.RequestID, &log.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.AuditLog{}, false, nil
		}
		return domain.AuditLog{}, false, err
	}
	log.UserID = scanNullInt64(userID)
	return log, true, nil
}

func (r *AuditRepo) List(ctx context.Context, page, pageSize int) ([]domain.AuditLog, int64, error) {
	var total int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_logs`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,user_id,username,action,resource_type,resource_id,before_value,after_value,ip_address,user_agent,request_id,created_at FROM audit_logs ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.AuditLog, 0)
	for rows.Next() {
		var item domain.AuditLog
		var userID sql.NullInt64
		if err := rows.Scan(&item.ID, &userID, &item.Username, &item.Action, &item.ResourceType, &item.ResourceID, &item.BeforeValue, &item.AfterValue, &item.IPAddress, &item.UserAgent, &item.RequestID, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		item.UserID = scanNullInt64(userID)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *AuditRepo) Update(ctx context.Context, id int64, log domain.AuditLog) (domain.AuditLog, error) {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE audit_logs SET username=?,action=?,resource_type=?,resource_id=?,before_value=?,after_value=?,ip_address=?,user_agent=?,request_id=? WHERE id=?`, log.Username, log.Action, log.ResourceType, log.ResourceID, log.BeforeValue, log.AfterValue, log.IPAddress, log.UserAgent, log.RequestID, id)
	if err != nil {
		return domain.AuditLog{}, err
	}
	log.ID = id
	return log, nil
}

func (r *AuditRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `DELETE FROM audit_logs WHERE id=?`, id)
	return err
}

