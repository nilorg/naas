package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type resource struct {
}

// AddWebRoute 添加web路由
// @Tags 		Resource（资源）
// @Summary		添加web路由
// @Accept  json
// @Produce	json
// @Param	resource_id		path	string	true	"资源ID"
// @Param 	body	body	service.ResourceAddWebRouteRequest	true	"body"
// @Success 200	{object}	Result
// @Router /resources/{resource_id}/web_routes [POST]
// @Security OAuth2AccessCode
func (*resource) AddWebRoute(ctx *gin.Context) {
	var (
		req service.ResourceAddWebRouteRequest
		err error
	)
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	resourceID := convert.ToUint64(ctx.Param("resource_id"))
	err = service.Resource.AddWebRoute(resourceID, &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}
