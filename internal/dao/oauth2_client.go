package dao

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/pkg/db"
)

// OAuth2Clienter oauth2 client 接口
type OAuth2Clienter interface {
	SelectByID(ctx context.Context, clientID string) (mc *model.OAuth2Client, err error)
}

type oauth2Client struct {
}

func (*oauth2Client) SelectByID(ctx context.Context, clientID string) (mc *model.OAuth2Client, err error) {
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
