package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/convert"
	"github.com/pkg/errors"
)

// 备注：
// 早期思路，后台管理系统不走OAuth2认证

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// PayloadFunc 设置jwt payload data
// data 是 Authenticator 返回的数据
func PayloadFunc(data interface{}) jwt.MapClaims {
	logger.Debugln("jwt PayloadFunc...")
	if v, ok := data.(*model.SessionAccount); ok {
		return jwt.MapClaims{
			"user_id":   convert.ToString(v.UserID),
			"user_name": v.UserName,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler 身份处理
func IdentityHandler(c *gin.Context) interface{} {
	logger.Debugln("jwt IdentityHandler...")
	claims := jwt.ExtractClaims(c)
	return &model.SessionAccount{
		UserID:   model.ConvertStringToID(claims["user_id"].(string)),
		UserName: claims["user_name"].(string),
	}
}

// Authenticator 身份验证
func Authenticator(ctx *gin.Context) (interface{}, error) {
	var loginVals login
	if err := ctx.ShouldBind(&loginVals); err != nil {
		return nil, errors.New("登录参数输入不全")
	}
	userName := loginVals.Username
	password := loginVals.Password

	result, err := service.Admin.Login(contexts.WithGinContext(ctx), userName, password)
	if err != nil {
		return nil, err
	}
	return &model.SessionAccount{
		UserID:   result.UserID,
		UserName: result.UserName,
	}, nil
}

// Authorizator 回调函数，该函数应执行已验证用户的授权。
//只有在身份验证成功之后。成功必归真，失败必归假。
//可选，默认成功。
func Authorizator(data interface{}, c *gin.Context) bool {
	logger.Debugf("Authorizator Data:%+v", data)

	//if v, ok := data.(*User); ok && v.UserName == "admin" {
	//	return true
	//}

	//return false
	return true
}

// Unauthorized 验证失败
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// Logout 管理员退出
func Logout(ctx *gin.Context) {

}
