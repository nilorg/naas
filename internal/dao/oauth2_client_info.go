package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2ClientInfoer oauth2 client 接口
type OAuth2ClientInfoer interface {
	SelectByClientID(ctx context.Context, clientID string) (mc *model.OAuth2ClientInfo, err error)
	Insert(ctx context.Context, mc *model.OAuth2ClientInfo) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	DeleteByClientID(ctx context.Context, clientID uint64) (err error)
	DeleteInClientIDs(ctx context.Context, clientIDs []uint64) (err error)
	Update(ctx context.Context, mc *model.OAuth2ClientInfo) (err error)
}

type oauth2ClientInfo struct {
}

func (*oauth2ClientInfo) SelectByClientID(ctx context.Context, clientID string) (mc *model.OAuth2ClientInfo, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2ClientInfo
	err = gdb.Where("client_id = ?", clientID).First(&dbResult).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}

func (*oauth2ClientInfo) Insert(ctx context.Context, mc *model.OAuth2ClientInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	return
}

func (*oauth2ClientInfo) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(model.OAuth2ClientInfo{}, id).Error
	return
}

func (*oauth2ClientInfo) DeleteByClientID(ctx context.Context, clientID uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("client_id = ?", clientID).Delete(model.OAuth2ClientInfo{}).Error
	return
}

func (*oauth2ClientInfo) DeleteInClientIDs(ctx context.Context, clientIDs []uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("client_id in (?)", clientIDs).Delete(model.OAuth2ClientInfo{}).Error
	return
}

func (*oauth2ClientInfo) Update(ctx context.Context, mc *model.OAuth2ClientInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Update(mc).Error
	return
}
