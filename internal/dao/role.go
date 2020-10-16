package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// Roleer ...
type Roleer interface {
	Insert(ctx context.Context, m *model.Role) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	SelectByCode(ctx context.Context, code model.Code) (m *model.Role, err error)
	SelectByRoot(ctx context.Context) (results []*model.Role, err error)
	SelectByParentCode(ctx context.Context, parentCode model.Code) (results []*model.Role, err error)
	Update(ctx context.Context, m *model.Role) (err error)
}

type role struct {
}

func (r *role) Insert(ctx context.Context, m *model.Role) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}

func (r *role) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Role{}, id).Error
	return
}

func (r *role) SelectByCode(ctx context.Context, code model.Code) (m *model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Role
	err = gdb.Where("code = ?", code).First(&dbResult).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (r *role) SelectByRoot(ctx context.Context) (results []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("ISNULL(parent_code) OR parent_code = ''").Find(&results).Error
	return
}

func (r *role) SelectByParentCode(ctx context.Context, parentCode model.Code) (results []*model.Role, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("parent_code = ?", parentCode).Find(&results).Error
	return
}

func (r *role) Update(ctx context.Context, m *model.Role) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	return
}
