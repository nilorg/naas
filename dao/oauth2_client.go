package dao

import (
	"github.com/nilorg/naas/model"
	"github.com/nilorg/naas/module/store"
)

type oauth2Client struct {
}

func (*oauth2Client) SelectByID(clientID string) (mc *model.OAuth2Client, err error) {
	var dbResult model.OAuth2Client
	err = store.DB.Where("client_id = ?", clientID).First(&dbResult).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}
