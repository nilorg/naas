package service

import (
	"context"

	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
)

type resource struct {
}

// Get resource
func (*resource) Get(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
	resource, err = dao.Resource.Select(ctx, id)
	return
}

// LoadPolicy 加载规则
func (*resource) LoadPolicy(ctx context.Context, resourceID model.ID) (results []*gormadapter.CasbinRule, err error) {
	return dao.Resource.LoadPolicy(ctx, resourceID)
}

// ResourceAddWebRouteRequest ...
type ResourceAddWebRouteRequest struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

// AddWebRoute 添加web路由
func (*resource) AddWebRoute(ctx context.Context, resourceID model.ID, req *ResourceAddWebRouteRequest) (err error) {
	v := &model.ResourceWebRoute{
		Name:       req.Name,
		Path:       req.Path,
		Method:     req.Method,
		ResourceID: resourceID,
	}
	err = dao.ResourceWebRoute.Insert(ctx, v)
	return
}
