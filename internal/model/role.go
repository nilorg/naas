package model

// Role 角色
type Role struct {
	CodeModel
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ParentCode  Code    `json:"parent_code"`
	ChildRoles  []*Role `json:"child_roles" gorm:"-"`
}

// RoleResourceWebFunction ...
type RoleResourceWebFunction struct {
	Model
	RoleCode              Code `json:"role_code" gorm:"column:role_code"`
	ResourceWebFunctionID ID   `json:"resource_web_function_id" gorm:"column:resource_web_function_id"`
}

// UserRole 用户权限
type UserRole struct {
	Model
	UserID   ID   `json:"user_id" gorm:"column:user_id"`
	RoleCode Code `json:"role_code" gorm:"column:role_code"`
}

// RoleResourceWebRoute 角色资源web路由
type RoleResourceWebRoute struct {
	Model
	RoleCode           Code `json:"role_code" gorm:"column:role_code"`
	ResourceWebRouteID ID   `json:"resource_web_route_id" gorm:"column:resource_web_route_id"`
}

// ResultRole 角色
type ResultRole struct {
	Role       *Role `json:"role"`
	ParentRole *Role `json:"parent_role"`
}
