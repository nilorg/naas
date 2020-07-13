package dao

import (
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/pkg/cache"
)

var (
	OAuth2Client         OAuth2Clienter
	OAuth2ClientInfo     OAuth2ClientInfoer
	OAuth2ClientScope    OAuth2ClientScoper
	OAuth2Scope          OAuth2Scoper
	Resource             Resourcer
	ResourceWebRoute     ResourceWebRouter
	Admin                Adminer
	User                 Userer
	UserInfo             UserInfoer
	Organization         Organizationer
	OrganizationRole     OrganizationRoleer
	Role                 Roleer
	UserRole             UserRoleer
	RoleResourceWebRoute RoleResourceWebRouter
)

// Init 初始化...
func Init() {
	OAuth2Client = &oauth2Client{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_client:")}
	OAuth2ClientInfo = &oauth2ClientInfo{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_client_info:")}
	OAuth2ClientScope = &oauth2ClientScope{}
	OAuth2Scope = &oauth2Scope{}
	ResourceWebRoute = &resourceWebRoute{}
	Admin = &admin{}
	User = &user{cache: cache.NewRedisCache(store.RedisClient, "naas:user:")}
	UserInfo = &userInfo{cache: cache.NewRedisCache(store.RedisClient, "naas:user_info:")}
	Organization = &organization{}
	OrganizationRole = &organizationRole{}
	Role = &role{}
	Resource = &resource{cache: cache.NewRedisCache(store.RedisClient, "naas:resource:")}
	UserRole = &userRole{cache: cache.NewRedisCache(store.RedisClient, "naas:user_role:")}
	RoleResourceWebRoute = &roleResourceWebRoute{}
}
