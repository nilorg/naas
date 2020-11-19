package model

// Organization ...
type Organization struct {
	Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Code        Code   `json:"code" gorm:"column:code"`
	ParentID    ID     `json:"parent_id" gorm:"column:parent_id"`
}

// UserOrganization 用户组织
type UserOrganization struct {
	Model
	UserID         ID `json:"user_id" gorm:"column:user_id"`
	OrganizationID ID `json:"organization_id" gorm:"column:organization_id"`
}

// ResultOrganization 返回组织信息
type ResultOrganization struct {
	Organization       *Organization `json:"organization"`
	ParentOrganization *Organization `json:"parent_organization"`
}
