package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type role struct {
}

// QueryChildren 查询方法
// @Tags 		Role（角色）
// @Summary		查询角色
// @Description	recursive:递归获取所有角色
// @Description	list:查询列表
// @Accept  json
// @Produce	json
// @Param q query string true "查询参数" Enums(recursive,list)
// @Success 200	{object}	Result{data=model.TableListData}
// @Success 200	{object}	Result{data=[]model.Role}
// @Router /roles [GET]
// @Security OAuth2AccessCode
func (r *role) QueryChildren() gin.HandlerFunc {
	return QueryChildren(map[string]gin.HandlerFunc{
		"recursive": r.Recursive,
		"list":      r.List,
	})
}

// Recursive 递归
func (*role) Recursive(ctx *gin.Context) {
	roles := service.Role.Recursive(contexts.WithGinContext(ctx))
	ctx.JSON(http.StatusOK, roles)
}

// List 查询列表
func (*role) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}

// AddResourceWebRoute 添加资源web路由
// @Tags 		Role（角色）
// @Summary		添加资源web路由
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Param	resource_web_route_id		path	string	true	"资源web路由ID"
// @Param 	body	body	service.ResourceAddWebRouteRequest	true	"body"
// @Success 200	{object}	Result
// @Router /roles/{role_code}/resource_web_route/{resource_web_route_id} [POST]
// @Security OAuth2AccessCode
func (*role) AddResourceWebRoute(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	resourceWebRouteID := convert.ToUint64(ctx.Param("resource_web_route_id"))
	err := service.Role.AddResourceWebRoute(contexts.WithGinContext(ctx), roleCode, resourceWebRouteID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}
