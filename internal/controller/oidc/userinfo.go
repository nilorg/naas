package oidc

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/oidc"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/sdk/convert"
)

// GetUserinfo ...
func GetUserinfo(ctx *gin.Context) {
	var (
		idTokenClaims  *oauth2.JwtClaims
		resultUserinfo oidc.Userinfo
		err            error
	)
	idTokenClaims = ctx.MustGet("idToken").(*oauth2.JwtClaims)
	user, err := service.User.GetOneByID(contexts.WithGinContext(ctx), convert.ToUint64(idTokenClaims.Subject))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.Status(404)
		} else {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	resultUserinfo.Sub = convert.ToString(user.ID)
	userInfoFlag := idTokenClaims.VerifyScope("profile", true) || idTokenClaims.VerifyScope("email", true)
	if userInfoFlag {
		userInfo, err := service.User.GetInfoOneByUserID(contexts.WithGinContext(ctx), convert.ToUint64(idTokenClaims.Subject))
		if err == nil && userInfo == nil {
			ctx.JSON(200, resultUserinfo)
			return
		}
		if idTokenClaims.VerifyScope("profile", true) {
			resultUserinfo.Name = user.Username
			resultUserinfo.Nickname = userInfo.Nickname
			resultUserinfo.Picture = userInfo.Picture
			resultUserinfo.Gender = convert.ToString(userInfo.Gender)
			resultUserinfo.UpdatedAt = userInfo.UpdatedAt.Unix()
		}
		if idTokenClaims.VerifyScope("email", true) {
			resultUserinfo.Email = userInfo.Email
			resultUserinfo.EmailVerified = userInfo.EmailVerified
		}
		if idTokenClaims.VerifyScope("phone", true) {
			resultUserinfo.PhoneNumber = userInfo.Phone
			resultUserinfo.PhoneNumberVerified = userInfo.PhoneVerified
		}
	}
	ctx.JSON(200, resultUserinfo)
}
