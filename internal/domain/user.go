package domain

import "time"

type User struct {
	Base
	Username     string       `json:"username"`
	PasswordHash string       `json:"-"`
	DisplayName  string       `json:"display_name"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	Status       string       `json:"status"`
	LastLoginAt  *time.Time   `json:"last_login_at,omitempty"`
	Roles        []Role       `json:"roles,omitempty"`
	Permissions  []Permission `json:"permissions,omitempty"`
}

const (
	UserStatusEnabled  = "enabled"
	UserStatusDisabled = "disabled"
)
