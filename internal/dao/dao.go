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
	ResourceRoute        ResourceRouter
	ResourceMenu         ResourceMenuer
	User                 Userer
	UserInfo             UserInfoer
	UserThird            UserThirder
	Organization         Organizationer
	Role                 Roleer
	UserRole             UserRoleer
	UserOrganization     UserOrganizationer
	RoleResourceRelation RoleResourceRelationer
	ResourceAction       ResourceActioner
)

// Init 初始化...
func Init() {
	OAuth2Client = &oauth2Client{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_client:")}
	OAuth2ClientInfo = &oauth2ClientInfo{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_client_info:")}
	OAuth2ClientScope = &oauth2ClientScope{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_client_scope:")}
	OAuth2Scope = &oauth2Scope{cache: cache.NewRedisCache(store.RedisClient, "naas:oauth2_scope:")}
	ResourceRoute = &resourceRoute{}
	ResourceMenu = &resourceMenu{}
	User = &user{cache: cache.NewRedisCache(store.RedisClient, "naas:user:")}
	UserThird = &userThird{cache: cache.NewRedisCache(store.RedisClient, "naas:user_third:")}
	UserInfo = &userInfo{cache: cache.NewRedisCache(store.RedisClient, "naas:user_info:")}
	Organization = &organization{}
	Role = &role{}
	Resource = &resource{cache: cache.NewRedisCache(store.RedisClient, "naas:resource:")}
	UserRole = &userRole{cache: cache.NewRedisCache(store.RedisClient, "naas:user_role:")}
	UserOrganization = &userOrganization{cache: cache.NewRedisCache(store.RedisClient, "naas:user_organization:")}
	RoleResourceRelation = &roleResourceRelation{}
	ResourceAction = &resourceAction{}
}
