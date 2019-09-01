package middleware

import (
	"net/url"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// Header 头处理
func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", "Nilorg OAuth2")
		c.Next()
	}
}

// AuthRequired 身份验证
func AuthRequired(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
	session := sessions.Default(ctx)
	currentAccount := session.Get("current_user")
	if currentAccount == nil {
		uri := *ctx.Request.URL
		redirectURI, _ := url.Parse("/oauth2/login")
		redirectURIQuery := url.Values{}
		redirectURIQuery.Set("client_id", clientID)
		redirectURIQuery.Set("login_redirect_uri", uri.String())
		redirectURI.RawQuery = redirectURIQuery.Encode()
		ctx.Redirect(302, redirectURI.String())
		return
	}
	ctx.Next()
}
