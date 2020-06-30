package service

import (
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
