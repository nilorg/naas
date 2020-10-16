package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/sdk/convert"
	"gorm.io/gorm"
)

type oauth2 struct {
}

// OAuth2ClientEditModel ...
type OAuth2ClientEditModel struct {
	RedirectURI string `json:"redirect_uri"`
	Name        string `json:"name"`
	Website     string `json:"website"`
	Profile     string `json:"profile"`
	Description string `json:"description"`
}

// CreateClient 创建客户端
func (*oauth2) CreateClient(ctx context.Context, create *OAuth2ClientEditModel) (err error) {
	tran := store.DB.Begin()
	ctx = store.NewDBContext(ctx, tran)
	client := &model.OAuth2Client{
		ClientSecret: uuid.New().String(),
		RedirectURI:  create.RedirectURI,
	}
	err = dao.OAuth2Client.Insert(ctx, client)
	if err != nil {
		tran.Rollback()
		return
	}
	clientInfo := &model.OAuth2ClientInfo{
		ClientID:    client.ClientID,
		Name:        create.Name,
		Website:     create.Website,
		Profile:     create.Profile,
		Description: create.Description,
	}
	if clientInfo.Profile == "" {
		clientInfo.Profile, err = createPicture("oauth2_client", convert.ToString(clientInfo.ClientID))
		if err != nil {
			tran.Rollback()
			return
		}
	}
	err = dao.OAuth2ClientInfo.Insert(ctx, clientInfo)
	if err != nil {
		tran.Rollback()
		return
	}
	err = tran.Commit().Error
	return
}

// UpdateClient 修改客户端
func (*oauth2) UpdateClient(ctx context.Context, id model.ID, update *OAuth2ClientEditModel) (err error) {
	tran := store.DB.Begin()
	ctx = store.NewDBContext(ctx, tran)
	var (
		oauth2Client     *model.OAuth2Client
		oauth2ClientInfo *model.OAuth2ClientInfo
	)
	oauth2Client, err = dao.OAuth2Client.SelectByID(ctx, id)
	if err != nil {
		tran.Rollback()
		return
	}
	if oauth2Client.RedirectURI != update.RedirectURI {
		err = dao.OAuth2Client.UpdateRedirectURI(ctx, id, update.RedirectURI)
		if err != nil {
			tran.Rollback()
			return
		}
		oauth2Client.RedirectURI = update.RedirectURI
	}

	oauth2ClientInfo, err = dao.OAuth2ClientInfo.SelectByClientID(ctx, id)
	if err != nil {
		tran.Rollback()
		return
	}
	oauth2ClientInfo.Name = update.Name
	oauth2ClientInfo.Profile = update.Profile
	oauth2ClientInfo.Description = update.Description
	oauth2ClientInfo.Website = update.Website

	err = dao.OAuth2ClientInfo.Update(ctx, oauth2ClientInfo)
	if err != nil {
		tran.Rollback()
		return
	}
	err = tran.Commit().Error
	return
}

// GetClient get oauth2 client.
func (o *oauth2) GetClient(ctx context.Context, id model.ID) (client *model.OAuth2Client, err error) {
	client, err = dao.OAuth2Client.SelectByID(ctx, id)
	return
}

// GetClient get oauth2 client info.
func (o *oauth2) GetClientInfo(ctx context.Context, id model.ID) (client *model.OAuth2ClientInfo, err error) {
	client, err = dao.OAuth2ClientInfo.SelectByClientID(ctx, id)
	return
}

// OAuth2ClientDetailInfo ...
type OAuth2ClientDetailInfo struct {
	ClientID    model.ID `json:"client_id"`
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Profile     string   `json:"profile"`
	Description string   `json:"description"`
	RedirectURI string   `json:"redirect_uri"`
}

// GetClientDetailInfo get oauth2 client info.
func (o *oauth2) GetClientDetailInfo(ctx context.Context, id model.ID) (clientDetail *OAuth2ClientDetailInfo, err error) {
	var (
		client     *model.OAuth2Client
		clientInfo *model.OAuth2ClientInfo
	)
	client, err = dao.OAuth2Client.SelectByID(ctx, id)
	if err != nil {
		return
	}
	clientInfo, err = dao.OAuth2ClientInfo.SelectByClientID(ctx, id)
	if err != nil {
		return
	}
	clientDetail = &OAuth2ClientDetailInfo{
		ClientID:    client.ClientID,
		Name:        clientInfo.Name,
		Website:     clientInfo.Website,
		Profile:     clientInfo.Profile,
		Description: clientInfo.Description,
		RedirectURI: client.RedirectURI,
	}
	return
}

// ResultClientInfo ...
type ResultClientInfo struct {
	ClientID    model.ID `json:"client_id"`
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Profile     string   `json:"profile"`
	Description string   `json:"description"`
	RedirectURI string   `json:"redirect_uri"`
}

// ClientListPaged ...
func (o *oauth2) ClientListPaged(ctx context.Context, start, limit int) (result []*ResultClientInfo, total int64, err error) {
	var (
		clientList []*model.OAuth2Client
	)
	clientList, total, err = dao.OAuth2Client.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, client := range clientList {
		clientInfo, clientInfoErr := o.GetClientInfo(ctx, client.ClientID)
		resultInfo := &ResultClientInfo{}
		if clientInfoErr == nil {
			resultInfo.ClientID = clientInfo.ClientID
			resultInfo.Name = clientInfo.Name
			resultInfo.Website = clientInfo.Website
			resultInfo.Profile = clientInfo.Profile
			resultInfo.Description = clientInfo.Description
			resultInfo.RedirectURI = client.RedirectURI
		}
		result = append(result, resultInfo)
	}
	return
}

// GetClientScope 获取客户端的范围
func (o *oauth2) GetClientAllScope(ctx context.Context, clientID model.ID) (scopes []*model.OAuth2ClientScope, err error) {
	scopes, err = dao.OAuth2ClientScope.SelectByOAuth2ClientID(ctx, clientID)
	return
}

// OAuth2ClientScopeInfo ...
type OAuth2ClientScopeInfo struct {
	Code        model.Code `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
}

func (o *oauth2) GetClientAllScopeInfo(ctx context.Context, clientID model.ID) (scopeInfos []*OAuth2ClientScopeInfo, err error) {
	var clientScopes []*model.OAuth2ClientScope
	clientScopes, err = dao.OAuth2ClientScope.SelectByOAuth2ClientID(ctx, clientID)
	if err != nil {
		return
	}
	var scope *model.OAuth2Scope
	for _, clientScope := range clientScopes {
		scope, err = dao.OAuth2Scope.Select(ctx, clientScope.ScopeCode)
		if err != nil {
			return
		}
		scopeInfos = append(scopeInfos, &OAuth2ClientScopeInfo{
			Code:        scope.Code,
			Name:        scope.Name,
			Description: scope.Description,
			Type:        scope.Type,
		})
	}
	return
}

// GetClientAllScopeCode 获取客户端的范围code
func (o *oauth2) GetClientAllScopeCode(ctx context.Context, clientID model.ID) (scopes []model.Code, err error) {
	var scopeArray []*model.OAuth2ClientScope
	scopeArray, err = dao.OAuth2ClientScope.SelectByOAuth2ClientID(ctx, clientID)
	if err != nil {
		return
	}
	for _, scope := range scopeArray {
		scopes = append(scopes, scope.ScopeCode)
	}
	return
}

// AllScope 获取所有的范围
func (o *oauth2) AllScope(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	scopes, err = dao.OAuth2Scope.SelectAll(ctx)
	return
}

// GetScopeOne 根据code获取scope
func (o *oauth2) GetScopeOne(ctx context.Context, code model.Code) (scope *model.OAuth2Scope, err error) {
	return dao.OAuth2Scope.Select(ctx, code)
}

// OAuth2EditScopeModel ...
type OAuth2EditScopeModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// EditScope 编辑/创建scope
func (o *oauth2) EditScope(ctx context.Context, code model.Code, edit *OAuth2EditScopeModel) (err error) {
	var (
		scope  *model.OAuth2Scope
		insert bool
	)
	scope, err = dao.OAuth2Scope.Select(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			insert = true
			scope = new(model.OAuth2Scope)
			scope.Code = code
			err = nil
		} else {
			return
		}
	}
	scope.Name = edit.Name
	scope.Description = edit.Description
	scope.Type = edit.Type
	if insert {
		err = dao.OAuth2Scope.Insert(ctx, scope)
	} else {
		err = dao.OAuth2Scope.Update(ctx, scope)
	}
	return
}

func (o *oauth2) ScopeListPaged(ctx context.Context, start, limit int) (scopes []*model.OAuth2Scope, total int64, err error) {
	scopes, total, err = dao.OAuth2Scope.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	return
}

func (o *oauth2) AllScopeCode(ctx context.Context) (scopeCodes []model.Code, err error) {
	var scopes []*model.OAuth2Scope
	scopes, err = dao.OAuth2Scope.SelectAll(ctx)
	if err != nil {
		return
	}
	for _, scope := range scopes {
		scopeCodes = append(scopeCodes, scope.Code)
	}
	return
}
