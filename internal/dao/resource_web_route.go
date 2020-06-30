package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// ResourceWebRouter ...
type ResourceWebRouter interface {
	Insert(ctx context.Context, resourceWebRoute *model.ResourceWebRoute) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (resourceWebRoute *model.ResourceWebRoute, err error)
	Update(ctx context.Context, resourceWebRoute *model.ResourceWebRoute) (err error)
}

type resourceWebRoute struct {
}

func (*resourceWebRoute) Insert(ctx context.Context, resourceWebRoute *model.ResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resourceWebRoute).Error
	return
}
func (*resourceWebRoute) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceWebRoute{}, id).Error
	return
}
func (*resourceWebRoute) Select(ctx context.Context, id uint64) (resourceWebRoute *model.ResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.ResourceWebRoute
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	resourceWebRoute = &dbResult
	return
}
func (*resourceWebRoute) Update(ctx context.Context, resourceWebRoute *model.ResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resourceWebRoute).Update(resourceWebRoute).Error
	return
}
