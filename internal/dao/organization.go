package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// Organizationer ...
type Organizationer interface {
	Insert(ctx context.Context, m *model.Organization) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (m *model.Organization, err error)
	Update(ctx context.Context, m *model.Organization) (err error)
}

type organization struct {
}

func (o *organization) Insert(ctx context.Context, m *model.Organization) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}
func (o *organization) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Unscoped().Delete(&model.Organization{}, id).Error
	return
}
func (o *organization) Select(ctx context.Context, id uint64) (m *model.Organization, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Organization
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}
func (o *organization) Update(ctx context.Context, m *model.Organization) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Update(m).Error
	if err != nil {
		return
	}
	return
}
