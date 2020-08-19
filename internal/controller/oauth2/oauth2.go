package oauth2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/pkg/logger"

	"github.com/nilorg/pkg/slice"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/pkg/token"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/oauth2"
)

var (
	oauth2Server *oauth2.Server
	//// SourceScope ...
	//SourceScope = []string{
	//	"openid",
	//	"profile",
	//	"email",
	//	"phone",
	//}
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
		client, err = service.OAuth2.GetClient(contexts.WithContext(context.Background()), model.ConvertStringToID(basic.ID))
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
		user, err = service.User.GetUserByUsername(contexts.WithContext(context.Background()), username)
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
		client, err = service.OAuth2.GetClient(contexts.WithContext(context.Background()), model.ConvertStringToID(clientID))
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
			Scope:       RemoveRepeat(scope),
		}
		err = store.RedisClient.Set(context.Background(), key.WrapOAuth2Code(code), value, time.Minute).Err()
		if err != nil {
			logger.Errorf("store.RedisClient.Set Error: %s", err)
			err = oauth2.ErrServerError
		}
		return
	}
	oauth2Server.VerifyCode = func(code, clientID, redirectURI string) (value *oauth2.CodeValue, err error) {
		value = &oauth2.CodeValue{}
		redisKey := key.WrapOAuth2Code(code)
		ctx := context.Background()
		err = store.RedisClient.Get(ctx, redisKey).Scan(value)
		if err != nil {
			logger.Errorf("store.RedisClient.Get Error: %s", err)
			err = oauth2.ErrAccessDenied
			return
		}
		// 删除Key
		_ = store.RedisClient.Del(ctx, redisKey)
		if value.ClientID != clientID || (strings.HasPrefix(redirectURI, value.RedirectURI) && redirectURI != value.RedirectURI) {
			err = oauth2.ErrAccessDenied
		}
		return
	}
	oauth2Server.VerifyScope = func(scope []string, clientID string) (err error) {
		// 表示权限范围，如果与客户端申请的范围一致，此项可省略。
		if len(scope) == 0 {
			return
		}
		var scopes []model.Code
		scopes, err = service.OAuth2.GetClientAllScopeCode(contexts.WithContext(context.Background()), model.ConvertStringToID(clientID))
		if err != nil {
			err = oauth2.ErrInvalidScope
			return
		}
		if !slice.IsSubset(scope, model.ConvertCodeSliceToStringSlice(scopes)) {
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
		logger.Debugf("oauth2Server.VerifyIntrospectionToken....")
		key := tokenRevocationKey(token, clientID, tokenTypeHint...)
		var exsit bool
		exsit, err = store.RedisClient.HExists(context.Background(), key, token).Result()
		if err != nil {
			logger.Errorf("store.RedisClient.HExists: %s", err)
			err = oauth2.ErrServerError
			return
		}
		if exsit {
			err = oauth2.ErrExpiredToken
			return
		}
		var tokenClaims *oauth2.JwtClaims
		tokenClaims, err = oauth2.ParseJwtClaimsToken(token, global.JwtPrivateKey.Public())
		if err != nil {
			logger.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			err = oauth2.ErrServerError
			return
		}
		if !tokenClaims.VerifyAudience([]string{clientID}, false) {
			logger.Debugf("tokenClaims.VerifyAudience.....false")
			err = oauth2.ErrInvalidClient
		}
		resp = new(oauth2.IntrospectionResponse)
		resp.Active = true
		if verr := tokenClaims.Valid(); verr != nil {
			logger.Debugf("tokenClaims.Valid: %s", verr)
			resp.Active = false
			return
		}
		resp.ClientID = clientID
		resp.Scope = tokenClaims.Scope
		resp.Sub = tokenClaims.Subject
		resp.Aud = clientID
		resp.Exp = tokenClaims.ExpiresAt
		resp.Iss = tokenClaims.IssuedAt
		var user *model.User
		user, err = service.User.GetOneByID(contexts.WithContext(context.Background()), model.ConvertStringToID(tokenClaims.Subject))
		if err == nil && user != nil {
			resp.Username = user.Username
		}
		return
	}
	oauth2Server.TokenRevocation = func(token, clientID string, tokenTypeHint ...string) {
		key := tokenRevocationKey(token, clientID, tokenTypeHint...)
		tokenClaims, err := oauth2.ParseJwtClaimsToken(token, global.JwtPrivateKey.Public())
		if err != nil {
			logger.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			return
		}
		exp := time.Unix(tokenClaims.ExpiresAt, 0)
		err = store.RedisClient.HSet(context.Background(), key, token, time.Now().Sub(exp)).Err()
		if err != nil {
			logger.Errorf("store.RedisClient.Set: %s", err)
			return
		}
	}
	oauth2Server.GenerateAccessToken = token.NewGenerateAccessToken(global.JwtPrivateKey, viper.GetBool("server.oidc.enabled") && viper.GetBool("server.oidc.userinfo_endpoint_enabled"))
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

// RemoveRepeat 过滤重复元素
func RemoveRepeat(slc []string) (result []string) {
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func tokenRevocationKey(token, clientID string, tokenTypeHint ...string) string {
	if len(tokenTypeHint) > 0 {
		return fmt.Sprintf("oauth2_token_revocation:%s:%s", clientID, tokenTypeHint[0])
	}
	return fmt.Sprintf("oauth2_token_revocation:%s:access_token", clientID)
}
