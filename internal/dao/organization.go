package dao

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// Organizationer ...
type Organizationer interface {
	Insert(ctx context.Context, m *model.Organization) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids ...model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.Organization, err error)
	SelectByRoot(ctx context.Context) (results []*model.Organization, err error)
	SelectByParentID(ctx context.Context, parentID model.ID) (results []*model.Organization, err error)
	Update(ctx context.Context, m *model.Organization) (err error)
	ListPaged(ctx context.Context, start, limit int) (list []*model.Organization, total int64, err error)
	ListByName(ctx context.Context, name string, limit int) (list []*model.Organization, err error)
	ExistByCode(ctx context.Context, code model.Code) (exist bool, err error)
	ExistByID(ctx context.Context, id model.ID) (exist bool, err error)
}

type organization struct {
}

func (o *organization) Insert(ctx context.Context, m *model.Organization) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}
func (o *organization) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Organization{}, id).Error
	return
}
func (o *organization) Select(ctx context.Context, id model.ID) (m *model.Organization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	m = new(model.Organization)
	err = gdb.Model(m).Where("id = ?", id).Take(m).Error
	if err != nil {
		m = nil
		return
	}
	return
}

func (o *organization) SelectByRoot(ctx context.Context) (results []*model.Organization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("ISNULL(parent_id) OR parent_id = 0").Find(&results).Error
	return
}

func (o *organization) SelectByParentID(ctx context.Context, parentID model.ID) (results []*model.Organization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("parent_id = ?", parentID).Find(&results).Error
	return
}

func (o *organization) Update(ctx context.Context, m *model.Organization) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	if err != nil {
		return
	}
	return
}

func (o *organization) ListPaged(ctx context.Context, start, limit int) (list []*model.Organization, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.Organization{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (o *organization) ListByName(ctx context.Context, name string, limit int) (list []*model.Organization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.Organization{}).Where("name like ?", fmt.Sprintf("%%%s%%", name))
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Offset(0).Limit(limit).Find(&list).Error
	return
}

func (o *organization) DeleteInIDs(ctx context.Context, ids ...model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in ?", ids).Delete(model.Organization{}).Error
	return
}

func (o *organization) ExistByCode(ctx context.Context, code model.Code) (exist bool, err error) {
	return o.exist(ctx, "code = ?", code)
}

func (o *organization) ExistByID(ctx context.Context, id model.ID) (exist bool, err error) {
	return o.exist(ctx, "id = ?", id)
}

func (o *organization) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.Organization{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}
