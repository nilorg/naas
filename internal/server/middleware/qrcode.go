package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/pkg/tools/key"
)

// WxQrcodeAuthRequired 微信扫码身份验证
func WxQrcodeAuthRequired(ctx *gin.Context) {
	wx := IsMicroMessenger(ctx)
	if !wx {
		ctx.String(http.StatusBadRequest, "请使用微信扫码")
		ctx.Abort()
		return
	}
	session := sessions.Default(ctx)
	account := session.Get(key.SessionAccount)
	if account == nil {
		// 重定向去微信授权
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/third/wx/scanqrcode?scan_redirect_uri=%s", url.QueryEscape(ctx.Request.RequestURI)))
		ctx.Abort()
		return
	}

	ctx.Next()
}
