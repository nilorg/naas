package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type oauth2 struct {
}

func (*oauth2) GetScopes(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2Scope
		err    error
	)
	scopes, err = service.OAuth2.AllScope()
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, scopes)
}

func (*oauth2) GetClientScopes(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2ClientScope
		err    error
	)
	clientID := ctx.Param("client_id")
	scopes, err = service.OAuth2.GetClientAllScope(convert.ToUint64(clientID))
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, scopes)
}

func (*oauth2) CreateClient(ctx *gin.Context) {
	var (
		create service.OAuth2ClientEditModel
		err    error
	)
	err = ctx.ShouldBindJSON(&create)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.OAuth2.CreateClient(&create)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*oauth2) UpdateClient(ctx *gin.Context) {
	var (
		update service.OAuth2ClientEditModel
		err    error
	)
	clientID := convert.ToUint64(ctx.Param("client_id"))
	err = ctx.ShouldBindJSON(&update)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.OAuth2.UpdateClient(clientID, &update)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*oauth2) ClientListByPaged(ctx *gin.Context) {
	var (
		result []*service.ResultClientInfo
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.OAuth2.ClientListPaged(pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
