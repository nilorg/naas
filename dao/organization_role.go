package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/model"
	"github.com/nilorg/pkg/db"
)

// OrganizationRoleer ...
type OrganizationRoleer interface {
	Insert(ctx context.Context, m *model.OrganizationRole) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (m *model.OrganizationRole, err error)
	Update(ctx context.Context, m *model.OrganizationRole) (err error)
}

type organizationRole struct {
}

func (o *organizationRole) Insert(ctx context.Context, m *model.OrganizationRole) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}
func (o *organizationRole) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Unscoped().Delete(&model.OrganizationRole{}, id).Error
	return
}
func (o *organizationRole) Select(ctx context.Context, id uint64) (m *model.OrganizationRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OrganizationRole
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}
func (o *organizationRole) Update(ctx context.Context, m *model.OrganizationRole) (err error) {
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
