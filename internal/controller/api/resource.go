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

// AddRoute 添加路由
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		添加路由
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceRouteEdit	true	"body"
// @Success 200	{object}	Result
// @Router /resource/routes [POST]
// @Security OAuth2AccessCode
func (*resource) AddRoute(ctx *gin.Context) {
	var (
		req service.ResourceRouteEdit
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.AddRoute(contexts.WithGinContext(ctx), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateRoute 修改资源路由
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		修改一个资源路由
// @Description	根据资源路由ID,修改一个资源路由
// @Accept  json
// @Produce	json
// @Param 	resource_route_id	path	string	true	"resource route id"
// @Param 	body	body	service.ResourceRouteEdit	true	"路由需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_route_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateRoute(ctx *gin.Context) {
	var (
		route service.ResourceRouteEdit
		err   error
	)
	id := model.ConvertStringToID(ctx.Param("resource_route_id"))
	err = ctx.ShouldBindJSON(&route)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateRoute(contexts.WithGinContext(ctx), id, &route)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateRoute 删除资源路由
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		删除一个资源路由
// @Description	根据资源路由ID,删除一个资源路由
// @Accept  json
// @Produce	json
// @Param 	resource_route_id	path	string	true	"resource route id"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_route_id} [DELETE]
// @Security OAuth2AccessCode
func (*resource) DeleteRoute(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Query("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteRoute(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListRoutePaged 查询资源服务器Rule
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		查询资源服务器Rule
// @Description	查询资源服务器Rule翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/routes [GET]
// @Security OAuth2AccessCode
func (*resource) ListRoutePaged(ctx *gin.Context) {
	var (
		result []*model.ResultResourceRoute
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListRoutePaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// GetRouteOne 获取一个路由
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		获取一个路由
// @Description	根据路由ID,获取一个路由
// @Accept  json
// @Produce	json
// @Param 	resource_route_id	path	string	true	"resource route id"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_route_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetRouteOne(ctx *gin.Context) {
	var (
		menu *model.ResourceRoute
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_route_id"))
	menu, err = service.Resource.GetResourceRoute(contexts.WithGinContext(ctx), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, menu)
}

// ========================

// AddMenu 添加菜单
// @Tags 		ResourceRoute（资源服务器路由）
// @Summary		添加菜单
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceMenuEdit	true	"body"
// @Success 200	{object}	Result
// @Router /resource/menus [POST]
// @Security OAuth2AccessCode
func (*resource) AddMenu(ctx *gin.Context) {
	var (
		req service.ResourceMenuEdit
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.AddMenu(contexts.WithGinContext(ctx), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UpdateMenu 修改资源菜单
// @Tags 		ResourceMenu（资源服务器菜单）
// @Summary		修改一个资源菜单
// @Description	根据资源菜单ID,修改一个资源菜单
// @Produce	json
// @Param 	resource_menu_id	path	string	true	"resource menu id"
// @Param 	body	body	service.ResourceMenuEdit	true	"路由需要修改的信息"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_menu_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) UpdateMenu(ctx *gin.Context) {
	var (
		menu service.ResourceMenuEdit
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_menu_id"))
	err = ctx.ShouldBindJSON(&menu)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.UpdateMenu(contexts.WithGinContext(ctx), id, &menu)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// DeleteMenu 删除资源菜单
// @Tags 		ResourceMenu（资源服务器菜单）
// @Summary		删除一个资源菜单
// @Description	根据资源菜单ID,修改一个资源菜单
// @Produce	json
// @Param 	resource_menu_id	path	string	true	"resource menu id"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_menu_id} [DELETE]
// @Security OAuth2AccessCode
func (*resource) DeleteMenu(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Query("ids"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Resource.DeleteMenu(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListMenuPaged 查询资源服务器Menu
// @Tags 		ResourceMenu（资源服务器菜单）
// @Summary		查询资源服务器菜单
// @Description	查询资源服务器菜单翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resource/menus [GET]
// @Security OAuth2AccessCode
func (*resource) ListMenuPaged(ctx *gin.Context) {
	var (
		result []*model.ResultResourceMenu
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListMenuPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// GetMenuOne 获取一个菜单
// @Tags 		ResourceMenu（资源服务器菜单）
// @Summary		获取一个菜单
// @Description	根据菜单ID,获取一个菜单
// @Accept  json
// @Produce	json
// @Param 	resource_menu_id	path	string	true	"resource menu id"
// @Success 200	{object}	Result
// @Router /resource/routes/{resource_menu_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetMenuOne(ctx *gin.Context) {
	var (
		menu *model.ResourceMenu
		err  error
	)
	id := model.ConvertStringToID(ctx.Param("resource_menu_id"))
	menu, err = service.Resource.GetResourceMenu(contexts.WithGinContext(ctx), id)
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
