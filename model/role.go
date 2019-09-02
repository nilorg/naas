package model

// Role 角色
type Role struct {
	Model
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	WebFunctions []WebFunction `json:"web_functions" gorm:"-"`
}

// RoleWebFunction ...
type RoleWebFunction struct {
	RoleID        uint64 `json:"role_id" gorm:"column:role_id"`
	WebFunctionID uint64 `json:"web_function_id" gorm:"column:web_function_id"`
}

// UserRole 用户权限
type UserRole struct {
	Model
	UserID uint64 `json:"user_id" gorm:"column:user_id"`
	RoleID uint64 `json:"role_id" gorm:"column:role_id"`
}
