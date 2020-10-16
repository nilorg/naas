package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// Adminer ...
type Adminer interface {
	SelectByUsername(ctx context.Context, username string) (ma *model.Admin, err error)
	Insert(ctx context.Context, ma *model.Admin) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (ma *model.Admin, err error)
	Update(ctx context.Context, ma *model.Admin) (err error)
}

type admin struct {
}

func (*admin) SelectByUsername(ctx context.Context, username string) (ma *model.Admin, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.Admin
	err = gdb.Where("username = ?", username).First(&dbResult).Error
	if err != nil {
		return
	}
	ma = &dbResult
	return
}

func (*admin) Insert(ctx context.Context, ma *model.Admin) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(ma).Error
	return
}

func (*admin) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Admin{}, id).Error
	return
}

func (*admin) Select(ctx context.Context, id model.ID) (ma *model.Admin, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(ma).Where("id = ?", id).Scan(ma).Error
	if err != nil {
		ma = nil
		return
	}
	return
}

func (*admin) Update(ctx context.Context, ma *model.Admin) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(ma).Save(ma).Error
	if err != nil {
		return
	}
	return
}
