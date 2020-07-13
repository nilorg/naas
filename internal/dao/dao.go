package dao

import (
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/pkg/cache"
)

var (
	OAuth2Client         OAuth2Clienter        = &oauth2Client{}
	OAuth2ClientInfo     OAuth2ClientInfoer    = &oauth2ClientInfo{}
	OAuth2ClientScope    OAuth2ClientScoper    = &oauth2ClientScope{}
	OAuth2Scope          OAuth2Scoper          = &oauth2Scope{}
	Resource             Resourcer             = &resource{cache: cache.NewRedisCache(store.RedisClient, "naas_resource:")}
	ResourceWebRoute     ResourceWebRouter     = &resourceWebRoute{}
	Admin                Adminer               = &admin{}
	User                 Userer                = &user{}
	UserInfo             UserInfoer            = &userInfo{}
	Organization         Organizationer        = &organization{}
	OrganizationRole     OrganizationRoleer    = &organizationRole{}
	Role                 Roleer                = &role{}
	UserRole             UserRoleer            = &userRole{cache: cache.NewRedisCache(store.RedisClient, "naas_user_role:")}
	RoleResourceWebRoute RoleResourceWebRouter = &roleResourceWebRoute{}
)
