package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type oauth2 struct {
}

// QueryChildren 查询方法
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
func (o *oauth2) QueryChildren() gin.HandlerFunc {
	return QueryChildren(map[string]gin.HandlerFunc{
		"paged": o.listByPaged,
		"all":   o.allScope,
	})
}

func (*oauth2) allScope(ctx *gin.Context) {
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

func (*oauth2) listByPaged(ctx *gin.Context) {
	var (
		scopes []*model.OAuth2Scope
		err    error
	)
	pagination := model.NewPagination(ctx)
	scopes, pagination.Total, err = service.OAuth2.ScopeListPaged(pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, scopes))
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
	scopes, err = service.OAuth2.GetClientAllScope(convert.ToUint64(clientID))
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
	err = service.OAuth2.CreateClient(&create)
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
	clientID := convert.ToUint64(ctx.Param("client_id"))
	clientDetailInfo, err = service.OAuth2.GetClientDetailInfo(clientID)
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
	result, pagination.Total, err = service.OAuth2.ClientListPaged(pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
