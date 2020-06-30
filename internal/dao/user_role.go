package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// UserRoleer ...
type UserRoleer interface {
	Insert(ctx context.Context, m *model.UserRole) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (m *model.UserRole, err error)
	Update(ctx context.Context, m *model.UserRole) (err error)
	SelectAllByUserID(ctx context.Context, userID uint64) (m []*model.UserRole, err error)
}

type userRole struct {
}

func (u *userRole) Insert(ctx context.Context, m *model.UserRole) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	return
}
func (u *userRole) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserRole{}, id).Error
	return
}
func (u *userRole) Select(ctx context.Context, id uint64) (m *model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.UserRole
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}
func (u *userRole) Update(ctx context.Context, m *model.UserRole) (err error) {
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

func (u *userRole) SelectAllByUserID(ctx context.Context, userID uint64) (roles []*model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.UserRole{}).Where("user_id = ?", userID).Find(&roles).Error
	return
}
