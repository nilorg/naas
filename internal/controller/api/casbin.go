package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type casbin struct {
}

// ListResourceRoutes 获取资源服务器路由
// @Tags 		Casbin
// @Summary		获取资源服务器路由翻页列表
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/resource/{resource_server_id}/routes [GET]
// @Security OAuth2AccessCode
func (*casbin) ListResourceRoutes(ctx *gin.Context) {
	var (
		list []*model.ResourceRoute
		err  error
	)
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Casbin.ListResourceRoutePagedByResourceServerID(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit(), resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// ListRoleResourceRoutes 根据角色获取资源服务器路由
// @Tags 		Casbin
// @Summary		根据角色获取资源服务器路由翻页列表
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/role/{role_code}/resource/{resource_server_id}/routes [GET]
// @Security OAuth2AccessCode
func (*casbin) ListRoleResourceRoutes(ctx *gin.Context) {
	var (
		list []*model.RoleResourceRelation
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	list, err = service.Casbin.ListRoleResourceRouteByRoleCodeAndResourceServerID(contexts.WithGinContext(ctx), roleCode, resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, list)
}

// AddResourceRoute 添加资源路由
// @Tags 		Casbin
// @Summary		添加资源路由
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Success 200	{object}	Result
// @Router /casbin/role/{role_code}/resource_routes [PUT]
// @Security OAuth2AccessCode
func (*casbin) AddResourceRoute(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	var (
		req service.CasbinAddResourceRouteModel
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Casbin.AddResourceRoute(contexts.WithGinContext(ctx), model.Code(roleCode), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ===============

// AddResourceMenu 添加资源菜单
// @Tags 		Casbin
// @Summary		添加资源菜单
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Success 200	{object}	Result
// @Router /casbin/role/{role_code}/resource_menus [PUT]
// @Security OAuth2AccessCode
func (*casbin) AddResourceMenu(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	var (
		req service.CasbinAddResourceMenuModel
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Casbin.AddResourceMenu(contexts.WithGinContext(ctx), model.Code(roleCode), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListResourceMenus 获取资源服务器菜单
// @Tags 		Casbin
// @Summary		获取资源服务器菜单翻页列表
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/resource/{resource_server_id}/menus [GET]
// @Security OAuth2AccessCode
func (*casbin) ListResourceMenus(ctx *gin.Context) {
	var (
		list []*model.ResourceMenu
		err  error
	)
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Casbin.ListResourceMenuPagedByResourceServerID(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit(), resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// ListRoleResourceMenus 根据角色获取资源服务器菜单
// @Tags 		Casbin
// @Summary		根据角色获取资源服务器菜单翻页列表
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/role/{role_code}/resource/{resource_server_id}/menus [GET]
// @Security OAuth2AccessCode
func (*casbin) ListRoleResourceMenus(ctx *gin.Context) {
	var (
		list []*model.RoleResourceRelation
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	list, err = service.Casbin.ListRoleResourceMenuByRoleCodeAndResourceServerID(contexts.WithGinContext(ctx), roleCode, resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, list)
}

// ===============

// AddResourceAction 添加资源动作
// @Tags 		Casbin
// @Summary		添加资源动作
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Success 200	{object}	Result
// @Router /casbin/role/{role_code}/resource_actions [PUT]
// @Security OAuth2AccessCode
func (*casbin) AddResourceAction(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	var (
		req service.CasbinAddResourceActionModel
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Casbin.AddResourceAction(contexts.WithGinContext(ctx), model.Code(roleCode), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListResourceActions 获取资源服务器动作翻页列表
// @Tags 		Casbin
// @Summary		获取资源服务器动作翻页列表
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/resource/{resource_server_id}/actions [GET]
// @Security OAuth2AccessCode
func (*casbin) ListResourceActions(ctx *gin.Context) {
	var (
		list []*model.ResourceAction
		err  error
	)
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Casbin.ListResourceActionPagedByResourceServerID(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit(), resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// ListRoleResourceActions 根据角色获取资源服务器动作
// @Tags 		Casbin
// @Summary		根据角色获取资源服务器动作
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/role/{role_code}/resource/{resource_server_id}/actions [GET]
// @Security OAuth2AccessCode
func (*casbin) ListRoleResourceActions(ctx *gin.Context) {
	var (
		list []*model.RoleResourceRelation
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	list, err = service.Casbin.ListRoleResourceActionByRoleCodeAndResourceServerID(contexts.WithGinContext(ctx), roleCode, resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, list)
}
