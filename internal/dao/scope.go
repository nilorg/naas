package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// Scoper ...
type Scoper interface {
	Insert(ctx context.Context, m *model.Scope) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (m *model.Scope, err error)
	SelectAll(ctx context.Context) (m []*model.Scope, err error)
	Update(ctx context.Context, m *model.Scope) (err error)
}

type scope struct {
}

func (s *scope) Insert(ctx context.Context, m *model.Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}

func (s *scope) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Role{}, id).Error
	return
}

func (s *scope) SelectAll(ctx context.Context) (m []*model.Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.Scope{}).Find(&m).Error
	return
}

func (s *scope) Select(ctx context.Context, id uint64) (m *model.Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Scope
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (s *scope) Update(ctx context.Context, m *model.Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Update(m).Error
	if err != nil {
		return
	}
	return
}
