package model

// Organization ...
type Organization struct {
	Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Code        string `json:"code" gorm:"column:code"`
}

// OrganizationRole 组织权限
type OrganizationRole struct {
	Model
	UserID uint64 `json:"user_id" gorm:"column:user_id"`
	RoleID uint64 `json:"role_id" gorm:"column:role_id"`
}
