package third

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/sdk/random"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type weixin struct {
}

// QrConnect 微信扫码登录
// @Tags 		Third（第三方）
// @Summary		微信扫码登录回调
// @Description	微信扫码登录回调
// @Router /third/wx/qrconnect [GET]
func (*weixin) QrConnect(ctx *gin.Context) {
	state := random.AZaz09(32)
	session := sessions.Default(ctx)
	session.Set(key.SessionWeixinSnsapiLoginState, state)
	err := session.Save()
	if err != nil {
		logrus.Errorln(err)
		ctx.String(http.StatusBadRequest, "生成state随机数错误")
		return
	}
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	loginRedirectURI := ctx.Query("login_redirect_uri")
	clientID := ctx.Query("client_id")
	values := make(url.Values)
	values.Set("appid", viper.GetString("weixin.fwh.appid"))
	values.Set("response_type", "code")
	values.Set("scope", "snsapi_login")
	values.Set("state", state)
	values.Set("redirect_uri", fmt.Sprintf("%s://%s/third/wx/callback?source=qrconnect&client_id=%s&login_redirect_uri=%s", scheme, ctx.Request.URL.Host, clientID, loginRedirectURI))
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", "https://open.weixin.qq.com/connect/qrconnect", values.Encode()))
}

// CallBack 微信回调
// @Tags 		Third（第三方）
// @Summary		微信回调
// @Description	微信回调
// @Router /third/wx/callback [GET]
func (*weixin) CallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.String(http.StatusBadRequest, "缺失参数code")
		return
	}
	state := ctx.Query("state")
	if code == "" {
		ctx.String(http.StatusBadRequest, "缺失参数state")
		return
	}
	loginRedirectURI := ctx.Query("login_redirect_uri")
	if loginRedirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	clientID := ctx.Query("client_id")
	source := ctx.Query("source")
	if source == "qrconnect" {
		session := sessions.Default(ctx)
		inState := session.Get(key.SessionWeixinSnsapiLoginState)
		if state != inState.(string) {
			ctx.String(http.StatusBadRequest, "微信回调state不匹配")
			return
		}
		su, err := service.User.LoginForWeixinCode(contexts.WithGinContext(ctx), code)
		if err != nil {
			if err == errors.ErrThirdUserNotFound {
				// TODO：去绑定页面，或者初始化。让用户选择
				ctx.Redirect(http.StatusFound, fmt.Sprintf("/third/wx/bind?client_id=%s&login_redirect_uri=%s", clientID, loginRedirectURI))
			} else {
				ctx.String(http.StatusBadRequest, "微信登录错误")
			}
			return
		}
		session.Set(key.SessionAccount, su)
		session.Delete(key.SessionWeixinSnsapiLoginState)
		err = session.Save()
		if err != nil {
			ctx.String(http.StatusBadRequest, "用户信息存储失败")
			return
		}
		ctx.Redirect(http.StatusFound, loginRedirectURI)
	} else {
		ctx.String(http.StatusBadRequest, "微信回调来源错误")
	}
}

// Init 微信绑定
// @Tags 		Third（第三方）
// @Summary		微信绑定
// @Description	微信绑定
// @Router /third/wx/bind [GET]
func (*weixin) Bind(ctx *gin.Context) {
	loginRedirectURI := ctx.Query("login_redirect_uri")
	if loginRedirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	clientID := ctx.Query("client_id")
	loginURI := fmt.Sprintf("/oauth2/login?client_id=%s&login_redirect_uri=%s", clientID, loginRedirectURI)
	initURI := fmt.Sprintf("/third/wx/init?client_id=%s&login_redirect_uri=%s", clientID, loginRedirectURI)
	ctx.HTML(http.StatusOK, "third_bind.tmpl", gin.H{
		"login_uri": loginURI,
		"init_uri":  initURI,
	})
}

// Init 微信初始化
// @Tags 		Third（第三方）
// @Summary		微信初始化
// @Description	微信初始化
// @Router /third/wx/init [GET]
func (*weixin) Init(ctx *gin.Context) {
	loginRedirectURI := ctx.Query("login_redirect_uri")
	if loginRedirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount)
	cu := currentAccount.(*model.SessionAccount)
	if cu.Action != model.SessionAccountActionBindWx {
		ctx.String(http.StatusBadRequest, "微信初始化不符合")
		return
	}
	if cu.WxOpenID == "" {
		ctx.String(http.StatusBadRequest, "未找到微信OpenID")
		return
	}
	parentCtx := contexts.WithGinContext(ctx)
	su, err := service.User.InitFromWeixinOpenID(parentCtx, cu.WxOpenID)
	if err != nil {
		ctx.String(http.StatusBadRequest, "微信初始化错误")
		return
	}
	session.Set(key.SessionAccount, su)
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusBadRequest, "用户信息存储失败")
		return
	}
	ctx.Redirect(http.StatusFound, loginRedirectURI)
}
