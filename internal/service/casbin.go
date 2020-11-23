package service

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type casbinService struct{}

func (cs *casbinService) InitLoadAllPolicy() {
	ctx := contexts.WithContext(context.Background())
	cs.initRoleForUser(ctx)
	cs.initRoleResourceWebRoute(ctx)
}

func (*casbinService) initRoleResourceWebRoute(ctx context.Context) {
	var (
		err       error
		relations []*model.RoleResourceRelation
		flag      bool
	)
	// 获取所有角色对应的资源路由
	relations, err = dao.RoleResourceRelation.SelectAll(ctx, model.RoleResourceRelationTypeWebRoute)
	if err != nil {
		logrus.Errorf("dao.RoleResourceWebRoute.SelectAll: %s", err)
		return
	}
	for _, relation := range relations {
		var resourceWebRoute *model.ResourceWebRoute
		resourceWebRoute, err = dao.ResourceWebRoute.Select(ctx, relation.RelationID)
		if err != nil {
			logrus.Errorf("dao.ResourceWebRoute.Select: %s", err)
			return
		}
		sub, dom, obj, act := formatWebRoutePolicy(relation.RoleCode, resourceWebRoute)
		flag, err = casbin.Enforcer.AddPolicy(sub, dom, obj, act)
		if err != nil {
			logrus.Errorf("casbin.Enforcer.AddPolicy: %s", err)
			continue
		}
		logrus.Infof("[ResourceWebRoute]casbin.Enforcer.AddPolicy-Flag: %v", flag)
	}
	err = casbin.Enforcer.SavePolicy()
	if err != nil {
		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
	}
}

func (*casbinService) initRoleForUser(ctx context.Context) {
	var (
		err       error
		userRoles []*model.UserRole
		flag      bool
	)
	// 获取所有用户的角色
	userRoles, err = dao.UserRole.SelectAll(ctx)
	if err != nil {
		logrus.Errorf("dao.UserRole.SelectAll: %s", err)
		return
	}
	for _, userRole := range userRoles {
		user, role, domain := formatRoleForUserInDomain(userRole.UserID, userRole.OrganizationID, userRole.RoleCode)
		flag, err = casbin.Enforcer.AddRoleForUserInDomain(user, role, domain)
		if err != nil {
			logrus.Errorf("casbin.Enforcer.AddRoleForUserInDomain: %s", err)
			err = nil
			continue
		}
		logrus.Infof("casbin.Enforcer.AddRoleForUserInDomain-Flag: %v", flag)
	}
	err = casbin.Enforcer.SavePolicy()
	if err != nil {
		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
	}
}

func formatWebRoutePolicy(roleCode model.Code, resourceWebRoute *model.ResourceWebRoute) (sub, dom, obj, act string) {
	sub = fmt.Sprintf("role:%s", roleCode)                                        // 希望访问资源的角色
	dom = fmt.Sprintf("resource:%d:web_route", resourceWebRoute.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = resourceWebRoute.Path                                                   // 要访问的资源
	act = resourceWebRoute.Method                                                 // 用户对资源执行的操作
	return
}

func formatRoleForUserInDomain(userID, organizationID model.ID, roleCode model.Code) (user, role, domain string) {
	user = fmt.Sprintf("user:%v", userID)
	role = fmt.Sprintf("role:%v", roleCode)
	domain = fmt.Sprintf("organization:%v", organizationID)
	return
}

func formatWebMenuPolicy(roleCode model.Code, resourceWebMenu *model.ResourceWebMenu) (sub, dom, obj, act string) {
	// Enforcer.AddPolicy("role:reader", "domain1", "data1", "read")
	sub = fmt.Sprintf("role:%s", roleCode)                                      // 希望访问资源的角色
	dom = fmt.Sprintf("resource:%d:web_menu", resourceWebMenu.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = fmt.Sprintf("web_menu:%v", resourceWebMenu.ID)                        // 要访问的资源
	act = "show"                                                                // 角色对资源执行的操作
	return
}

// ListResourceWebRoutePaged ...
func (cs *casbinService) ListResourceWebRoutePagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceWebRoute, total int64, err error) {
	list, total, err = dao.ResourceWebRoute.ListPagedByResourceServerID(ctx, start, limit, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// ListRoleResourceWebRouteByRoleCodeAndResourceServerID ...
func (cs *casbinService) ListRoleResourceWebRouteByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	list, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCodeAndResourceServerID(ctx, model.RoleResourceRelationTypeWebRoute, roleCode, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// CasbinAddResourceWebRouteModel ...
type CasbinAddResourceWebRouteModel struct {
	ResourceWebRouteIDs []model.ID `json:"resource_web_route_ids"`
	ResourceServerID    model.ID   `json:"resource_server_id"`
}

// AddResourceWebRoute 添加web路由资源角色
func (cs *casbinService) AddResourceWebRoute(ctx context.Context, roleCode model.Code, create *CasbinAddResourceWebRouteModel) (err error) {
	err = Role.AddRoleResourceRelation(ctx, roleCode, model.RoleResourceRelationTypeWebRoute, create.ResourceServerID, create.ResourceWebRouteIDs...)
	return
}

// CasbinAddResourceWebMenuModel ...
type CasbinAddResourceWebMenuModel struct {
	ResourceWebMenuIDs []model.ID `json:"resource_web_menu_ids"`
	ResourceServerID   model.ID   `json:"resource_server_id"`
}

// AddResourceWebMenu 添加web菜单资源角色
func (cs *casbinService) AddResourceWebMenu(ctx context.Context, roleCode model.Code, create *CasbinAddResourceWebMenuModel) (err error) {
	err = Role.AddRoleResourceRelation(ctx, roleCode, model.RoleResourceRelationTypeWebMenu, create.ResourceServerID, create.ResourceWebMenuIDs...)
	return
}

// ListResourceWebMenuPagedByResourceServerID ...
func (cs *casbinService) ListResourceWebMenuPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceWebMenu, total int64, err error) {
	list, total, err = dao.ResourceWebMenu.ListPagedByResourceServerID(ctx, start, limit, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// ListRoleResourceWebMenuByRoleCodeAndResourceServerID ...
func (cs *casbinService) ListRoleResourceWebMenuByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	list, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCodeAndResourceServerID(ctx, model.RoleResourceRelationTypeWebMenu, roleCode, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}
