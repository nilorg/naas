package service

import (
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
)

type resource struct {
}

// Get resource
func (*resource) Get(id uint64) (resource *model.Resource, err error) {
	resource, err = dao.Resource.Select(store.NewDBContext(), id)
	return
}

// LoadPolicy 加载规则
func (*resource) LoadPolicy(resourceID uint64) (results []*gormadapter.CasbinRule, err error) {
	return dao.Resource.LoadPolicy(store.NewDBContext(), resourceID)
}

// ResourceAddWebRouteRequest ...
type ResourceAddWebRouteRequest struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

// AddWebRoute 添加web路由
func (*resource) AddWebRoute(resourceID uint64, req *ResourceAddWebRouteRequest) (err error) {
	v := &model.ResourceWebRoute{
		Name:       req.Name,
		Path:       req.Path,
		Method:     req.Method,
		ResourceID: resourceID,
	}
	err = dao.ResourceWebRoute.Insert(store.NewDBContext(), v)
	return
}
