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
	// RoleResourceRelationTypeRoute 路由
	RoleResourceRelationTypeRoute RoleResourceRelationType = 1
	// RoleResourceRelationTypeMenu 菜单
	RoleResourceRelationTypeMenu RoleResourceRelationType = 2
	// RoleResourceRelationTypeAction 动作
	RoleResourceRelationTypeAction RoleResourceRelationType = 3
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
