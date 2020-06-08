package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// UserInfoer ...
type UserInfoer interface {
	SelectByUserID(ctx context.Context, userID uint64) (mu *model.UserInfo, err error)
	Insert(ctx context.Context, mu *model.UserInfo) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	DeleteByUserID(ctx context.Context, userID uint64) (err error)
	DeleteInUserIDs(ctx context.Context, userIDs []uint64) (err error)
	Select(ctx context.Context, id uint64) (mu *model.UserInfo, err error)
	Update(ctx context.Context, mu *model.UserInfo) (err error)
}

type userInfo struct {
}

func (*userInfo) SelectByUserID(ctx context.Context, userID uint64) (mu *model.UserInfo, err error) {
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
	err = gdb.Model(mu).Create(mu).Error
	return
}

func (*userInfo) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserInfo{}, id).Error
	return
}

func (*userInfo) DeleteByUserID(ctx context.Context, userID uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("user_id = ?", userID).Delete(model.UserInfo{}).Error
	return
}

func (*userInfo) DeleteInUserIDs(ctx context.Context, userIDs []uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("user_id in (?)", userIDs).Delete(model.UserInfo{}).Error
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
