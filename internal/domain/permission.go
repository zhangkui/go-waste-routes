package domain

type Permission struct {
	Base
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

type RolePermission struct {
	RoleID       int64 `json:"role_id"`
	PermissionID  int64 `json:"permission_id"`
}
