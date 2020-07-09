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
