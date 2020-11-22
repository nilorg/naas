package model

// Role 角色
type Role struct {
	CodeModel
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	ParentCode     Code    `json:"parent_code"`
	OrganizationID ID      `json:"organization_id" gorm:"column:organization_id"`
	ChildRoles     []*Role `json:"child_roles" gorm:"-"`
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
	UserID         ID   `json:"user_id" gorm:"column:user_id"`
	RoleCode       Code `json:"role_code" gorm:"column:role_code"`
	OrganizationID ID   `json:"organization_id" gorm:"column:organization_id"`
}

// RoleResourceRelationType 角色资源关系类型
type RoleResourceRelationType int

var (
	// RoleResourceRelationTypeWebRoute web路由
	RoleResourceRelationTypeWebRoute RoleResourceRelationType = 1
	// RoleResourceRelationTypeWebMenu web菜单
	RoleResourceRelationTypeWebMenu RoleResourceRelationType = 2
)

// RoleResourceRelation 角色资源关系表
type RoleResourceRelation struct {
	Model
	RoleCode         Code                     `json:"role_code" gorm:"column:role_code"`
	RelationType     RoleResourceRelationType `json:"relation_type" gorm:"column:relation_type"`
	RelationID       ID                       `json:"relation_id" gorm:"column:relation_id"`
	ResourceServerID ID                       `json:"resource_server_id" gorm:"column:resource_server_id"`
}

// ResultRole 角色
type ResultRole struct {
	Role         *Role         `json:"role"`
	ParentRole   *Role         `json:"parent_role"`
	Organization *Organization `json:"organization"`
}
