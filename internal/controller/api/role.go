package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/strings"
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
		"recursive":   r.Recursive,
		"tree_select": r.RecursiveTreeSelect,
		"tree_node":   r.RecursiveTreeNode,
		"list":        r.List,
	})
}

// Recursive 递归
func (*role) Recursive(ctx *gin.Context) {
	organizationID := model.ConvertStringToID(ctx.Query("organization_id"))
	roles := service.Role.Recursive(contexts.WithGinContext(ctx), organizationID)
	writeData(ctx, roles)
}

// RecursiveTreeSelect 递归 tree select
func (*role) RecursiveTreeSelect(ctx *gin.Context) {
	organizationID := model.ConvertStringToID(ctx.Query("organization_id"))
	roles := service.Role.Recursive(contexts.WithGinContext(ctx), organizationID)
	treeSelects := model.RecursiveRoleToTreeSelect(roles)
	writeData(ctx, treeSelects)
}

// RecursiveTreeNode 递归 tree node
func (*role) RecursiveTreeNode(ctx *gin.Context) {
	organizationID := model.ConvertStringToID(ctx.Query("organization_id"))
	parentCtx := contexts.WithGinContext(ctx)
	roles := service.Role.Recursive(parentCtx, organizationID)
	nodes := model.RecursiveRoleToTreeNode(roles)
	id := model.ConvertStringToCode(ctx.Query("id"))
	if id != "" {
		role, err := service.Role.GetOneByCode(parentCtx, id)
		if err == nil {
			nodes = append(nodes, &model.ResultTreeNode{
				ID:     string(role.Code),
				PID:    string(role.ParentCode),
				Title:  role.Name,
				Value:  role.Code,
				IsLeaf: len(role.ChildRoles) == 0,
			})
		}
	}
	writeData(ctx, nodes)
}

// List 查询列表
func (*role) List(ctx *gin.Context) {
	var (
		result []*model.ResultRole
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.Role.ListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}

// AddResourceWebRoute 添加资源web路由
// @Tags 		Role（角色）
// @Summary		添加资源web路由
// @Accept  json
// @Produce	json
// @Param	role_code		path	string	true	"角色Code"
// @Param	resource_web_route_id		path	string	true	"资源web路由ID"
// @Success 200	{object}	Result
// @Router /roles/{role_code}/resource_web_route/{resource_web_route_id} [POST]
// @Security OAuth2AccessCode
func (*role) AddResourceWebRoute(ctx *gin.Context) {
	roleCode := ctx.Param("role_code")
	resourceWebRouteID := model.ConvertStringToID(ctx.Param("resource_web_route_id"))
	err := service.Role.AddResourceWebRoute(contexts.WithGinContext(ctx), model.Code(roleCode), resourceWebRouteID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Create 创建角色
// @Tags 		Role（角色）
// @Summary		创建角色
// @Description	创建角色
// @Accept  json
// @Produce	json
// @Param 	body	body	service.RoleEditModel	true	"body"
// @Success 200	{object}	Result
// @Router /roles [POST]
// @Security OAuth2AccessCode
func (*role) Create(ctx *gin.Context) {
	var (
		m   service.RoleEditModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Role.Create(contexts.WithGinContext(ctx), &m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Delete 删除一个角色
// @Tags 		Role（角色）
// @Summary		删除一个角色
// @Description	根据角色Code,删除一个角色
// @Accept  json
// @Produce	json
// @Param 	codes	query	string	true	"role code"
// @Success 200	{object}	Result
// @Router /roles [DELETE]
// @Security OAuth2AccessCode
func (*role) Delete(ctx *gin.Context) {
	var (
		err error
	)
	codesStringSplit := strings.Split(ctx.Query("codes"), ",")
	var codesUint64Split []model.Code
	for _, code := range codesStringSplit {
		codesUint64Split = append(codesUint64Split, model.ConvertStringToCode(code))
	}
	err = service.Role.DeleteByCodes(contexts.WithGinContext(ctx), codesUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Update 修改一个角色
// @Tags 		Role（角色）
// @Summary		修改一个角色
// @Description	根据角色Code,修改一个角色
// @Accept  json
// @Produce	json
// @Param 	body	body	service.RoleEditModel	true	"角色需要修改的信息"
// @Success 200	{object}	Result
// @Router /roles [PUT]
// @Security OAuth2AccessCode
func (*role) Update(ctx *gin.Context) {
	var (
		role service.RoleEditModel
		err  error
	)
	err = ctx.ShouldBindJSON(&role)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.Role.Update(contexts.WithGinContext(ctx), &role)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// GetOne 获取一个角色
// @Tags 		Role（角色）
// @Summary		获取一个角色
// @Description	根据角色Code,获取一个角色
// @Accept  json
// @Produce	json
// @Param 	role_code	path	string	true	"role code"
// @Success 200	{object}	Result
// @Router /roles/{role_code} [GET]
// @Security OAuth2AccessCode
func (*role) GetOne(ctx *gin.Context) {
	var (
		role *model.Role
		err  error
	)
	roleCode := model.ConvertStringToCode(ctx.Param("role_code"))
	role, err = service.Role.GetOneByCode(contexts.WithGinContext(ctx), roleCode)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, role)
}
