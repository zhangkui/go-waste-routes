package service

import "go-waste-routes/internal/domain"

var DefaultPermissions = []domain.Permission{
	{Code: "user.manage", Name: "User Manage", Module: "user"},
	{Code: "role.manage", Name: "Role Manage", Module: "role"},
	{Code: "customer.manage", Name: "Customer Manage", Module: "customer"},
	{Code: "plan.manage", Name: "Plan Manage", Module: "plan"},
	{Code: "route.manage", Name: "Route Manage", Module: "route"},
	{Code: "task.manage", Name: "Task Manage", Module: "task"},
	{Code: "weighing.manage", Name: "Weighing Manage", Module: "weighing"},
	{Code: "invoice.manage", Name: "Invoice Manage", Module: "invoice"},
}

var DefaultRole = domain.Role{Code: "operator", Name: "日常操作人员"}
var AdminRole = domain.Role{Code: "admin", Name: "业务管理员", BuiltIn: true}
