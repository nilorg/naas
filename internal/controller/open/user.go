package open

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/controller/api"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type user struct {
}

type createUserFromWeixinModel struct {
	WxUnionID string `json:"wx_union_id" gorm:"column:wx_union_id"` // 微信unionid
}

// CreateUserFromWeixin 从微信创建用户
// @Tags 		Open（开放接口）
// @Summary		从微信创建用户
// @Description	从微信创建用户
// @Accept  json
// @Produce	json
// @Param 	body	body	createUserFromWeixinModel	true	"body"
// @Success 200	{object}	api.Result
// @Router /open/users/wx [POST]
func (*user) CreateUserFromWeixin(ctx *gin.Context) {
	var (
		m   createUserFromWeixinModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		api.Writer(ctx, err)
		return
	}
	err = service.User.CreateFromWeixin(contexts.WithGinContext(ctx), m.WxUnionID)
	if err != nil {
		api.Writer(ctx, err)
		return
	}
	api.Writer(ctx, nil)
}
