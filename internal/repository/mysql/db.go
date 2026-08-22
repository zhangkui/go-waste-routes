package mysql

import (
	"database/sql"
	"time"
)

type DB struct {
	Conn *sql.DB
}

func New(conn *sql.DB) *DB { return &DB{Conn: conn} }

func (db *DB) Ping() error {
	return db.Conn.Ping()
}

func (db *DB) WithTx(fn func(*sql.Tx) error) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) Now() time.Time { return time.Now().UTC() }
