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
	ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total uint64, err error)
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
