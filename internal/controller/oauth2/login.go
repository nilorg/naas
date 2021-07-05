package oauth2

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/geetest/gt3"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/geetest"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoginPage 登录页面
func LoginPage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	geetestEnabled := viper.GetBool("geetest.enabled")
	var (
		err        error
		clientInfo *model.OAuth2ClientInfo
	)
	loginRedirectURI := ctx.Query("login_redirect_uri")
	clientID := strings.TrimSpace(ctx.Query("client_id"))
	if clientID == "" {
		ctx.String(http.StatusBadRequest, "oauth2 client_id is empty")
		return
	}
	clientInfo, err = service.OAuth2.GetClientInfo(contexts.WithGinContext(ctx), model.ConvertStringToID(clientID))
	if errMsg != "" {
		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error":              errMsg,
			"client_info":        clientInfo,
			"client_id":          clientID,
			"login_redirect_uri": loginRedirectURI,
			"geetest_enabled":    geetestEnabled,
		})
		return
	} else if err != nil {
		ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error":              err.Error(),
			"client_id":          clientID,
			"login_redirect_uri": loginRedirectURI,
			"geetest_enabled":    geetestEnabled,
		})
		return
	}

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"client_info":        clientInfo,
		"client_id":          clientID,
		"login_redirect_uri": loginRedirectURI,
		"geetest_enabled":    geetestEnabled,
	})
}

// Login login post
func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	var err error
	session := sessions.Default(ctx)

	// 极验验证
	if viper.GetBool("geetest.enabled") {
		challenge := ctx.PostForm(gt3.GeetestChallenge)
		seccode := ctx.PostForm(gt3.GeetestSeccode)
		status := session.Get(gt3.GeetestServerStatusSessionKey)
		if status == nil {
			err = SetErrorMessage(ctx, "未找到极验验证授权信息")
			if err != nil {
				logrus.Errorln(err)
			}
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		} else if status.(int) == 1 {
			var res *gt3.ValidateResponse
			res, err = geetest.GeetestClient.Validate(challenge, seccode)
			if err != nil {
				err = SetErrorMessage(ctx, err.Error())
				if err != nil {
					logrus.Errorln(err)
				}
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
				return
			}
			if res.Seccode == "false" {
				err = SetErrorMessage(ctx, "验证码未通过")
				if err != nil {
					logrus.Errorln(err)
				}
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			}
		} else {
			err = SetErrorMessage(ctx, "极验验证授权状态错误")
			if err != nil {
				logrus.Errorln(err)
			}
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		}
	}
	// 登录验证
	suser, err := service.User.Login(contexts.WithGinContext(ctx), username, password)
	if err != nil {
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
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
