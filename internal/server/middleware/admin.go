package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/oauth2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// JWTAuthRequiredForAPI 身份验证
// 该中间件只适合在本项目中使用，中间件使用public key作为token验证、Redis作为过期token验证
// 如果其他项目要使用，验证需要通过grpc或者oauth2内省端点验证
func JWTAuthRequiredForAPI(key interface{}, oauth2ClientID string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, ok := parseAuth(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Is Empty",
			})
			return
		}
		var (
			exsit       bool
			tokenClaims *oauth2.JwtClaims
			err         error
		)
		rdsKey := fmt.Sprintf("oauth2_token_revocation:%s:access_token", oauth2ClientID)
		exsit, err = store.RedisClient.HExists(context.Background(), rdsKey, tok).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrServerError.Error(),
			})
			return
		}
		if exsit {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrExpiredToken.Error(),
			})
			return
		}
		tokenClaims, err = oauth2.ParseJwtClaimsToken(tok, key)
		if err != nil {
			logrus.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
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
			logrus.Errorf("token valid: %s", err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		ctx.Set("token", tokenClaims)
		ctx.Next()
	}
}

// JWTAuthRequired 身份验证
// 该中间件只适合在本项目中使用，中间件使用public key作为token验证、Redis作为过期token验证
// 如果其他项目要使用，验证需要通过grpc或者oauth2内省端点验证
func JWTAuthRequired(key interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, ok := parseAuth(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization Is Empty",
			})
			return
		}
		var (
			exsit       bool
			tokenClaims *oauth2.JwtClaims
			err         error
		)
		rdsKey := "oauth2_token_revocation:access_token"
		exsit, err = store.RedisClient.HExists(context.Background(), rdsKey, tok).Result()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrServerError.Error(),
			})
			return
		}
		if exsit {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": oauth2.ErrExpiredToken.Error(),
			})
			return
		}
		tokenClaims, err = oauth2.ParseJwtClaimsToken(tok, key)
		if err != nil {
			logrus.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
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
			logrus.Errorf("token valid: %s", err)
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
		usr, userInfo, err := service.User.GetInfoOneByCache(contexts.WithGinContext(ctx), model.ConvertStringToID(tokenClaims.Subject))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}
		if usr.Username != viper.GetString("server.admin.super_user") {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": oauth2.ErrAccessDenied.Error(),
			})
			return
		}
		ctx.Set("current_user", &model.SessionAccount{
			UserID:   usr.ID,
			UserName: usr.Username,
			Nickname: userInfo.Nickname,
			Picture:  userInfo.Picture,
		})
		ctx.Next()
	}
}
