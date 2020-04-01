package service

import (
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
)

type oauth2 struct {
}

// GetClient get oauth2 client.
func (o *oauth2) GetClient(id string) (client *model.OAuth2Client, err error) {
	client, err = dao.OAuth2Client.SelectByID(store.NewDBContext(), id)
	return
}
