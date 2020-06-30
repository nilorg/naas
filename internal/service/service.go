package service

var (
	// Admin 管理员服务
	Admin = &admin{}
	// User 用户服务
	User = &user{}
	// OAuth2 service
	OAuth2 = &oauth2{}
	// Role service
	Role = &role{}
	// Casbin service
	Casbin = &casbinService{}
	// Resource service
	Resource = &resource{}
)
