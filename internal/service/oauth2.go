package service

import (
	"github.com/jinzhu/gorm"
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

// GetClient get oauth2 client info.
func (o *oauth2) GetClientInfo(id uint64) (client *model.OAuth2ClientInfo, err error) {
	client, err = dao.OAuth2ClientInfo.SelectByClientID(store.NewDBContext(), id)
	return
}

// ResultClientInfo ...
type ResultClientInfo struct {
	*model.OAuth2Client
	*model.OAuth2ClientInfo
}

// ClientListPaged ...
func (o *oauth2) ClientListPaged(start, limit int) (result []*ResultClientInfo, total uint64, err error) {
	var (
		clientList []*model.OAuth2Client
	)
	clientList, total, err = dao.OAuth2Client.ListPaged(store.NewDBContext(), start, limit)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = nil
		}
		return
	}
	for _, client := range clientList {
		clientInfo, clientInfoErr := o.GetClientInfo(client.ClientID)
		resultInfo := &ResultClientInfo{
			OAuth2Client: client,
		}
		if clientInfoErr == nil {
			resultInfo.OAuth2ClientInfo = clientInfo
		}
		result = append(result, resultInfo)
	}
	return
}

// GetClientScope 获取客户端的范围
// TODO: 后期添加缓存
func (o *oauth2) GetClientAllScope(clientID uint64) (scopes []*model.OAuth2ClientScope, err error) {
	scopes, err = dao.OAuth2ClientScope.SelectByOAuth2ClientID(store.NewDBContext(), clientID)
	return
}

// AllScope 获取所有的范围
// TODO: 后期添加缓存
func (o *oauth2) AllScope() (scopes []*model.OAuth2Scope, err error) {
	scopes, err = dao.OAuth2Scope.SelectAll(store.NewDBContext())
	return
}
