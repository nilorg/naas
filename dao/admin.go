package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/model"
	"github.com/nilorg/pkg/db"
)

// Adminer ...
type Adminer interface {
	SelectByUsername(ctx context.Context, username string) (ma *model.Admin, err error)
}

type admin struct {
}

func (*admin) SelectByUsername(ctx context.Context, username string) (ma *model.Admin, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
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
