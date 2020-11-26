package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/strings"
)

type resource struct {
}

// GetServerOne 获取一个资源
// @Tags 		Resource（资源）
// @Summary		获取一个资源
// @Description	根据资源ID,获取一个资源
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource id"
// @Success 200	{object}	Result
// @Router /resource/servers/{resource_server_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetServerOne(ctx *gin.Context) {
	var (
		r   *model.Resource
		err error
	)
	id := model.ConvertStringToID(ctx.Param("resource_server_id"))
	r, err = service.Resource.GetServer(contexts.WithGinContext(ctx), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, r)
}

func (*resource) DeleteServer(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Param("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteByIDs(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListServerByPaged 查询资源
// @Tags 		Resource（资源）
// @Summary		查询资源
// @Description	查询资源翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/servers [GET]
// @Security OAuth2AccessCode
func (*resource) ListServerByPaged(ctx *gin.Context) {
	var (
		result []*model.ResultResourceServer
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListServerPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// CreateServer 创建资源
// @Tags 		Resource（资源）
// @Summary		创建资源
// @Description	创建资源
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceEditModel	true	"body"
// @Success 200	{object}	Result
// @Router /resource/servers [POST]
// @Security OAuth2AccessCode
func (*resource) CreateServer(ctx *gin.Context) {
	var (
		m   service.ResourceEditModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.CreateServer(contexts.WithGinContext(ctx), &m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateServer 修改一个资源
// @Tags 		Resource（资源）
// @Summary		修改一个资源
// @Description	根据资源ID,修改一个资源
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource id"
// @Param 	body	body	service.ResourceEditModel	true	"资源需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/servers/{resource_server_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateServer(ctx *gin.Context) {
	var (
		org service.ResourceEditModel
		err error
	)
	resID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	err = ctx.ShouldBindJSON(&org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateServer(contexts.WithGinContext(ctx), resID, &org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ========================

// AddWebRoute 添加web路由
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		添加web路由
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceWebRouteEdit	true	"body"
// @Success 200	{object}	Result
// @Router /resource/web_routes [POST]
// @Security OAuth2AccessCode
func (*resource) AddWebRoute(ctx *gin.Context) {
	var (
		req service.ResourceWebRouteEdit
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.AddWebRoute(contexts.WithGinContext(ctx), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateWebRoute 修改资源web路由
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		修改一个资源web路由
// @Description	根据资源web路由ID,修改一个资源web路由
// @Accept  json
// @Produce	json
// @Param 	resource_web_route_id	path	string	true	"resource web route id"
// @Param 	body	body	service.ResourceWebRouteEdit	true	"Web路由需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_route_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateWebRoute(ctx *gin.Context) {
	var (
		route service.ResourceWebRouteEdit
		err   error
	)
	id := model.ConvertStringToID(ctx.Param("resource_web_route_id"))
	err = ctx.ShouldBindJSON(&route)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateWebRoute(contexts.WithGinContext(ctx), id, &route)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateWebRoute 删除资源web路由
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		删除一个资源web路由
// @Description	根据资源web路由ID,删除一个资源web路由
// @Accept  json
// @Produce	json
// @Param 	resource_web_route_id	path	string	true	"resource web route id"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_route_id} [DELETE]
// @Security OAuth2AccessCode
func (*resource) DeleteWebRoute(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Query("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteWebRoute(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListWebRoutePaged 查询资源服务器WebRule
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		查询资源服务器WebRule
// @Description	查询资源服务器WebRule翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/web_routes [GET]
// @Security OAuth2AccessCode
func (*resource) ListWebRoutePaged(ctx *gin.Context) {
	var (
		result []*model.ResultResourceWebRoute
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListWebRoutePaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// GetWebRouteOne 获取一个Web路由
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		获取一个Web路由
// @Description	根据路由ID,获取一个Web路由
// @Accept  json
// @Produce	json
// @Param 	resource_web_route_id	path	string	true	"resource web route id"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_route_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetWebRouteOne(ctx *gin.Context) {
	var (
		menu *model.ResourceWebRoute
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_web_route_id"))
	menu, err = service.Resource.GetResourceWebRoute(contexts.WithGinContext(ctx), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, menu)
}

// ========================

// AddWebMenu 添加web菜单
// @Tags 		ResourceWebRoute（资源服务器Web路由）
// @Summary		添加web菜单
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceWebMenuEdit	true	"body"
// @Success 200	{object}	Result
// @Router /resource/web_menus [POST]
// @Security OAuth2AccessCode
func (*resource) AddWebMenu(ctx *gin.Context) {
	var (
		req service.ResourceWebMenuEdit
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.AddWebMenu(contexts.WithGinContext(ctx), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateWebMenu 修改资源web菜单
// @Tags 		ResourceWebMenu（资源服务器Web菜单）
// @Summary		修改一个资源web菜单
// @Description	根据资源web菜单ID,修改一个资源web菜单
// @Produce	json
// @Param 	resource_web_menu_id	path	string	true	"resource web menu id"
// @Param 	body	body	service.ResourceWebMenuEdit	true	"Web路由需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_menu_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateWebMenu(ctx *gin.Context) {
	var (
		menu service.ResourceWebMenuEdit
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_web_menu_id"))
	err = ctx.ShouldBindJSON(&menu)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateWebMenu(contexts.WithGinContext(ctx), id, &menu)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateWebMenu 删除资源web菜单
// @Tags 		ResourceWebMenu（资源服务器Web菜单）
// @Summary		删除一个资源web菜单
// @Description	根据资源web菜单ID,修改一个资源web菜单
// @Produce	json
// @Param 	resource_web_menu_id	path	string	true	"resource web menu id"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_menu_id} [DELETE]
// @Security OAuth2AccessCode
func (*resource) DeleteWebMenu(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Query("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteWebMenu(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListWebMenuPaged 查询资源服务器WebMenu
// @Tags 		ResourceWebMenu（资源服务器Web菜单）
// @Summary		查询资源服务器Web菜单
// @Description	查询资源服务器Web菜单翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/web_menus [GET]
// @Security OAuth2AccessCode
func (*resource) ListWebMenuPaged(ctx *gin.Context) {
	var (
		result []*model.ResultResourceWebMenu
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListWebMenuPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// GetWebMenuOne 获取一个Web菜单
// @Tags 		ResourceWebMenu（资源服务器Web菜单）
// @Summary		获取一个Web菜单
// @Description	根据菜单ID,获取一个Web菜单
// @Accept  json
// @Produce	json
// @Param 	resource_web_menu_id	path	string	true	"resource web menu id"
// @Success 200	{object}	Result
// @Router /resource/web_routes/{resource_web_menu_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetWebMenuOne(ctx *gin.Context) {
	var (
		menu *model.ResourceWebMenu
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_web_menu_id"))
	menu, err = service.Resource.GetResourceWebMenu(contexts.WithGinContext(ctx), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, menu)
}

// ========================

// AddAction 添加动作
// @Tags 		ResourceAction（资源动作）
// @Summary		添加动作
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceActionEdit	true	"body"
// @Success 200	{object}	Result
// @Router /resource/actions [POST]
// @Security OAuth2AccessCode
func (*resource) AddAction(ctx *gin.Context) {
	var (
		req service.ResourceActionEdit
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.AddAction(contexts.WithGinContext(ctx), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateAction 修改动作
// @Tags 		ResourceAction（资源动作）
// @Summary		修改一个资源动作
// @Description	根据资源动作ID,修改一个资源动作
// @Produce	json
// @Param 	resource_action_id	path	string	true	"resource action id"
// @Param 	body	body	service.ResourceActionEdit	true	"资源动作需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/actions/{resource_action_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateAction(ctx *gin.Context) {
	var (
		action service.ResourceActionEdit
		err    error
	)
	id := model.ConvertStringToID(ctx.Param("resource_action_id"))
	err = ctx.ShouldBindJSON(&action)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateAction(contexts.WithGinContext(ctx), id, &action)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// DeleteAction 删除资源动作
// @Tags 		ResourceAction（资源动作）
// @Summary		删除一个资源动作
// @Description	根据资源动作ID,修改一个资源动作
// @Produce	json
// @Param 	resource_action_id	path	string	true	"resource action id"
// @Success 200	{object}	Result
// @Router /resource/actions/{resource_action_id} [DELETE]
// @Security OAuth2AccessCode
func (*resource) DeleteAction(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Query("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteAction(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListActionPaged 查询资源服务器动作
// @Tags 		ResourceAction（资源动作）
// @Summary		查询资源服务器动作
// @Description	查询资源服务器动作翻页数据
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/actions [GET]
// @Security OAuth2AccessCode
func (*resource) ListActionPaged(ctx *gin.Context) {
	var (
		list []*model.ResultResourceAction
		err  error
	)
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Resource.ListActionPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// GetActionOne 获取一个动作
// @Tags 		ResourceAction（资源动作）
// @Summary		获取一个动作
// @Description	根据菜单ID,获取一个动作
// @Accept  json
// @Produce	json
// @Param 	resource_action_id	path	string	true	"resource action id"
// @Success 200	{object}	Result
// @Router /resource/actions/{resource_action_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetActionOne(ctx *gin.Context) {
	var (
		action *model.ResourceAction
		err    error
	)
	id := model.ConvertStringToID(ctx.Param("resource_action_id"))
	action, err = service.Resource.GetResourceAction(contexts.WithGinContext(ctx), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, action)
}
