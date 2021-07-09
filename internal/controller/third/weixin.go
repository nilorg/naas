package third

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
// @Summary		微信扫码登录重定向
// @Description	微信扫码登录重定向
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
	loginRedirectURI = url.QueryEscape(loginRedirectURI)
	values := make(url.Values)
	values.Set("appid", viper.GetString("weixin.kfpt.app_id"))
	values.Set("response_type", "code")
	values.Set("scope", "snsapi_login")
	values.Set("state", state)
	values.Set("redirect_uri", fmt.Sprintf("%s://%s/third/wx/callback?source=qrconnect&login_redirect_uri=%s", scheme, ctx.Request.Host, loginRedirectURI))
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?%s", "https://open.weixin.qq.com/connect/qrconnect", values.Encode()))
}

// ScanQrCode 扫码登录
// @Tags 		Third（第三方）
// @Summary		扫码登录
// @Description	扫码登录
// @Router /third/wx/scanqrcode [GET]
func (*weixin) ScanQrCode(ctx *gin.Context) {
	state := random.AZaz09(32)
	session := sessions.Default(ctx)
	session.Set(key.SessionQrcodeWeixinState, state)
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
	scanRedirectURI := ctx.Query("scan_redirect_uri")
	scanRedirectURI = url.QueryEscape(scanRedirectURI)
	values := make(url.Values)
	values.Set("appid", viper.GetString("weixin.fwh.app_id"))
	values.Set("response_type", "code")
	values.Set("scope", "snsapi_userinfo")
	values.Set("state", state)
	values.Set("redirect_uri", fmt.Sprintf("%s://%s/third/wx/callback?source=scanqrcode&scan_redirect_uri=%s", scheme, ctx.Request.Host, scanRedirectURI))
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", values.Encode()))
}

// CallBack 微信回调
// @Tags 		Third（第三方）
// @Summary		微信回调
// @Description	微信回调
// @Router /third/wx/callback [GET]
func (wx *weixin) CallBack(ctx *gin.Context) {

	source := ctx.Query("source")
	if source == "qrconnect" { // pc微信扫码
		wx.callBackForQrconnect(ctx)
	} else if source == "scanqrcode" { // 自研扫码
		wx.callBackForScanQrcode(ctx)
	} else {
		ctx.String(http.StatusBadRequest, "微信回调来源错误")
	}
}

func (*weixin) callBackForQrconnect(ctx *gin.Context) {
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
	session := sessions.Default(ctx)
	inState := session.Get(key.SessionWeixinSnsapiLoginState)
	if state != inState.(string) {
		ctx.String(http.StatusBadRequest, "微信回调state不匹配")
		return
	}
	su, st, err := service.User.LoginForWeixinKfptCode(contexts.WithGinContext(ctx), code)
	if err != nil {
		if err == errors.ErrThirdUserNotFound {
			session.Set(key.SessionThird, st)
			err = session.Save()
			if err != nil {
				ctx.String(http.StatusBadRequest, "用户信息存储失败")
				return
			}
			ctx.Redirect(http.StatusFound, fmt.Sprintf("/third/bind?source=qrconnect&redirect_uri=%s", url.QueryEscape(loginRedirectURI)))
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
}

func (*weixin) callBackForScanQrcode(ctx *gin.Context) {
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
	scanRedirectURI := ctx.Query("scan_redirect_uri")
	if scanRedirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	session := sessions.Default(ctx)
	inState := session.Get(key.SessionQrcodeWeixinState)
	if state != inState.(string) {
		ctx.String(http.StatusBadRequest, "微信回调state不匹配")
		return
	}
	su, st, err := service.User.LoginForWeixinFwhCode(contexts.WithGinContext(ctx), code)
	if err != nil {
		if err == errors.ErrThirdUserNotFound {
			session.Set(key.SessionThird, st)
			err = session.Save()
			if err != nil {
				ctx.String(http.StatusBadRequest, "用户信息存储失败")
				return
			}
			ctx.Redirect(http.StatusFound, fmt.Sprintf("/third/bind?source=scanqrcode&redirect_uri=%s", url.QueryEscape(scanRedirectURI)))
		} else {
			ctx.String(http.StatusBadRequest, "微信登录错误")
		}
		return
	}
	session.Set(key.SessionAccount, su)
	session.Delete(key.SessionQrcodeWeixinState)
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusBadRequest, "用户信息存储失败")
		return
	}
	ctx.Redirect(http.StatusFound, scanRedirectURI)
}

// Init 微信初始化
// @Tags 		Third（第三方）
// @Summary		微信初始化
// @Description	微信初始化
// @Router /third/wx/init [GET]
func (wx *weixin) Init(ctx *gin.Context) {
	source := ctx.Query("source")
	if source == "qrconnect" { // pc微信扫码
		wx.initForQrconnect(ctx)
	} else if source == "scanqrcode" { // 自研扫码
		wx.initForScanQrcode(ctx)
	} else {
		ctx.String(http.StatusBadRequest, "初始化来源出错")
	}
}

func (*weixin) initForQrconnect(ctx *gin.Context) {
	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	session := sessions.Default(ctx)
	stb := session.Get(key.SessionThird).(*model.SessionThirdBind)
	if stb.Type != model.UserThirdTypeWxUnionID {
		ctx.String(http.StatusBadRequest, "微信初始化不符合")
		return
	}
	if strings.TrimSpace(stb.ThirdID) == "" {
		ctx.String(http.StatusBadRequest, "未找到微信OpenID")
		return
	}
	parentCtx := contexts.WithGinContext(ctx)
	su, err := service.User.InitFromWeixinUnionID(parentCtx, stb.ThirdID)
	if err != nil {
		ctx.String(http.StatusBadRequest, "微信初始化错误")
		return
	}
	session.Set(key.SessionAccount, su)
	session.Delete(key.SessionThird)
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusBadRequest, "用户信息存储失败")
		return
	}
	ctx.Redirect(http.StatusFound, redirectURI)
}

func (*weixin) initForScanQrcode(ctx *gin.Context) {
	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		ctx.String(http.StatusBadRequest, "未找到重定向地址")
		return
	}
	session := sessions.Default(ctx)
	stb := session.Get(key.SessionThird).(*model.SessionThirdBind)
	if stb.Type != model.UserThirdTypeWxUnionID {
		ctx.String(http.StatusBadRequest, "微信初始化不符合")
		return
	}
	if strings.TrimSpace(stb.ThirdID) == "" {
		ctx.String(http.StatusBadRequest, "未找到微信OpenID")
		return
	}
	parentCtx := contexts.WithGinContext(ctx)
	su, err := service.User.InitFromWeixinUnionID(parentCtx, stb.ThirdID)
	if err != nil {
		ctx.String(http.StatusBadRequest, "微信初始化错误")
		return
	}
	session.Set(key.SessionAccount, su)
	session.Delete(key.SessionThird)
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusBadRequest, "用户信息存储失败")
		return
	}
	ctx.Redirect(http.StatusFound, redirectURI)
}
