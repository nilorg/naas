package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type casbin struct {
}

// GetScopeOne 获取资源服务器Web路由
// @Tags 		Casbin
// @Summary		获取资源服务器Web路由翻页列表
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/resource/{resource_server_id}/web_routes [GET]
// @Security OAuth2AccessCode
func (*casbin) ListResourceWebRoutes(ctx *gin.Context) {
	var (
		list []*model.ResourceWebRoute
		err  error
	)
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Casbin.ListResourceWebRoutePagedByResourceServerID(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit(), resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// GetScopeOne 根据角色获取资源服务器Web路由
// @Tags 		Casbin
// @Summary		根据角色获取资源服务器Web路由翻页列表
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/role/{role_code}/resource/{resource_server_id}/web_routes [GET]
// @Security OAuth2AccessCode
func (*casbin) ListRoleResourceWebRoutes(ctx *gin.Context) {
	var (
		list []*model.RoleResourceRelation
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	list, err = service.Casbin.ListRoleResourceWebRouteByRoleCodeAndResourceServerID(contexts.WithGinContext(ctx), roleCode, resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, list)
}

// AddResourceWebRoute 添加资源web路由
// @Tags 		Casbin
// @Summary		添加资源web路由
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Success 200	{object}	Result
// @Router /casbin/role/{role_code}/resource_web_routes [PUT]
// @Security OAuth2AccessCode
func (*casbin) AddResourceWebRoute(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	var (
		req service.CasbinAddResourceWebRouteModel
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Casbin.AddResourceWebRoute(contexts.WithGinContext(ctx), model.Code(roleCode), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// AddResourceWebMenu 添加资源web菜单
// @Tags 		Casbin
// @Summary		添加资源web菜单
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Success 200	{object}	Result
// @Router /casbin/role/{role_code}/resource_web_menus [PUT]
// @Security OAuth2AccessCode
func (*casbin) AddResourceWebMenu(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	var (
		req service.CasbinAddResourceWebMenuModel
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Casbin.AddResourceWebMenu(contexts.WithGinContext(ctx), model.Code(roleCode), &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// GetScopeOne 获取资源服务器Web菜单
// @Tags 		Casbin
// @Summary		获取资源服务器Web菜单翻页列表
// @Accept  json
// @Produce	json
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/resource/{resource_server_id}/web_menus [GET]
// @Security OAuth2AccessCode
func (*casbin) ListResourceWebMenus(ctx *gin.Context) {
	var (
		list []*model.ResourceWebMenu
		err  error
	)
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	pagination := model.NewPagination(ctx)
	list, pagination.Total, err = service.Casbin.ListResourceWebMenuPagedByResourceServerID(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit(), resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, list))
}

// GetScopeOne 根据角色获取资源服务器Web菜单
// @Tags 		Casbin
// @Summary		根据角色获取资源服务器Web菜单翻页列表
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Param 	resource_server_id	path	string	true	"resource server id"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /casbin/role/{role_code}/resource/{resource_server_id}/web_menus [GET]
// @Security OAuth2AccessCode
func (*casbin) ListRoleResourceWebMenus(ctx *gin.Context) {
	var (
		list []*model.RoleResourceRelation
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	resourceServerID := model.ConvertStringToID(ctx.Param("resource_server_id"))
	list, err = service.Casbin.ListRoleResourceWebMenuByRoleCodeAndResourceServerID(contexts.WithGinContext(ctx), roleCode, resourceServerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, list)
}
