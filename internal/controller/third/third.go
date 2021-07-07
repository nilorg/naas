package third

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/geetest/gt3"
	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/geetest"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/pkg/ginextension"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Weixin = &weixin{}
	QrCode = &qrcode{}
)

// BindPage 绑定页面
func BindPage(ctx *gin.Context) {
	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	source := ctx.Query("source")
	loginURI := fmt.Sprintf("/third/login?source=%s&login_redirect_uri=%s", source, redirectURI)
	wx := ginextension.IsMicroMessenger(ctx)
	initURI := fmt.Sprintf("/third/wx/init?source=%s&redirect_uri=%s", source, redirectURI)
	ctx.HTML(http.StatusOK, "third_bind.tmpl", gin.H{
		"login_uri":   loginURI,
		"wx_status":   wx,
		"wx_init_uri": initURI,
	})
}

// LoginPage 登录页面
func LoginPage(ctx *gin.Context) {
	errMsg := oauth2.GetErrorMessage(ctx)
	geetestEnabled := viper.GetBool("geetest.enabled")
	if errMsg != "" {
		ctx.HTML(http.StatusOK, "third_login.tmpl", gin.H{
			"error":           errMsg,
			"geetest_enabled": geetestEnabled,
		})
		return
	}

	ctx.HTML(http.StatusOK, "third_login.tmpl", gin.H{
		"geetest_enabled": geetestEnabled,
	})
}

// Login login post
func Login(ctx *gin.Context) {
	loginRedirectURI := ctx.Query("login_redirect_uri")
	if loginRedirectURI == "" {
		ctx.String(http.StatusBadRequest, "第三方绑定，未找到重定向地址")
		return
	}
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
			err = oauth2.SetErrorMessage(ctx, "未找到极验验证授权信息")
			if err != nil {
				logrus.Errorln(err)
			}
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		} else if status.(int) == 1 {
			var res *gt3.ValidateResponse
			res, err = geetest.GeetestClient.Validate(challenge, seccode)
			if err != nil {
				err = oauth2.SetErrorMessage(ctx, err.Error())
				if err != nil {
					logrus.Errorln(err)
				}
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
				return
			}
			if res.Seccode == "false" {
				err = oauth2.SetErrorMessage(ctx, "验证码未通过")
				if err != nil {
					logrus.Errorln(err)
				}
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			}
		} else {
			err = oauth2.SetErrorMessage(ctx, "极验验证授权状态错误")
			if err != nil {
				logrus.Errorln(err)
			}
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		}
	}
	parentCtx := contexts.WithGinContext(ctx)
	stb := session.Get(key.SessionThird).(*model.SessionThirdBind)
	// 登录验证
	suser, err := service.User.Login(parentCtx, username, password)
	if err != nil {
		err = oauth2.SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorln(err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	session.Set(key.SessionAccount, suser)
	// 绑定第三方账户
	err = service.User.BindThird(parentCtx, suser.UserID, stb.ThirdID, stb.Type)
	if err != nil {
		logrus.Errorf("service.User.BindThird: %s", err)
		err = oauth2.SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorf("SetErrorMessage: %s", err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	session.Delete(key.SessionThird)
	err = session.Save()
	if err != nil {
		logrus.Errorf("Login-Success-session.Save: %s", err)
		err = oauth2.SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorf("SetErrorMessage: %s", err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.Redirect(http.StatusFound, loginRedirectURI)
}
