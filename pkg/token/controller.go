package token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
)

func writeData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"data":   data,
	})
}

func writeError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "error",
		"error":  err.Error(),
	})
}

// AuthToken 使用code获取Token
func AuthToken(oauth2Client *oauth2.Client, redirectURI string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		clientID := ctx.Query("client_id")
		token, err := oauth2Client.TokenAuthorizationCode(code, redirectURI, clientID)
		if err != nil {
			logger.Errorf("oauth2Client.TokenAuthorizationCode: %s", err)
			writeError(ctx, err)
			return
		}
		writeData(ctx, token)
	}
}

// AuthRefreshToken 刷新Token
func AuthRefreshToken(oauth2Client *oauth2.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken := ctx.Query("refresh_token")
		token, err := oauth2Client.RefreshToken(refreshToken)
		if err != nil {
			logger.Errorf("oauth2Client.RefreshToken: %s", err)
			writeError(ctx, err)
			return
		}
		writeData(ctx, token)
	}
}
