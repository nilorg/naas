package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// ResourceRouter ...
type ResourceRouter interface {
	Insert(ctx context.Context, resourceRoute *model.ResourceRoute) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids ...model.ID) (err error)
	Select(ctx context.Context, id model.ID) (resourceRoute *model.ResourceRoute, err error)
	Update(ctx context.Context, resourceRoute *model.ResourceRoute) (err error)
	ListByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceRoute, err error)
	ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceRoute, err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceRoute, total int64, err error)
	ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceRoute, total int64, err error)
}

type resourceRoute struct {
}

func (*resourceRoute) Insert(ctx context.Context, resourceRoute *model.ResourceRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resourceRoute).Error
	return
}
func (*resourceRoute) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceRoute{}, id).Error
	return
}
func (*resourceRoute) Select(ctx context.Context, id model.ID) (resourceRoute *model.ResourceRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	resourceRoute = new(model.ResourceRoute)
	err = gdb.Model(resourceRoute).Where("id = ?", id).Take(resourceRoute).Error
	if err != nil {
		resourceRoute = nil
	}
	return
}
func (*resourceRoute) Update(ctx context.Context, resourceRoute *model.ResourceRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resourceRoute).Save(resourceRoute).Error
	return
}

func (*resourceRoute) ListByResourceID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceRoute{}).Where("resource_server_id = ?", resourceID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (*resourceRoute) ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ResourceRoute{}).Where("id in ?", ids).Find(&list).Error
	return
}

func (*resourceRoute) DeleteInIDs(ctx context.Context, ids ...model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in ?", ids).Delete(model.ResourceRoute{}).Error
	return
}

func (r *resourceRoute) ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceRoute, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceRoute{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceRoute) ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceRoute, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceRoute{}).Where("resource_server_id = ?", resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}
