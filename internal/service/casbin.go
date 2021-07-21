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
	// cs.initRoleForUser(ctx)
	cs.initRoleResourceRoute(ctx)
	cs.initScopeResourceRoute(ctx)
}

func (*casbinService) initRoleResourceRoute(ctx context.Context) {
	var (
		err       error
		relations []*model.RoleResourceRelation
		flag      bool
	)
	// 获取所有角色对应的资源路由
	relations, err = dao.RoleResourceRelation.SelectAll(ctx, model.RoleResourceRelationTypeRoute)
	if err != nil {
		logrus.Errorf("dao.RoleResourceRelation.SelectAll: %s", err)
		return
	}
	for _, relation := range relations {
		var resourceRoute *model.ResourceRoute
		resourceRoute, err = dao.ResourceRoute.Select(ctx, relation.RelationID)
		if err != nil {
			logrus.Errorf("dao.ResourceRoute.Select: %s", err)
			return
		}
		sub, dom, obj, act := formatRoutePolicyForRole(relation.RoleCode, resourceRoute)
		flag, err = casbin.Enforcer.AddPolicy(sub, dom, obj, act)
		if err != nil {
			logrus.Errorf("casbin.Enforcer.AddPolicy: %s", err)
			continue
		}
		logrus.Infof("[ResourceRoute]casbin.Enforcer.AddPolicy-Flag: %v", flag)
	}
	err = casbin.Enforcer.SavePolicy()
	if err != nil {
		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
	}
}

func (*casbinService) initScopeResourceRoute(ctx context.Context) {
	var (
		err       error
		relations []*model.ScopeResourceRelation
		flag      bool
	)
	// 获取所有范围对应的资源路由
	relations, err = dao.ScopeResourceRelation.SelectAll(ctx, model.ScopeResourceRelationTypeRoute)
	if err != nil {
		logrus.Errorf("dao.ScopeResourceRelation.SelectAll: %s", err)
		return
	}
	for _, relation := range relations {
		var resourceRoute *model.ResourceRoute
		resourceRoute, err = dao.ResourceRoute.Select(ctx, relation.RelationID)
		if err != nil {
			logrus.Errorf("dao.ResourceRoute.Select: %s", err)
			return
		}
		sub, dom, obj, act := formatRoutePolicyForScope(relation.ScopeCode, resourceRoute)
		flag, err = casbin.Enforcer.AddPolicy(sub, dom, obj, act)
		if err != nil {
			logrus.Errorf("casbin.Enforcer.AddPolicy: %s", err)
			continue
		}
		logrus.Infof("[ResourceRoute]casbin.Enforcer.AddPolicy-Flag: %v", flag)
	}
	err = casbin.Enforcer.SavePolicy()
	if err != nil {
		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
	}
}

// func (*casbinService) initRoleForUser(ctx context.Context) {
// 	var (
// 		err       error
// 		userRoles []*model.UserRole
// 		flag      bool
// 	)
// 	// 获取所有用户的角色
// 	userRoles, err = dao.UserRole.SelectAll(ctx)
// 	if err != nil {
// 		logrus.Errorf("dao.UserRole.SelectAll: %s", err)
// 		return
// 	}
// 	for _, userRole := range userRoles {
// 		user, role, domain := formatRoleForUserInDomain(userRole.UserID, userRole.OrganizationID, userRole.RoleCode)
// 		flag, err = casbin.Enforcer.AddRoleForUserInDomain(user, role, domain)
// 		if err != nil {
// 			logrus.Errorf("casbin.Enforcer.AddRoleForUserInDomain: %s", err)
// 			err = nil
// 			continue
// 		}
// 		logrus.Infof("casbin.Enforcer.AddRoleForUserInDomain-Flag: %v", flag)
// 	}
// 	err = casbin.Enforcer.SavePolicy()
// 	if err != nil {
// 		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
// 	}
// }

func formatRoutePolicyForRole(roleCode model.Code, resourceRoute *model.ResourceRoute) (sub, dom, obj, act string) {
	sub = fmt.Sprintf("role:%s", roleCode)                                 // 希望访问资源的角色
	dom = fmt.Sprintf("resource:%d:route", resourceRoute.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = resourceRoute.Path                                               // 要访问的资源
	act = resourceRoute.Method                                             // 用户对资源执行的操作
	return
}

func formatRoutePolicyForScope(scopeCode model.Code, resourceRoute *model.ResourceRoute) (sub, dom, obj, act string) {
	sub = fmt.Sprintf("scope:%s", scopeCode)                               // 希望访问资源的范围
	dom = fmt.Sprintf("resource:%d:route", resourceRoute.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = resourceRoute.Path                                               // 要访问的资源
	act = resourceRoute.Method                                             // 用户对资源执行的操作
	return
}

func formatRoleForUserInDomain(userID, organizationID model.ID, roleCode model.Code) (user, role, domain string) {
	user = fmt.Sprintf("user:%v", userID)
	role = fmt.Sprintf("role:%v", roleCode)
	domain = fmt.Sprintf("organization:%v", organizationID)
	return
}

func formatMenuPolicyForRole(roleCode model.Code, resourceMenu *model.ResourceMenu) (sub, dom, obj, act string) {
	// Enforcer.AddPolicy("role:reader", "domain1", "data1", "read")
	sub = fmt.Sprintf("role:%s", roleCode)                               // 希望访问资源的角色
	dom = fmt.Sprintf("resource:%d:menu", resourceMenu.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = fmt.Sprintf("menu:%v", resourceMenu.ID)                        // 要访问的资源
	act = "show"                                                         // 角色对资源执行的操作
	return
}

func formatActionPolicyForRole(roleCode model.Code, resourceAction *model.ResourceAction) (sub, dom, obj, act string) {
	// Enforcer.AddPolicy("role:reader", "domain1", "data1", "read")
	sub = fmt.Sprintf("role:%s", roleCode)                                   // 希望访问资源的角色
	dom = fmt.Sprintf("resource:%d:action", resourceAction.ResourceServerID) // 域/域租户,这里以资源为单位
	obj = fmt.Sprintf("object:%v", resourceAction.Group)                     // 要访问的资源
	act = model.ConvertCodeToString(resourceAction.Code)                     // 角色对资源执行的操作
	return
}

// ListResourceRoutePagedByResourceServerID ...
func (cs *casbinService) ListResourceRoutePagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceRoute, total int64, err error) {
	list, total, err = dao.ResourceRoute.ListPagedByResourceServerID(ctx, start, limit, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// ListRoleResourceRouteByRoleCodeAndResourceServerID ...
func (cs *casbinService) ListRoleResourceRouteByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	list, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCodeAndResourceServerID(ctx, model.RoleResourceRelationTypeRoute, roleCode, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// CasbinAddResourceRouteModel ...
type CasbinAddResourceRouteModel struct {
	ResourceRouteIDs []model.ID `json:"resource_route_ids"`
	ResourceServerID model.ID   `json:"resource_server_id"`
}

// AddResourceRoute 添加路由资源角色
func (cs *casbinService) AddResourceRoute(ctx context.Context, roleCode model.Code, create *CasbinAddResourceRouteModel) (err error) {
	err = Role.AddRoleResourceRelation(ctx, roleCode, model.RoleResourceRelationTypeRoute, create.ResourceServerID, create.ResourceRouteIDs...)
	return
}

// CasbinAddResourceMenuModel ...
type CasbinAddResourceMenuModel struct {
	ResourceMenuIDs  []model.ID `json:"resource_menu_ids"`
	ResourceServerID model.ID   `json:"resource_server_id"`
}

// AddResourceMenu 添加菜单资源角色
func (cs *casbinService) AddResourceMenu(ctx context.Context, roleCode model.Code, create *CasbinAddResourceMenuModel) (err error) {
	err = Role.AddRoleResourceRelation(ctx, roleCode, model.RoleResourceRelationTypeMenu, create.ResourceServerID, create.ResourceMenuIDs...)
	return
}

// ListResourceMenuPagedByResourceServerID ...
func (cs *casbinService) ListResourceMenuPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceMenu, total int64, err error) {
	list, total, err = dao.ResourceMenu.ListPagedByResourceServerID(ctx, start, limit, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// ListRoleResourceMenuByRoleCodeAndResourceServerID ...
func (cs *casbinService) ListRoleResourceMenuByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	list, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCodeAndResourceServerID(ctx, model.RoleResourceRelationTypeMenu, roleCode, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// CasbinAddResourceActionModel ...
type CasbinAddResourceActionModel struct {
	ResourceActionIDs []model.ID `json:"resource_action_ids"`
	ResourceServerID  model.ID   `json:"resource_server_id"`
}

// AddResourceAction 添加资源动作
func (cs *casbinService) AddResourceAction(ctx context.Context, roleCode model.Code, create *CasbinAddResourceActionModel) (err error) {
	err = Role.AddRoleResourceRelation(ctx, roleCode, model.RoleResourceRelationTypeAction, create.ResourceServerID, create.ResourceActionIDs...)
	return
}

// ListResourceActionPagedByResourceServerID ...
func (cs *casbinService) ListResourceActionPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceAction, total int64, err error) {
	list, total, err = dao.ResourceAction.ListPagedByResourceServerID(ctx, resourceServerID, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}

// ListRoleResourceActionByRoleCodeAndResourceServerID ...
func (cs *casbinService) ListRoleResourceActionByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceRelation, err error) {
	list, err = dao.RoleResourceRelation.ListByRelationTypeAndRoleCodeAndResourceServerID(ctx, model.RoleResourceRelationTypeAction, roleCode, resourceServerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}
	return
}
