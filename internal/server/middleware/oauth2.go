package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/sdk/convert"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// OAuth2AuthRequired 身份验证
func OAuth2AuthRequired(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount)
	if currentAccount == nil {
		uri := *ctx.Request.URL
		redirectURI, _ := url.Parse("/oauth2/login")
		redirectURIQuery := url.Values{}
		redirectURIQuery.Set("client_id", clientID)
		redirectURIQuery.Set("login_redirect_uri", uri.String())
		redirectURI.RawQuery = redirectURIQuery.Encode()
		ctx.Redirect(302, redirectURI.String())
		ctx.Abort()
	} else {
		ctx.Next()
	}
}

// OAuth2AuthUserinfoRequired 身份验证
func OAuth2AuthUserinfoRequired(key interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, ok := parseAuth(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Is Empty",
			})
			return
		}
		var (
			idTokenClaims *oauth2.JwtClaims
			err           error
		)
		idTokenClaims, err = oauth2.ParseJwtClaimsToken(tok, key)
		if err != nil {
			logrus.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if idTokenClaims == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if err = idTokenClaims.Valid(); err != nil {
			logrus.Errorf("token valid: %s", err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if !idTokenClaims.VerifyScope("openid", true) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrInvalidScope.Error(),
			})
			return
		}
		ctx.Set("idToken", idTokenClaims)
		ctx.Next()
	}
}

// OAuth2AuthScopeRequired 验证scope
func OAuth2AuthScopeRequired(scopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			client       *model.OAuth2Client
			clientScopes []model.Code
			err          error
		)
		clientID := ctx.Query("client_id")
		clientSecret := ctx.Query("client_secret")
		if v := ctx.PostForm("client_id"); v != "" {
			clientID = v
		}
		if v := ctx.PostForm("client_secret"); v != "" {
			clientSecret = v
		}
		client, err = service.OAuth2.GetClient(contexts.WithGinContext(ctx), model.ConvertStringToID(clientID))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrUnauthorizedClient.Error(),
			})
			return
		}
		if convert.ToString(client.ClientID) != clientID || client.ClientSecret != clientSecret {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrUnauthorizedClient.Error(),
			})
			return
		}
		clientScopes, _ = service.OAuth2.GetClientAllScopeCode(contexts.WithGinContext(ctx), model.ConvertStringToID(clientID))
		pass := false
		for i := 0; i < len(scopes); i++ {
			for j := 0; j < len(clientScopes); j++ {
				if scopes[i] == scopes[i] {
					pass = true
					goto PASS_LABEL // 跳出循环
				}
			}
		}
	PASS_LABEL:
		if !pass {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		ctx.Next()
	}
}
