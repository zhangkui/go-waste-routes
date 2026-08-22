package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Migrate(ctx context.Context, db *sql.DB, directory string) error {
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(128) PRIMARY KEY, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)"); err != nil {
		return err
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	var names []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		if err := migrateOne(ctx, db, directory, name); err != nil {
			return err
		}
	}
	return nil
}

func migrateOne(ctx context.Context, db *sql.DB, directory, name string) error {
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version=?", name).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	content, err := os.ReadFile(filepath.Join(directory, name))
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, string(content)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("migration %s: %w", name, err)
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES(?)", name); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
