package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func placeholders(count int) string {
	if count <= 0 {
		return ""
	}
	values := make([]string, count)
	for index := range values {
		values[index] = "?"
	}
	return strings.Join(values, ",")
}

func inClause(count int) string {
	if count <= 0 {
		return "(NULL)"
	}
	return "(" + placeholders(count) + ")"
}

func scanNullString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

func scanNullInt64(value sql.NullInt64) *int64 {
	if value.Valid {
		result := value.Int64
		return &result
	}
	return nil
}

func scanNullTime(value sql.NullTime) *time.Time {
	if value.Valid {
		result := value.Time
		return &result
	}
	return nil
}

func joinColumns(columns []string) string { return strings.Join(columns, ",") }

func columnAssignments(columns []string) string {
	items := make([]string, 0, len(columns))
	for _, column := range columns {
		items = append(items, fmt.Sprintf("%s=?", column))
	}
	return strings.Join(items, ",")
}
