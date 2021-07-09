package oauth2

import (
	"net/http"
	"strings"
	"time"

	"github.com/nilorg/pkg/slice"
	sdkStrings "github.com/nilorg/sdk/strings"
)

// VerifyClientFunc 验证客户端委托
type VerifyClientFunc func(basic *ClientBasic) (err error)

// VerifyClientIDFunc 验证客户端ID委托
type VerifyClientIDFunc func(clientID string) (err error)

// VerifyRedirectURIFunc 验证RedirectURI委托
type VerifyRedirectURIFunc func(clientID, redirectURI string) (err error)

// GenerateCodeFunc 生成Code委托
type GenerateCodeFunc func(clientID, openID, redirectURI string, scope []string) (code string, err error)

// VerifyCodeFunc 验证Code委托
type VerifyCodeFunc func(code, clientID, redirectURI string) (value *CodeValue, err error)

// VerifyPasswordFunc 验证账号密码委托
type VerifyPasswordFunc func(username, password string) (openID string, err error)

// VerifyScopeFunc 验证范围委托
type VerifyScopeFunc func(scope []string, clientID string) (err error)

// GenerateAccessTokenFunc 生成AccessToken委托
type GenerateAccessTokenFunc func(issuer, clientID, scope, openID string, code *CodeValue) (token *TokenResponse, err error)

// GenerateDeviceAuthorizationFunc 生成设备授权
type GenerateDeviceAuthorizationFunc func(issuer, verificationURI, clientID string, scope []string) (resp *DeviceAuthorizationResponse, err error)

// ParseAccessTokenFunc 解析AccessToken为JwtClaims委托
type ParseAccessTokenFunc func(accessToken string) (claims *JwtClaims, err error)

// RefreshAccessTokenFunc 刷新AccessToken委托
type RefreshAccessTokenFunc func(clientID, refreshToken string) (token *TokenResponse, err error)

// VerifyDeviceCodeFunc 验证DeviceCode委托
type VerifyDeviceCodeFunc func(deviceCode, clientID string) (value *DeviceCodeValue, err error)

// VerifyIntrospectionTokenFunc 验证IntrospectionToken委托
type VerifyIntrospectionTokenFunc func(token, clientID string, tokenTypeHint ...string) (resp *IntrospectionResponse, err error)

// TokenRevocationFunc Token撤销委托
// https://tools.ietf.org/html/rfc7009#section-2.2
type TokenRevocationFunc func(token, clientID string, tokenTypeHint ...string)

// CustomGrantTypeAuthenticationFunc 自定义GrantType身份验证委托
type CustomGrantTypeAuthenticationFunc func(client *ClientBasic, req *http.Request) (openID string, err error)

// VerifyGrantTypeFunc 验证授权类型委托
type VerifyGrantTypeFunc func(clientID, grantType string) (err error)

// NewDefaultGenerateAccessToken 创建默认生成AccessToken方法
func NewDefaultGenerateAccessToken(jwtVerifyKey []byte) GenerateAccessTokenFunc {
	return func(issuer, clientID, scope, openID string, codeVlue *CodeValue) (token *TokenResponse, err error) {
		scopeSplit := sdkStrings.Split(scope, " ")
		accessJwtClaims := NewJwtClaims(issuer, clientID, scope, openID)
		if codeVlue != nil {
			if len(scopeSplit) > 0 && !slice.IsEqual(scopeSplit, codeVlue.Scope) {
				accessJwtClaims = NewJwtClaims(issuer, clientID, strings.Join(codeVlue.Scope, " "), openID)
			}
		}
		var tokenStr string
		tokenStr, err = NewHS256JwtClaimsToken(accessJwtClaims, jwtVerifyKey)
		if err != nil {
			err = ErrServerError
			return
		}

		refreshAccessJwtClaims := NewJwtClaims(issuer, clientID, ScopeRefreshToken, "")
		refreshAccessJwtClaims.ID = tokenStr
		var refreshTokenStr string
		refreshTokenStr, err = NewHS256JwtClaimsToken(refreshAccessJwtClaims, jwtVerifyKey)
		if err != nil {
			err = ErrServerError
			return
		}
		token = &TokenResponse{
			AccessToken:  tokenStr,
			TokenType:    TokenTypeBearer,
			ExpiresIn:    accessJwtClaims.ExpiresAt,
			RefreshToken: refreshTokenStr,
			Scope:        scope,
		}
		return
	}
}

// NewDefaultRefreshAccessToken 创建默认刷新AccessToken方法
func NewDefaultRefreshAccessToken(jwtVerifyKey []byte) RefreshAccessTokenFunc {
	return func(clientID, refreshToken string) (token *TokenResponse, err error) {
		var refreshTokenClaims *JwtClaims
		refreshTokenClaims, err = ParseHS256JwtClaimsToken(refreshToken, jwtVerifyKey)
		if err != nil {
			return
		}
		if refreshTokenClaims.Subject != clientID {
			err = ErrUnauthorizedClient
			return
		}
		if refreshTokenClaims.Scope != ScopeRefreshToken {
			err = ErrInvalidScope
			return
		}
		refreshTokenClaims.ExpiresAt = time.Now().Add(AccessTokenExpire).Unix()

		var tokenClaims *JwtClaims
		tokenClaims, err = ParseHS256JwtClaimsToken(refreshTokenClaims.ID, jwtVerifyKey)
		if err != nil {
			return
		}
		if tokenClaims.Subject != clientID {
			err = ErrUnauthorizedClient
			return
		}
		tokenClaims.ExpiresAt = time.Now().Add(AccessTokenExpire).Unix()

		var refreshTokenStr string
		refreshTokenStr, err = NewHS256JwtClaimsToken(refreshTokenClaims, jwtVerifyKey)
		if err != nil {
			return
		}
		var tokenStr string
		tokenStr, err = NewHS256JwtClaimsToken(tokenClaims, jwtVerifyKey)
		token = &TokenResponse{
			AccessToken:  tokenStr,
			RefreshToken: refreshTokenStr,
			TokenType:    TokenTypeBearer,
			ExpiresIn:    refreshTokenClaims.ExpiresAt,
			Scope:        tokenClaims.Scope,
		}
		return
	}
}

// NewDefaultParseAccessToken 创建默认解析AccessToken方法
func NewDefaultParseAccessToken(jwtVerifyKey []byte) ParseAccessTokenFunc {
	return func(accessToken string) (claims *JwtClaims, err error) {
		return ParseHS256JwtClaimsToken(accessToken, jwtVerifyKey)
	}
}
