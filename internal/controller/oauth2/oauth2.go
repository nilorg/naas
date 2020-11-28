package oauth2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nilorg/sdk/convert"

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
	"github.com/nilorg/pkg/slice"
	sdkStrings "github.com/nilorg/sdk/strings"
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
		oauth2.ServerDeviceVerificationURI("/oauth2/device/activate"),
		oauth2.ServerIntrospectEndpointEnabled(viper.GetBool("server.oauth2.introspection_endpoint_enabled")),
		oauth2.ServerTokenRevocationEnabled(viper.GetBool("server.oauth2.revocation_endpoint_enabled")),
	)
	// TODO: 这个方法需要优化
	oauth2Server.VerifyClientID = func(clientID string) (err error) {
		var client *model.OAuth2Client
		client, err = service.OAuth2.GetClient(contexts.WithContext(context.Background()), model.ConvertStringToID(clientID))
		if err != nil {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		if convert.ToString(client.ClientID) != clientID {
			err = oauth2.ErrUnauthorizedClient
		}
		return
	}
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
			logrus.Errorf("store.RedisClient.Set Error: %s", err)
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
			logrus.Errorf("store.RedisClient.Get Error: %s", err)
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
	oauth2Server.GenerateDeviceAuthorization = func(issuer, verificationURI, clientID string, scope []string) (resp *oauth2.DeviceAuthorizationResponse, err error) {
		ctx := context.Background()

		deviceCode := oauth2.RandomDeviceCode()
		deviceCodeKey := key.WrapOAuth2DeviceCode(deviceCode)
		userCode := oauth2.RandomUserCode()
		userCodekey := key.WrapOAuth2UserCode(userCode)
		expires := time.Now().Add(time.Minute)
		err = store.RedisClient.HSet(
			ctx, userCodekey,
			"device_code", deviceCode,
			"client_id", clientID,
			"scope", strings.Join(scope, ","),
		).Err()
		if err != nil {
			err = oauth2.ErrServerError
			return
		}
		err = store.RedisClient.ExpireAt(ctx, userCodekey, expires).Err()
		if err != nil {
			logrus.Errorln(err)
			err = oauth2.ErrServerError
			return
		}
		err = store.RedisClient.HSet(
			ctx, deviceCodeKey,
			"user_code", userCode,
			"client_id", clientID,
			"scope", strings.Join(scope, ","),
			"status", "0",
			"open_id", "",
		).Err()
		if err != nil {
			err = oauth2.ErrServerError
			return
		}
		err = store.RedisClient.ExpireAt(ctx, deviceCodeKey, expires).Err()
		if err != nil {
			logrus.Errorln(err)
			err = oauth2.ErrServerError
			return
		}
		resp = new(oauth2.DeviceAuthorizationResponse)
		// 有效的时间长度（以秒为单位）。
		// 如果在那个时候用户没有完成授权流程，并且您的设备也没有轮询来检索有关用户决定的信息，则您可能需要从步骤1重新开始此过程。
		resp.ExpiresIn = expires.Unix()
		// 您的设备应在两次轮询请求之间等待的时间长度（以秒为单位）。
		// 例如，如果值为5，则您的设备应每五秒钟向NAAS授权服务器发送一次轮询请求。
		resp.Interval = 5
		// NAAS唯一分配的值，用于标识运行请求授权的应用的设备。
		// 用户将从具有更丰富输入功能的另一台设备授权该设备。
		// 例如，用户可能使用笔记本电脑或移动电话来授权在电视上运行的应用。
		// 在这种情况下，标识电视。 device_code 此代码可让运行该应用的设备安全地确定用户是否已授予访问权限。
		resp.DeviceCode = deviceCode
		// 区分大小写的值，用于向NAAS标识应用程序请求访问的范围。
		// 您的用户界面将指示用户在具有更丰富输入功能的单独设备上输入此值。
		// 然后，当提示用户授予对您的应用程序的访问权限时，NAAS使用该值显示正确的范围集。
		resp.UserCode = userCode
		// 用户必须在单独的设备上导航到的URL，以输入和授予或拒绝对您的应用程序的访问。
		// 您的用户界面还将显示此值。user_code
		resp.VerificationURI = issuer + verificationURI
		resp.VerificationURIComplete = fmt.Sprintf("%s?user_code=%s", issuer+verificationURI, userCode)
		return
	}
	oauth2Server.VerifyDeviceCode = func(deviceCode, clientID string) (value *oauth2.DeviceCodeValue, err error) {
		ctx := context.Background()
		deviceCodeKey := key.WrapOAuth2DeviceCode(deviceCode)
		if store.RedisClient.Exists(ctx, deviceCodeKey).Val() != 1 {
			err = oauth2.ErrAccessDenied
			return
		}
		var deviceValue map[string]string
		if deviceValue, err = store.RedisClient.HGetAll(ctx, deviceCodeKey).Result(); err != nil {
			logrus.Errorln(err)
			err = oauth2.ErrServerError
			return
		}
		if deviceValue["client_id"] != clientID {
			err = oauth2.ErrInvalidClient
			return
		}
		// 检查用户是否同意授权
		if deviceValue["status"] != "1" {
			err = oauth2.ErrAuthorizationPending
		} else {
			// 删除验证数据
			value = new(oauth2.DeviceCodeValue)
			// value.ClientID = clientID
			// value.DeviceCode = deviceCode
			// value.UserCode = userCode
			value.Scope = RemoveRepeat(sdkStrings.Split(deviceValue["scope"], ","))
			value.OpenID = deviceValue["open_id"]

			store.RedisClient.Del(ctx, deviceCodeKey)
		}
		return
	}
	oauth2Server.VerifyIntrospectionToken = func(token, clientID string, tokenTypeHint ...string) (resp *oauth2.IntrospectionResponse, err error) {
		logrus.Debugf("oauth2Server.VerifyIntrospectionToken....")
		key := tokenRevocationKey(token, clientID, tokenTypeHint...)
		var exsit bool
		exsit, err = store.RedisClient.HExists(context.Background(), key, token).Result()
		if err != nil {
			logrus.Errorf("store.RedisClient.HExists: %s", err)
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
			logrus.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			err = oauth2.ErrServerError
			return
		}
		if !tokenClaims.VerifyAudience([]string{clientID}, false) {
			logrus.Debugf("tokenClaims.VerifyAudience.....false")
			err = oauth2.ErrInvalidClient
		}
		resp = new(oauth2.IntrospectionResponse)
		resp.Active = true
		if verr := tokenClaims.Valid(); verr != nil {
			logrus.Debugf("tokenClaims.Valid: %s", verr)
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
			logrus.Errorf("oauth2.ParseJwtClaimsToken: %s", err)
			return
		}
		exp := time.Unix(tokenClaims.ExpiresAt, 0)
		err = store.RedisClient.HSet(context.Background(), key, token, time.Now().Sub(exp)).Err()
		if err != nil {
			logrus.Errorf("store.RedisClient.Set: %s", err)
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
