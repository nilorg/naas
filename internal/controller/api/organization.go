package api

import (
	"github.com/nilorg/sdk/strings"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type organization struct {
}

// GetOne 获取一个组织
// @Tags 		Organization（组织）
// @Summary		获取一个组织
// @Description	根据组织ID,获取一个组织
// @Accept  json
// @Produce	json
// @Param 	org_id	path	string	true	"org id"
// @Success 200	{object}	Result
// @Router /organizations/{org_id} [GET]
// @Security OAuth2AccessCode
func (*organization) GetOne(ctx *gin.Context) {
	var (
		org *model.Organization
		err error
	)
	orgID := model.ConvertStringToID(ctx.Param("org_id"))
	org, err = service.Organization.GetOneByID(contexts.WithGinContext(ctx), orgID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, org)
}

// Create 创建组织
// @Tags 		Organization（组织）
// @Summary		创建组织
// @Description	创建组织
// @Accept  json
// @Produce	json
// @Param 	body	body	service.OrganizationEditModel	true	"body"
// @Success 200	{object}	Result
// @Router /organizations [POST]
// @Security OAuth2AccessCode
func (*organization) Create(ctx *gin.Context) {
	var (
		m   service.OrganizationEditModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Organization.Create(contexts.WithGinContext(ctx), &m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Delete 删除一个组织
// @Tags 		Organization（组织）
// @Summary		删除一个组织
// @Description	根据组织ID,删除一个组织
// @Accept  json
// @Produce	json
// @Param 	org_id	path	string	true	"organization id"
// @Success 200	{object}	Result
// @Router /organizations/{org_id} [DELETE]
// @Security OAuth2AccessCode
func (*organization) Delete(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Param("org_id"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.Organization.DeleteByIDs(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Update 修改一个组织
// @Tags 		Organization（组织）
// @Summary		修改一个组织
// @Description	根据组织ID,修改一个组织
// @Accept  json
// @Produce	json
// @Param 	org_id	path	string	true	"organization id"
// @Param 	body	body	service.OrganizationEditModel	true	"组织需要修改的信息"
// @Success 200	{object}	Result
// @Router /organizations/{org_id} [PUT]
// @Security OAuth2AccessCode
func (*organization) Update(ctx *gin.Context) {
	var (
		org service.OrganizationEditModel
		err error
	)
	orgID := model.ConvertStringToID(ctx.Param("org_id"))
	err = ctx.ShouldBindJSON(&org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Organization.Update(contexts.WithGinContext(ctx), orgID, &org)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// QueryChildren 查询方法
// @Tags 		Organization（组织）
// @Summary		查询组织
// @Description	recursive:递归获取所有组织
// @Description	list:查询列表
// @Accept  json
// @Produce	json
// @Param q query string true "查询参数" Enums(recursive,list)
// @Success 200	{object}	Result{data=model.TableListData}
// @Success 200	{object}	Result{data=[]model.Organization}
// @Router /organizations [GET]
// @Security OAuth2AccessCode
func (o *organization) QueryChildren() gin.HandlerFunc {
	return QueryChildren(map[string]gin.HandlerFunc{
		"list": o.ListByPaged,
	})
}

func (*organization) ListByPaged(ctx *gin.Context) {
	var (
		result []*model.ResultOrganization
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Organization.ListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
