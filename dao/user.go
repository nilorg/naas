package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/model"
	"github.com/nilorg/pkg/db"
)

// Userer ...
type Userer interface {
	SelectByUsername(ctx context.Context, username string) (mu *model.User, err error)
	Insert(ctx context.Context, mu *model.User) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Select(ctx context.Context, id uint64) (mu *model.User, err error)
	Update(ctx context.Context, mu *model.User) (err error)
}

type user struct {
}

func (*user) SelectByUsername(ctx context.Context, username string) (mu *model.User, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.User
	err = gdb.Where("username = ?", username).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (*user) Insert(ctx context.Context, mu *model.User) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mu).Error
	return
}

func (*user) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Unscoped().Delete(&model.User{}, id).Error
	return
}

func (*user) Select(ctx context.Context, id uint64) (mu *model.User, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.User
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (*user) Update(ctx context.Context, mu *model.User) (err error) {
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
