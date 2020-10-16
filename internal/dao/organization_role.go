package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
	"gorm.io/gorm"
)

// OrganizationRoleer ...
type OrganizationRoleer interface {
	Insert(ctx context.Context, m *model.OrganizationRole) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.OrganizationRole, err error)
	Update(ctx context.Context, m *model.OrganizationRole) (err error)
}

type organizationRole struct {
}

func (o *organizationRole) Insert(ctx context.Context, m *model.OrganizationRole) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromGormV2Context(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}
func (o *organizationRole) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromGormV2Context(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.OrganizationRole{}, id).Error
	return
}
func (o *organizationRole) Select(ctx context.Context, id model.ID) (m *model.OrganizationRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromGormV2Context(ctx)
	if err != nil {
		return
	}
	m = new(model.OrganizationRole)
	err = gdb.Model(m).Where("id = ?", id).Scan(m).Error
	if err != nil {
		m = nil
		return
	}
	return
}
func (o *organizationRole) Update(ctx context.Context, m *model.OrganizationRole) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromGormV2Context(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	if err != nil {
		return
	}
	return
}
