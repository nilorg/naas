package model

// Role 角色
type Role struct {
	CodeModel
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ParentCode  string  `json:"parent_code"`
	ChildRoles  []*Role `json:"child_roles" gorm:"-"`
}

// RoleResourceWebFunction ...
type RoleResourceWebFunction struct {
	Model
	RoleCode              string `json:"role_code" gorm:"column:role_code"`
	ResourceWebFunctionID uint64 `json:"resource_web_function_id" gorm:"column:resource_web_function_id"`
}

// UserRole 用户权限
type UserRole struct {
	Model
	UserID   uint64 `json:"user_id" gorm:"column:user_id"`
	RoleCode string `json:"role_code" gorm:"column:role_code"`
}

// RoleResourceWebRoute 角色资源web路由
type RoleResourceWebRoute struct {
	Model
	RoleCode           string `json:"role_code" gorm:"column:role_code"`
	ResourceWebRouteID uint64 `json:"resource_web_route_id" gorm:"column:resource_web_route_id"`
}
