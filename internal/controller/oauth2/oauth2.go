package oauth2

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/pkg/logger"

	"github.com/nilorg/naas/pkg/tools"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/token"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
)

var (
	oauth2Server *oauth2.Server
	sourceScope  = []string{
		"read_write",
		"read",
		"write",
	}
)

// Init 初始化
func Init() {
	oauth2Server = oauth2.NewServer(
		oauth2.ServerIssuer(viper.GetString("server.oauth2.issuer")),
		oauth2.ServerDeviceAuthorizationEndpointEnabled(viper.GetBool("server.oauth2.device_authorization_endpoint_enabled")),
		oauth2.ServerIntrospectEndpointEnabled(viper.GetBool("server.oauth2.introspection_endpoint_enabled")),
		oauth2.ServerTokenRevocationEnabled(viper.GetBool("server.oauth2.revocation_endpoint_enabled")),
	)
	oauth2Server.VerifyClient = func(basic *oauth2.ClientBasic) (err error) {
		var client *model.OAuth2Client
		client, err = service.OAuth2.GetClient(basic.ID)
		if err != nil {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		if convert.ToString(client.ClientID) != basic.ID || client.ClientSecret != basic.Secret {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		return
	}
	oauth2Server.VerifyPassword = func(username, password string) (openID string, err error) {
		var user *model.User
		user, err = service.User.GetUserByUsername(username)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if user.Username != username || user.Password != password {
			err = oauth2.ErrAccessDenied
		}
		return
	}
	oauth2Server.VerifyRedirectURI = func(clientID, redirectURI string) (err error) {
		var client *model.OAuth2Client
		client, err = service.OAuth2.GetClient(clientID)
		if err != nil {
			err = oauth2.ErrAccessDenied
			return
		}
		if strings.Index(redirectURI, client.RedirectURI) == -1 {
			err = oauth2.ErrInvalidRedirectURI
		}
		return
	}
	oauth2Server.GenerateCode = func(clientID, openID, redirectURI string, scope []string) (code string, err error) {
		code = oauth2.RandomCode()
		value := &oauth2.CodeValue{
			ClientID:    clientID,
			OpenID:      openID,
			RedirectURI: redirectURI,
			Scope:       scope,
		}
		err = store.RedisClient.Set(key.WrapOAuth2Code(code), value, time.Minute).Err()
		if err != nil {
			logger.Errorf("store.RedisClient.Set Error: %s", err)
			err = oauth2.ErrServerError
		}
		return
	}
	oauth2Server.VerifyCode = func(code, clientID, redirectURI string) (value *oauth2.CodeValue, err error) {
		value = &oauth2.CodeValue{}
		redisKey := key.WrapOAuth2Code(code)
		err = store.RedisClient.Get(redisKey).Scan(value)
		if err != nil {
			logger.Errorf("store.RedisClient.Get Error: %s", err)
			err = oauth2.ErrAccessDenied
			return
		}
		// 删除Key
		_ = store.RedisClient.Del(redisKey)
		if value.ClientID != clientID || (strings.HasPrefix(redirectURI, value.RedirectURI) && redirectURI != value.RedirectURI) {
			err = oauth2.ErrAccessDenied
			return
		}
		// 包含
		if !tools.ScopeIsSubset(value.Scope, sourceScope) {
			err = oauth2.ErrInvalidScope
		}
		return
	}
	oauth2Server.VerifyScope = func(scope []string) (err error) {
		// 表示权限范围，如果与客户端申请的范围一致，此项可省略。
		if len(scope) == 0 {
			return
		}
		if !tools.ScopeIsSubset(scope, sourceScope) {
			err = oauth2.ErrInvalidScope
		}
		return
	}
	oauth2Server.GenerateDeviceAuthorization = func(issuer, verificationURI, clientID, scope string) (resp *oauth2.DeviceAuthorizationResponse, err error) {
		return
	}
	oauth2Server.VerifyDeviceCode = func(deviceCode, clientID string) (value *oauth2.DeviceCodeValue, err error) {
		return
	}
	oauth2Server.VerifyIntrospectionToken = func(token, clientID string, tokenTypeHint ...string) (resp *oauth2.IntrospectionResponse, err error) {
		return
	}
	oauth2Server.TokenRevocation = func(token, clientID string, tokenTypeHint ...string) {

	}
	oauth2Server.GenerateAccessToken = token.NewGenerateAccessToken(global.JwtPrivateKey)
	oauth2Server.RefreshAccessToken = token.NewRefreshAccessToken(global.JwtPrivateKey)
	oauth2Server.ParseAccessToken = token.NewParseAccessToken(global.JwtPrivateKey)
	oauth2Server.Init()
}

// SetErrorMessage set a error message
func SetErrorMessage(ctx *gin.Context, msg string) error {
	session := sessions.Default(ctx)
	session.Set("error_message", msg)
	return session.Save()
}

// GetErrorMessage return the first error message
func GetErrorMessage(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	value := session.Get("error_message")
	if value != nil {
		session.Delete("error_message")
		_ = session.Save()
		return value.(string)
	}
	return ""
}

// Token ...
func Token(ctx *gin.Context) {
	oauth2Server.HandleToken(ctx.Writer, ctx.Request)
}

// DeviceCode ...
func DeviceCode(ctx *gin.Context) {
	oauth2Server.HandleDeviceAuthorization(ctx.Writer, ctx.Request)
}

// TokenIntrospection ...
func TokenIntrospection(ctx *gin.Context) {
	oauth2Server.HandleTokenIntrospection(ctx.Writer, ctx.Request)
}

// TokenRevoke ...
func TokenRevoke(ctx *gin.Context) {
	oauth2Server.HandleTokenRevocation(ctx.Writer, ctx.Request)
}
