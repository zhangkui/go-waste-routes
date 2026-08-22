package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-waste-routes/internal/domain"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `INSERT INTO users(username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at) VALUES(?,?,?,?,?,?,?,?,?,?)`
	result, err := r.db.Conn.ExecContext(ctx, query, user.Username, user.PasswordHash, user.DisplayName, user.Phone, user.Email, user.Status, user.LastLoginAt, time.Now(), time.Now(), nil)
	if err != nil {
		return domain.User{}, err
	}
	id, _ := result.LastInsertId()
	user.ID = id
	return user, nil
}

func (r *UserRepo) Get(ctx context.Context, id int64) (domain.User, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at FROM users WHERE id=? AND deleted_at IS NULL`, id)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return domain.User{}, false, nil
	}
	if err != nil {
		return domain.User{}, false, err
	}
	return user, true, nil
}

func (r *UserRepo) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	total := int64(0)
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at FROM users WHERE deleted_at IS NULL ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.User, 0)
	for rows.Next() {
		item, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, id int64, user domain.User) (domain.User, error) {
	query := `UPDATE users SET username=?,password_hash=?,display_name=?,phone=?,email=?,status=?,last_login_at=?,updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.Conn.ExecContext(ctx, query, user.Username, user.PasswordHash, user.DisplayName, user.Phone, user.Email, user.Status, user.LastLoginAt, time.Now(), id)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = id
	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE users SET deleted_at=? WHERE id=? AND deleted_at IS NULL`, time.Now(), id)
	return err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (domain.User, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at FROM users WHERE username=? AND deleted_at IS NULL`, username)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return domain.User{}, false, nil
	}
	if err != nil {
		return domain.User{}, false, err
	}
	return user, true, nil
}

func (r *UserRepo) ReplaceRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	return r.db.WithTx(func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id=?`, userID); err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_roles(user_id,role_id) VALUES(?,?)`, userID, roleID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *UserRepo) SetStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE users SET status=?,updated_at=? WHERE id=?`, status, time.Now(), id)
	return err
}

func (r *UserRepo) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	_, err := r.db.Conn.ExecContext(ctx, `UPDATE users SET password_hash=?,updated_at=? WHERE id=?`, passwordHash, time.Now(), id)
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (domain.User, bool, error) {
	row := r.db.Conn.QueryRowContext(ctx, `SELECT id,username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at FROM users WHERE email=? AND deleted_at IS NULL`, email)
	user, err := scanUser(row)
	if err == sql.ErrNoRows {
		return domain.User{}, false, nil
	}
	if err != nil {
		return domain.User{}, false, err
	}
	return user, true, nil
}

func (r *UserRepo) Search(ctx context.Context, keyword string, limit int) ([]domain.User, error) {
	keyword = fmt.Sprintf("%%%s%%", strings.TrimSpace(keyword))
	rows, err := r.db.Conn.QueryContext(ctx, `SELECT id,username,password_hash,display_name,phone,email,status,last_login_at,created_at,updated_at,deleted_at FROM users WHERE deleted_at IS NULL AND (username LIKE ? OR display_name LIKE ? OR email LIKE ?) ORDER BY id DESC LIMIT ?`, keyword, keyword, keyword, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.User, 0)
	for rows.Next() {
		item, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *UserRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if err := r.db.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE status=? AND deleted_at IS NULL`, status).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func scanUser(scanner interface{ Scan(...any) error }) (domain.User, error) {
	var user domain.User
	var lastLogin sql.NullTime
	var deletedAt sql.NullTime
	if err := scanner.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.Phone, &user.Email, &user.Status, &lastLogin, &user.CreatedAt, &user.UpdatedAt, &deletedAt); err != nil {
		return domain.User{}, err
	}
	user.LastLoginAt = scanNullTime(lastLogin)
	user.DeletedAt = scanNullTime(deletedAt)
	return user, nil
}
