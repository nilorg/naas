package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// ResourceMenuer ...
type ResourceMenuer interface {
	Insert(ctx context.Context, resourceMenu *model.ResourceMenu) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids ...model.ID) (err error)
	Select(ctx context.Context, id model.ID) (resourceMenu *model.ResourceMenu, err error)
	SelectByParentID(ctx context.Context, parentID model.ID) (resourceMenus []*model.ResourceMenu, err error)
	Update(ctx context.Context, resourceMenu *model.ResourceMenu) (err error)
	ListByResourceServerID(ctx context.Context, resourceServerID model.ID, limit int) (list []*model.ResourceMenu, err error)
	ListRootByResourceServerID(ctx context.Context, resourceServerID model.ID) (list []*model.ResourceMenu, err error)
	ListByResourceServerIDAndParentID(ctx context.Context, resourceServerID, parentID model.ID, limit int) (list []*model.ResourceMenu, err error)
	ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceMenu, err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceMenu, total int64, err error)
	ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceMenu, total int64, err error)
	CountByParentID(ctx context.Context, parentID model.ID) (count int64, err error)
}

type resourceMenu struct {
}

func (*resourceMenu) Insert(ctx context.Context, resourceMenu *model.ResourceMenu) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resourceMenu).Error
	return
}
func (*resourceMenu) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceMenu{}, id).Error
	return
}

func (*resourceMenu) Select(ctx context.Context, id model.ID) (resourceMenu *model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	resourceMenu = new(model.ResourceMenu)
	err = gdb.Model(resourceMenu).Where("id = ?", id).Take(resourceMenu).Error
	if err != nil {
		resourceMenu = nil
	}
	return
}

func (*resourceMenu) SelectByParentID(ctx context.Context, parentID model.ID) (resourceMenus []*model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.ResourceMenu{}).Where("parent_id = ?", parentID).Find(&resourceMenus).Error
	return
}

func (*resourceMenu) Update(ctx context.Context, resourceMenu *model.ResourceMenu) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resourceMenu).Save(resourceMenu).Error
	return
}

func (*resourceMenu) ListByResourceServerID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceMenu{}).Where("resource_server_id = ?", resourceID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (*resourceMenu) ListRootByResourceServerID(ctx context.Context, resourceServerID model.ID) (list []*model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceMenu{}).Where("(ISNULL(parent_id) OR parent_id = '') and resource_server_id = ?", resourceServerID)
	err = exp.Find(&list).Error
	return
}

func (*resourceMenu) ListByResourceServerIDAndParentID(ctx context.Context, resourceID, parentID model.ID, limit int) (list []*model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceMenu{}).Where("resource_server_id = ? and parent_id = ?", resourceID, parentID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (*resourceMenu) ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ResourceMenu{}).Where("id in ?", ids).Find(&list).Error
	return
}

func (*resourceMenu) DeleteInIDs(ctx context.Context, ids ...model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in ?", ids).Delete(model.ResourceMenu{}).Error
	return
}

func (r *resourceMenu) ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceMenu, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceMenu{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceMenu) ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceMenu, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceMenu{}).Where("resource_server_id = ?", resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceMenu) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var count int64
	count, err = r.count(ctx, query, args...)
	if count > 0 {
		exist = true
	}
	return
}

func (r *resourceMenu) CountByParentID(ctx context.Context, parentID model.ID) (count int64, err error) {
	return r.count(ctx, "parent_id = ?", parentID)
}

func (r *resourceMenu) count(ctx context.Context, query interface{}, args ...interface{}) (count int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ResourceMenu{}).Where(query, args...).Count(&count).Error
	return
}
