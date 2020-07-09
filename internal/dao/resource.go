package dao

import (
	"context"
	"fmt"

	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// Resourcer ...
type Resourcer interface {
	Insert(ctx context.Context, resource *model.Resource) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (resource *model.Resource, err error)
	Update(ctx context.Context, resource *model.Resource) (err error)
	LoadPolicy(ctx context.Context, resourceID uint64) (results []*gormadapter.CasbinRule, err error)
}

type resource struct {
}

func (*resource) Insert(ctx context.Context, resource *model.Resource) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resource).Error
	return
}
func (*resource) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Resource{}, id).Error
	return
}
func (*resource) Select(ctx context.Context, id uint64) (resource *model.Resource, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Resource
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	resource = &dbResult
	return
}
func (*resource) Update(ctx context.Context, resource *model.Resource) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resource).Update(resource).Error
	return
}

// LoadPolicy 加载规则
func (*resource) LoadPolicy(ctx context.Context, resourceID uint64) (results []*gormadapter.CasbinRule, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("v1 like ?", fmt.Sprintf("resource:%d%%", resourceID)).Find(&results).Error
	return
}
