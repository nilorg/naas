package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// ResourceWebMenuer ...
type ResourceWebMenuer interface {
	Insert(ctx context.Context, resourceWebMenu *model.ResourceWebMenu) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids ...model.ID) (err error)
	Select(ctx context.Context, id model.ID) (resourceWebMenu *model.ResourceWebMenu, err error)
	SelectByParentID(ctx context.Context, parentID model.ID) (resourceWebMenus []*model.ResourceWebMenu, err error)
	Update(ctx context.Context, resourceWebMenu *model.ResourceWebMenu) (err error)
	ListByResourceServerID(ctx context.Context, resourceServerID model.ID, limit int) (list []*model.ResourceWebMenu, err error)
	ListRootByResourceServerID(ctx context.Context, resourceServerID model.ID) (list []*model.ResourceWebMenu, err error)
	ListByResourceServerIDAndParentID(ctx context.Context, resourceServerID, parentID model.ID, limit int) (list []*model.ResourceWebMenu, err error)
	ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceWebMenu, err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceWebMenu, total int64, err error)
	ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceWebMenu, total int64, err error)
	CountByParentID(ctx context.Context, parentID model.ID) (count int64, err error)
}

type resourceWebMenu struct {
}

func (*resourceWebMenu) Insert(ctx context.Context, resourceWebMenu *model.ResourceWebMenu) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resourceWebMenu).Error
	return
}
func (*resourceWebMenu) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.ResourceWebMenu{}, id).Error
	return
}

func (*resourceWebMenu) Select(ctx context.Context, id model.ID) (resourceWebMenu *model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	resourceWebMenu = new(model.ResourceWebMenu)
	err = gdb.Model(resourceWebMenu).Where("id = ?", id).Take(resourceWebMenu).Error
	if err != nil {
		resourceWebMenu = nil
	}
	return
}

func (*resourceWebMenu) SelectByParentID(ctx context.Context, parentID model.ID) (resourceWebMenus []*model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.ResourceWebMenu{}).Where("parent_id = ?", parentID).Find(&resourceWebMenus).Error
	return
}

func (*resourceWebMenu) Update(ctx context.Context, resourceWebMenu *model.ResourceWebMenu) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resourceWebMenu).Save(resourceWebMenu).Error
	return
}

func (*resourceWebMenu) ListByResourceServerID(ctx context.Context, resourceID model.ID, limit int) (list []*model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceWebMenu{}).Where("resource_server_id = ?", resourceID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (*resourceWebMenu) ListRootByResourceServerID(ctx context.Context, resourceServerID model.ID) (list []*model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceWebMenu{}).Where("(ISNULL(parent_id) OR parent_id = '') and resource_server_id = ?", resourceServerID)
	err = exp.Find(&list).Error
	return
}

func (*resourceWebMenu) ListByResourceServerIDAndParentID(ctx context.Context, resourceID, parentID model.ID, limit int) (list []*model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.ResourceWebMenu{}).Where("resource_server_id = ? and parent_id = ?", resourceID, parentID)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (*resourceWebMenu) ListInIDs(ctx context.Context, ids ...model.ID) (list []*model.ResourceWebMenu, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ResourceWebMenu{}).Where("id in ?", ids).Find(&list).Error
	return
}

func (*resourceWebMenu) DeleteInIDs(ctx context.Context, ids ...model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in ?", ids).Delete(model.ResourceWebMenu{}).Error
	return
}

func (r *resourceWebMenu) ListPaged(ctx context.Context, start, limit int) (list []*model.ResourceWebMenu, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceWebMenu{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceWebMenu) ListPagedByResourceServerID(ctx context.Context, start, limit int, resourceServerID model.ID) (list []*model.ResourceWebMenu, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.ResourceWebMenu{}).Where("resource_server_id = ?", resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *resourceWebMenu) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var count int64
	count, err = r.count(ctx, query, args...)
	if count > 0 {
		exist = true
	}
	return
}

func (r *resourceWebMenu) CountByParentID(ctx context.Context, parentID model.ID) (count int64, err error) {
	return r.count(ctx, "parent_id = ?", parentID)
}

func (r *resourceWebMenu) count(ctx context.Context, query interface{}, args ...interface{}) (count int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.ResourceWebMenu{}).Where(query, args...).Count(&count).Error
	return
}
