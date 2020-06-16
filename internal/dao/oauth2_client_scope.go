package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2ClientScoper oauth2 client 范围 接口
type OAuth2ClientScoper interface {
	SelectByOAuth2ClientID(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error)
	Insert(ctx context.Context, mc *model.OAuth2ClientScope) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Update(ctx context.Context, mc *model.OAuth2ClientScope) (err error)
}

type oauth2ClientScope struct {
}

func (*oauth2ClientScope) SelectByOAuth2ClientID(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("client_id = ?", clientID).Find(&scopes).Error
	return
}

func (*oauth2ClientScope) Insert(ctx context.Context, mc *model.OAuth2ClientScope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	return
}

func (*oauth2ClientScope) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(model.OAuth2ClientScope{}, id).Error
	return
}

func (*oauth2ClientScope) Update(ctx context.Context, mc *model.OAuth2ClientScope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Update(mc).Error
	return
}
