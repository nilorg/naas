package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2Clienter oauth2 client 接口
type OAuth2Clienter interface {
	Insert(ctx context.Context, mc *model.OAuth2Client) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	DeleteInIDs(ctx context.Context, ids []uint64) (err error)
	Update(ctx context.Context, mc *model.OAuth2Client) (err error)
	UpdateRedirectURI(ctx context.Context, id uint64, redirectURI string) (err error)
	SelectByID(ctx context.Context, clientID uint64) (mc *model.OAuth2Client, err error)
	ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total uint64, err error)
}

type oauth2Client struct {
}

func (*oauth2Client) Insert(ctx context.Context, mc *model.OAuth2Client) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	return
}

func (*oauth2Client) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(model.OAuth2Client{}, id).Error
	return
}

func (*oauth2Client) DeleteInIDs(ctx context.Context, ids []uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in (?)", ids).Delete(model.OAuth2Client{}).Error
	return
}

func (*oauth2Client) Update(ctx context.Context, mc *model.OAuth2Client) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Update(mc).Error
	return
}

func (*oauth2Client) UpdateRedirectURI(ctx context.Context, id uint64, redirectURI string) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Client{}).Where("client_id = ?", id).UpdateColumn("redirect_uri", redirectURI).Error
	return
}

func (*oauth2Client) SelectByID(ctx context.Context, clientID uint64) (mc *model.OAuth2Client, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2Client
	err = gdb.Where("client_id = ?", clientID).First(&dbResult).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}

func (*oauth2Client) ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total uint64, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.OAuth2Client{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&clientList).Error
	return
}
