package oauth2

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/core/util/key"
	"github.com/nilorg/naas/service"
	"github.com/nilorg/pkg/logger"
)

// LoginPage 登录页面
func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"error": GetErrorMessage(ctx),
	})
}

// Login login post
func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	suser, err := service.User.Login(username, password)
	if err != nil {
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logger.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}

	session := sessions.Default(ctx)
	session.Set(key.SessionAccount, suser)
	err = session.Save()
	if err != nil {
		_ = SetErrorMessage(ctx, err.Error())
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.Redirect(302, ctx.Query("login_redirect_uri"))
}
