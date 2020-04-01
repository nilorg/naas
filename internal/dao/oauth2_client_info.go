package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2ClientInfoer oauth2 client 接口
type OAuth2ClientInfoer interface {
	SelectByID(ctx context.Context, clientID string) (mc *model.OAuth2ClientInfo, err error)
}

type oauth2ClientInfo struct {
}

func (*oauth2ClientInfo) SelectByID(ctx context.Context, clientID string) (mc *model.OAuth2ClientInfo, err error) {
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
