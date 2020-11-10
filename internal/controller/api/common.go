package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type common struct {
}

// SelectQueryChildren 查询方法
// @Tags 		Role（角色）
// @Summary		查询角色
// @Description	org:组织
// @Description	xxxx:其他
// @Accept  json
// @Produce	json
// @Param q query string true "查询参数" Enums(org,xxxx)
// @Param pageSize query int true "页大小"
// @Success 200	{object}	Result{data=model.ResultSelect}
// @Router /common/select [GET]
// @Security OAuth2AccessCode
func (c *common) SelectQueryChildren() gin.HandlerFunc {
	return QueryChildren(map[string]gin.HandlerFunc{
		"organization":    c.SelectOrganizationList,
		"resource_server": c.SelectResourceServerList,
	})
}

// SelectOrganizationList 组织Select列表
func (*common) SelectOrganizationList(ctx *gin.Context) {
	id := ctx.Query("id")
	parentCtx := contexts.WithGinContext(ctx)
	if id != "" {
		org, err := service.Organization.GetOneByID(parentCtx, model.ConvertStringToID(id))
		if err != nil {
			writeError(ctx, err)
		} else {
			writeData(ctx, &model.ResultSelect{
				Label: org.Name,
				Value: org.ID,
			})
		}
	} else {
		name := ctx.Query("name")
		limit := convert.ToInt(ctx.Query("limit"))
		list, err := service.Organization.ListByName(parentCtx, name, limit)
		if err != nil {
			writeError(ctx, err)
		} else {
			var results []*model.ResultSelect
			for _, item := range list {
				results = append(results, &model.ResultSelect{
					Label: item.Name,
					Value: item.ID,
				})
			}
			writeData(ctx, results)
		}
	}
}

// SelectResourceServerList 资源服务器Select列表
func (*common) SelectResourceServerList(ctx *gin.Context) {
	id := ctx.Query("id")
	parentCtx := contexts.WithGinContext(ctx)
	if id != "" {
		res, err := service.Resource.GetServer(parentCtx, model.ConvertStringToID(id))
		if err != nil {
			writeError(ctx, err)
		} else {
			writeData(ctx, &model.ResultSelect{
				Label: res.Name,
				Value: res.ID,
			})
		}
	} else {
		name := ctx.Query("name")
		limit := convert.ToInt(ctx.Query("limit"))
		list, err := service.Resource.ListByName(parentCtx, name, limit)
		if err != nil {
			writeError(ctx, err)
		} else {
			var results []*model.ResultSelect
			for _, item := range list {
				results = append(results, &model.ResultSelect{
					Label: item.Name,
					Value: item.ID,
				})
			}
			writeData(ctx, results)
		}
	}
}
