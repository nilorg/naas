package token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
	"github.com/sirupsen/logrus"
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
		grantType := ctx.Query(oauth2.GrantTypeKey)
		var (
			token *oauth2.TokenResponse
			err   error
		)
		if grantType == oauth2.DeviceCodeKey {
			deviceCode := ctx.Query(oauth2.DeviceCodeKey)
			token, err = oauth2Client.TokenDeviceCode(deviceCode)
		} else if grantType == oauth2.AuthorizationCodeKey {
			code := ctx.Query(oauth2.CodeKey)
			clientID := ctx.Query(oauth2.ClientIDKey)
			token, err = oauth2Client.TokenAuthorizationCode(code, redirectURI, clientID)
		} else {
			writeError(ctx, oauth2.ErrInvalidGrant)
			return
		}
		if err != nil {
			logrus.Errorf("AuthToken: %s", err)
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
			logrus.Errorf("oauth2Client.RefreshToken: %s", err)
			writeError(ctx, err)
			return
		}
		writeData(ctx, token)
	}
}
