package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type casbin struct {
}

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
func (*casbin) ListRoleResourceWebRoutes(ctx *gin.Context) {
	var (
		list []*model.RoleResourceWebRoute
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
