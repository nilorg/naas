package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/server/auth"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/spf13/viper"
)

// NewJwtMiddleware 创建jwt授权中间件
func NewJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "naas",
		Key:             []byte(viper.GetString("jwt.secret")),
		Timeout:         viper.GetDuration("jwt.timeout") * time.Minute,
		MaxRefresh:      viper.GetDuration("jwt.max_refresh") * time.Minute,
		IdentityKey:     jwt.IdentityKey,
		PayloadFunc:     auth.PayloadFunc,
		IdentityHandler: auth.IdentityHandler,
		Authenticator:   auth.Authenticator,
		Authorizator:    auth.Authorizator,
		Unauthorized:    auth.Unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

// AdminAuthRequired 身份验证
func AdminAuthRequired(key interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, ok := parseAuth(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Is Empty",
			})
			return
		}
		var (
			tokenClaims *oauth2.JwtClaims
			err         error
		)
		tokenClaims, err = oauth2.ParseJwtClaimsToken(tok, key)
		if err != nil {
			logger.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if tokenClaims == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		if err = tokenClaims.Valid(); err != nil {
			logger.Errorf("token valid: %s", err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		ctx.Set("token", tokenClaims)
		ctx.Next()
	}
}

// AdminAuthSuperUserRequired 身份验证
// 初期先强制使用超级用户才能访问系统
func AdminAuthSuperUserRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenClaims := ctx.MustGet("token").(*oauth2.JwtClaims)
		if tokenClaims.Subject != viper.GetString("server.admin.super_user") {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		ctx.Next()
	}
}
