package model

// Role 角色
type Role struct {
	CodeModel
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	ParentCode   string        `json:"parent_code"`
	WebFunctions []WebFunction `json:"web_functions" gorm:"-"`
}

// RoleWebFunction ...
type RoleWebFunction struct {
	RoleCode      string `json:"role_code" gorm:"column:role_code"`
	WebFunctionID uint64 `json:"web_function_id" gorm:"column:web_function_id"`
}

// UserRole 用户权限
type UserRole struct {
	Model
	UserID   uint64 `json:"user_id" gorm:"column:user_id"`
	RoleCode string `json:"role_code" gorm:"column:role_code"`
}
