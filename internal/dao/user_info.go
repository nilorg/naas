package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// UserInfoer ...
type UserInfoer interface {
	SelectByUserID(ctx context.Context, userID string) (mu *model.UserInfo, err error)
	Insert(ctx context.Context, mu *model.UserInfo) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (mu *model.UserInfo, err error)
	Update(ctx context.Context, mu *model.UserInfo) (err error)
}

type userInfo struct {
}

func (*userInfo) SelectByUserID(ctx context.Context, userID string) (mu *model.UserInfo, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.UserInfo
	err = gdb.Where("user_id = ?", userID).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (*userInfo) Insert(ctx context.Context, mu *model.UserInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mu).Error
	return
}

func (*userInfo) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Unscoped().Delete(&model.UserInfo{}, id).Error
	return
}

func (*userInfo) Select(ctx context.Context, id uint64) (mu *model.UserInfo, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.UserInfo
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (*userInfo) Update(ctx context.Context, mu *model.UserInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Update(mu).Error
	if err != nil {
		return
	}
	return
}
