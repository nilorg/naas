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

// GetOne 获取一个资源
// @Tags 		Resource（资源）
// @Summary		获取一个资源
// @Description	根据资源ID,获取一个资源
// @Accept  json
// @Produce	json
// @Param 	resource_id	path	string	true	"resource id"
// @Success 200	{object}	Result
// @Router /resources/{resource_id} [GET]
// @Security OAuth2AccessCode
func (*resource) GetOne(ctx *gin.Context) {
	var (
		org *model.Resource
		err error
	)
	orgID := model.ConvertStringToID(ctx.Param("resource_id"))
	org, err = service.Resource.Get(contexts.WithGinContext(ctx), orgID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, org)
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
	resourceID := model.ConvertStringToID(ctx.Param("resource_id"))
	err = service.Resource.AddWebRoute(contexts.WithGinContext(ctx), resourceID, &req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*resource) Delete(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Param("res_id"), ",")
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

// ListByPaged 查询资源
// @Tags 		Resource（资源）
// @Summary		查询资源
// @Description	查询资源翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /resources [GET]
// @Security OAuth2AccessCode
func (*resource) ListByPaged(ctx *gin.Context) {
	var (
		result []*model.ResultResource
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Resource.ListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// Create 创建资源
// @Tags 		Resource（资源）
// @Summary		创建资源
// @Description	创建资源
// @Accept  json
// @Produce	json
// @Param 	body	body	service.ResourceEditModel	true	"body"
// @Success 200	{object}	Result
// @Router /resources [POST]
// @Security OAuth2AccessCode
func (*resource) Create(ctx *gin.Context) {
	var (
		m   service.ResourceEditModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.Create(contexts.WithGinContext(ctx), &m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Update 修改一个资源
// @Tags 		Resource（资源）
// @Summary		修改一个资源
// @Description	根据资源ID,修改一个资源
// @Accept  json
// @Produce	json
// @Param 	resource_id	path	string	true	"resource id"
// @Param 	body	body	service.ResourceEditModel	true	"资源需要修改的信息"
// @Success 200	{object}	Result
// @Router /resources/{resource_id} [PUT]
// @Security OAuth2AccessCode
func (*resource) Update(ctx *gin.Context) {
	var (
		org service.ResourceEditModel
		err error
	)
	resID := model.ConvertStringToID(ctx.Param("resource_id"))
	err = ctx.ShouldBindJSON(&org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Resource.Update(contexts.WithGinContext(ctx), resID, &org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}
