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

func (*user) Create(ctx *gin.Context) {
	var (
		user model.User
		err  error
	)
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.User.Create(user.Username, user.Password)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writeData(ctx, nil)
}

// UserUpdateModel ...
type UserUpdateModel struct {
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
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
		user UserUpdateModel
		err  error
	)
	userID := convert.ToUint64(ctx.Param("user_id"))
	err = ctx.ShouldBindJSON(&user)
	if err != nil {
		writeError(ctx, err)
		return
	}
	err = service.User.Update(userID, user.Username, user.Password)
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
