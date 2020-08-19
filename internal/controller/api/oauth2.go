package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type oauth2 struct {
}

// ScopeQueryChildren 查询方法
// @Tags 		OAuth2
// @Summary		查询scope
// @Description	paged:查询翻页列表
// @Description	all:查询所有
// @Accept  json
// @Produce	json
// @Param q query string true "查询参数" Enums(paged,all)
// @Param	current		query	int	false	"当前页"
// @Param	pageSize	query	int	false	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Success 200	{object}	Result{data=[]model.OAuth2Scope}
// @Router /oauth2/scopes [GET]
// @Security OAuth2AccessCode
func (o *oauth2) ScopeQueryChildren() gin.HandlerFunc {
	return QueryChildren(map[string]gin.HandlerFunc{
		"paged": o.scopeListByPaged,
		"all":   o.scopeAll,
	})
}

func (*oauth2) scopeAll(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2Scope
		err    error
	)
	scopes, err = service.OAuth2.AllScope(contexts.WithGinContext(ctx))
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, scopes)
}

func (*oauth2) scopeListByPaged(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2Scope
		err    error
	)
	pagination := model.NewPagination(ctx)
	scopes, pagination.Total, err = service.OAuth2.ScopeListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, scopes))
}

// GetScopeOne 获取一个scope
// @Tags 		OAuth2
// @Summary		获取一个scope
// @Accept  json
// @Produce	json
// @Param 	scop_code	path	string	true	"scop code"
// @Success 200	{object}	Result
// @Router /oauth2/scopes/{scop_code} [GET]
// @Security OAuth2AccessCode
func (*oauth2) GetScopeOne(ctx *gin.Context) {
	var (
		scope *model.OAuth2Scope
		err   error
	)
	scopCode := ctx.Param("scop_code")
	scope, err = service.OAuth2.GetScopeOne(contexts.WithGinContext(ctx), model.Code(scopCode))
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, scope)
}

// UpdateScope 修改scope
// @Tags 		OAuth2
// @Summary		修改scope
// @Accept  json
// @Produce	json
// @Param 	scop_code	path	string	true	"scop code"
// @Param 	body	body	service.OAuth2UpdateScopeModel	true	"body"
// @Success 200	{object}	Result
// @Router /oauth2/scopes/{scop_code} [PUT]
// @Security OAuth2AccessCode
func (*oauth2) UpdateScope(ctx *gin.Context) {
	var (
		scope service.OAuth2UpdateScopeModel
		err   error
	)
	scopCode := ctx.Param("scop_code")
	err = ctx.ShouldBindJSON(&scope)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.OAuth2.UpdateScope(contexts.WithGinContext(ctx), model.Code(scopCode), &scope)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// CreateScope 创建scope
// @Tags 		OAuth2
// @Summary		创建scope
// @Accept  json
// @Produce	json
// @Param 	body	body	service.OAuth2CreateScopeModel	true	"body"
// @Success 200	{object}	Result
// @Router /oauth2/scopes [POST]
// @Security OAuth2AccessCode
func (*oauth2) CreateScope(ctx *gin.Context) {
	var (
		scope service.OAuth2CreateScopeModel
		err   error
	)
	err = ctx.ShouldBindJSON(&scope)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.OAuth2.CreateScope(contexts.WithGinContext(ctx), &scope)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// GetClientScopes 查询客户端scope
// @Tags 		OAuth2
// @Summary		scope
// @Description	查询客户端scope
// @Accept  json
// @Produce	json
// @Param	client_id		query	string	true	"客户端ID"
// @Success 200	{object}	Result{data=model.OAuth2ClientScope}
// @Router /oauth2/clients/{client_id}/scopes [GET]
// @Security OAuth2AccessCode
func (*oauth2) GetClientScopes(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2ClientScope
		err    error
	)
	clientID := ctx.Param("client_id")
	scopes, err = service.OAuth2.GetClientAllScope(contexts.WithGinContext(ctx), model.ConvertStringToID(clientID))
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, scopes)
}

// CreateClient 创建客户端
// @Tags 		OAuth2
// @Summary		client
// @Description	创建客户端
// @Accept  json
// @Produce	json
// @Param 	body	body	service.OAuth2ClientEditModel	true	"body"
// @Success 200	{object}	Result
// @Router /oauth2/clients [POST]
// @Security OAuth2AccessCode
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
	err = service.OAuth2.CreateClient(contexts.WithGinContext(ctx), &create)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// GetClient 获取一个客户端
// @Tags 		OAuth2
// @Summary		获取一个客户端
// @Accept  json
// @Produce	json
// @Param 	client_id	path	string	true	"client id"
// @Success 200	{object}	Result{data=service.OAuth2ClientDetailInfo}
// @Router /oauth2/clients/{client_id} [GET]
// @Security OAuth2AccessCode
func (*oauth2) GetClient(ctx *gin.Context) {
	var (
		clientDetailInfo *service.OAuth2ClientDetailInfo
		err              error
	)
	clientID := model.ConvertStringToID(ctx.Param("client_id"))
	clientDetailInfo, err = service.OAuth2.GetClientDetailInfo(contexts.WithGinContext(ctx), clientID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, clientDetailInfo)
}

// UpdateClient 修改一个客户端
// @Tags 		OAuth2
// @Summary		client
// @Description	根据客户端ID,修改客户端信息
// @Accept  json
// @Produce	json
// @Param 	client_id	path	string	true	"client id"
// @Param 	body	body	service.OAuth2ClientEditModel	true	"客户端信息"
// @Success 200	{object}	Result
// @Router /oauth2/clients/{client_id} [PUT]
// @Security OAuth2AccessCode
func (*oauth2) UpdateClient(ctx *gin.Context) {
	var (
		update service.OAuth2ClientEditModel
		err    error
	)
	clientID := model.ConvertStringToID(ctx.Param("client_id"))
	err = ctx.ShouldBindJSON(&update)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.OAuth2.UpdateClient(contexts.WithGinContext(ctx), clientID, &update)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ClientListByPaged 查询客户端翻页
// @Tags 		OAuth2
// @Summary		client
// @Description	查询客户端翻页
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData{list=service.ResultClientInfo}}
// @Router /oauth2/clients [GET]
// @Security OAuth2AccessCode
func (*oauth2) ClientListByPaged(ctx *gin.Context) {
	var (
		result []*service.ResultClientInfo
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.OAuth2.ClientListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
