package oauth2

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/sirupsen/logrus"
)

// LoginPage 登录页面
func LoginPage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	if errMsg != "" {
		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": errMsg,
		})
		return
	}

	var (
		err        error
		clientInfo *model.OAuth2ClientInfo
	)
	clientID := ctx.Query("client_id")
	clientInfo, err = service.OAuth2.GetClientInfo(contexts.WithGinContext(ctx), model.ConvertStringToID(clientID))
	if err != nil {
		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"client_info": clientInfo,
	})
}

// Login login post
func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	suser, err := service.User.Login(contexts.WithGinContext(ctx), username, password)
	if err != nil {
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}

	session := sessions.Default(ctx)
	session.Set(key.SessionAccount, suser)
	err = session.Save()
	if err != nil {
		logrus.Errorf("Login-Success-session.Save: %s", err)
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorf("SetErrorMessage: %s", err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.Redirect(http.StatusFound, ctx.Query("login_redirect_uri"))
}
