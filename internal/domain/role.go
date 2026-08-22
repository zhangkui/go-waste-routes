package domain

type Role struct {
	Base
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	BuiltIn     bool         `json:"built_in"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type UserRole struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}
