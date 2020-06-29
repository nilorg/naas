package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/service"
)

type role struct {
}

// QueryChildren 查询方法
// @Tags 		角色
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
	roles := service.Role.Recursive()
	ctx.JSON(http.StatusOK, roles)
}

// List 查询列表
func (*role) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
