package api

import (
	"github.com/nilorg/sdk/strings"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/sdk/convert"
)

type user struct {
}

type userCreateModel struct {
	Username string `json:"username" example:"test"`
	Password string `json:"password" example:"123456"`
}

// Create 创建用户
// @Summary	创建用户
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
	err = service.User.Create(m.Username, m.Password)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*user) GetOne(ctx *gin.Context) {
	var (
		user *model.User
		err  error
	)
	userID := convert.ToUint64(ctx.Param("user_id"))
	user, err = service.User.GetOneCachedByID(userID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, user)
}

func (*user) Delete(ctx *gin.Context) {
	var (
		err error
	)
	idsStringSplit := strings.Split(ctx.Param("user_id"), ",")
	var idsUint64Split []uint64
	for _, id := range idsStringSplit {
		idsUint64Split = append(idsUint64Split, convert.ToUint64(id))
	}
	err = service.User.DeleteByIDs(idsUint64Split...)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*user) Update(ctx *gin.Context) {
	var (
		user service.UserUpdateModel
		err  error
	)
	userID := convert.ToUint64(ctx.Param("user_id"))
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.User.Update(userID, &user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

func (*user) ListByPaged(ctx *gin.Context) {
	var (
		result []*model.ResultUserInfo
		err    error
	)
	pagination := model.NewPagination(ctx)
	result, pagination.Total, err = service.User.ListPaged(pagination.GetSkip(), pagination.GetLimit())
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, model.NewTableListData(*pagination, result))
}
