package api

import (
	"github.com/nilorg/sdk/strings"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
)

type user struct {
}

type userCreateModel struct {
	Username string `json:"username" example:"test"`
	Password string `json:"password" example:"123456"`
}

// Create 创建用户
// @Tags 		User（用户）
// @Summary		创建用户
// @Description	创建用户
// @Accept  json
// @Produce	json
// @Param 	body	body	userCreateModel	true	"body"
// @Success 200	{object}	Result
// @Router /users [POST]
// @Security OAuth2AccessCode
func (*user) Create(ctx *gin.Context) {
	var (
		m   userCreateModel
		err error
	)
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.User.Create(contexts.WithGinContext(ctx), m.Username, m.Password)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// GetOne 获取一个用户
// @Tags 		User（用户）
// @Summary		获取一个用户
// @Description	根据用户ID,获取一个用户
// @Accept  json
// @Produce	json
// @Param 	user_id	path	string	true	"user id"
// @Success 200	{object}	Result
// @Router /users/{user_id} [GET]
// @Security OAuth2AccessCode
func (*user) GetOne(ctx *gin.Context) {
	var (
		user *model.User
		err  error
	)
	userID := model.ConvertStringToID(ctx.Param("user_id"))
	user, err = service.User.GetOneByID(contexts.WithGinContext(ctx), userID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, user)
}

// Delete 删除一个用户
// @Tags 		User（用户）
// @Summary		删除一个用户
// @Description	根据用户ID,删除一个用户
// @Accept  json
// @Produce	json
// @Param 	user_id	path	string	true	"user id"
// @Success 200	{object}	Result
// @Router /users/{user_id} [DELETE]
// @Security OAuth2AccessCode
func (*user) Delete(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Param("user_id"), ",")
	var idsUint64Split []model.ID
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, model.ConvertStringToID(id))
	}
	err = service.User.DeleteByIDs(contexts.WithGinContext(ctx), idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// Update 修改一个用户
// @Tags 		User（用户）
// @Summary		修改一个用户
// @Description	根据用户ID,修改一个用户
// @Accept  json
// @Produce	json
// @Param 	user_id	path	string	true	"user id"
// @Param 	body	body	service.UserUpdateModel	true	"用户需要修改的信息"
// @Success 200	{object}	Result
// @Router /users/{user_id} [PUT]
// @Security OAuth2AccessCode
func (*user) Update(ctx *gin.Context) {
	var (
		user service.UserUpdateModel
		err  error
	)
	userID := model.ConvertStringToID(ctx.Param("user_id"))
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.User.Update(contexts.WithGinContext(ctx), userID, &user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// ListByPaged 查询用户
// @Tags 		User（用户）
// @Summary		查询用户
// @Description	查询用户翻页数据
// @Accept  json
// @Produce	json
// @Param	current		query	int	true	"当前页"
// @Param	pageSize	query	int	true	"页大小"
// @Success 200	{object}	Result{data=model.TableListData}
// @Router /users [GET]
// @Security OAuth2AccessCode
func (*user) ListByPaged(ctx *gin.Context) {
	var (
		result []*model.ResultUserInfo
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.User.ListPaged(contexts.WithGinContext(ctx), pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
