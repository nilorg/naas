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
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (resourceWebRoute *model.ResourceWebRoute, err error)
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
func (*resourceWebRoute) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceWebRoute{}, id).Error
	return
}
func (*resourceWebRoute) Select(ctx context.Context, id model.ID) (resourceWebRoute *model.ResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	resourceWebRoute = new(model.ResourceWebRoute)
	err = gdb.Model(resourceWebRoute).Where("id = ?", id).Scan(resourceWebRoute).Error
	if err != nil {
		resourceWebRoute = nil
		return
	}
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
