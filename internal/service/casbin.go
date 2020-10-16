package service

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/sirupsen/logrus"
)

type casbinService struct{}

func (cs *casbinService) InitLoadAllPolicy() {
	ctx := contexts.WithContext(context.Background())
	cs.initRoleResourceWebRoute(ctx)
}

func (*casbinService) initRoleResourceWebRoute(ctx context.Context) {
	var (
		err                   error
		roleResourceWebRoutes []*model.RoleResourceWebRoute
		resourceWebRoute      *model.ResourceWebRoute
		flag                  bool
	)
	// 获取所有角色对应的资源路由
	roleResourceWebRoutes, err = dao.RoleResourceWebRoute.SelectAll(ctx)
	if err != nil {
		logrus.Errorf("dao.RoleResourceWebRoute.SelectAll: %s", err)
		return
	}
	for _, roleResourceWebRoute := range roleResourceWebRoutes {
		resourceWebRoute, err = dao.ResourceWebRoute.Select(ctx, roleResourceWebRoute.ResourceWebRouteID)
		if err != nil {
			logrus.Errorf("dao.ResourceWebRoute.Select: %s", err)
			return
		}
		sub, dom, obj, act := formatPolicy(roleResourceWebRoute.RoleCode, resourceWebRoute)
		flag, err = casbin.Enforcer.AddPolicy(sub, dom, obj, act)
		if err != nil {
			logrus.Errorf("casbin.Enforcer.AddPolicy: %s", err)
			continue
		}
		logrus.Infof("casbin.Enforcer.AddPolicy-Flag: %v", flag)
	}
	err = casbin.Enforcer.SavePolicy()
	if err != nil {
		logrus.Errorf("casbin.Enforcer.SavePolicy: %s", err)
	}
}

func formatPolicy(roleCode model.Code, roleResourceWebRoute *model.ResourceWebRoute) (sub, dom, obj, act string) {
	sub = fmt.Sprintf("role:%s", roleCode)                                      // 希望访问资源的用户
	dom = fmt.Sprintf("resource:%d:web_route", roleResourceWebRoute.ResourceID) // 域/域租户,这里以资源为单位
	obj = roleResourceWebRoute.Path                                             // 要访问的资源
	act = roleResourceWebRoute.Method                                           // 用户对资源执行的操作
	return
}
