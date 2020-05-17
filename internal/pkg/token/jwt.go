package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nilorg/oauth2"
)

func newJwtToken(claims jwt.Claims, key interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

// NewAccessToken ...
func NewAccessToken(claims *oauth2.JwtClaims, key interface{}) (string, error) {
	return newJwtToken(claims, key)
}

// ParseAccessToken ...
func ParseAccessToken(accessToken string, key interface{}) (claims *oauth2.JwtClaims, err error) {
	var token *jwt.Token
	token, err = jwt.ParseWithClaims(accessToken, &oauth2.JwtClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return key, nil
	})
	if token != nil {
		var ok bool
		if claims, ok = token.Claims.(*oauth2.JwtClaims); ok {
			return claims, nil
		}
	}
	return
}

// NewGenerateAccessToken 创建默认生成AccessToken方法
func NewGenerateAccessToken(key interface{}) oauth2.GenerateAccessTokenFunc {
	return func(issuer, clientID, scope, openID string) (token *oauth2.TokenResponse, err error) {
		accessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, scope, openID)
		var tokenStr string
		tokenStr, err = NewAccessToken(accessJwtClaims, key)
		if err != nil {
			err = oauth2.ErrServerError
		}

		refreshAccessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, oauth2.ScopeRefreshToken, "")
		refreshAccessJwtClaims.Id = tokenStr
		var refreshTokenStr string
		refreshTokenStr, err = newJwtToken(accessJwtClaims, key)
		if err != nil {
			err = oauth2.ErrServerError
		}
		token = &oauth2.TokenResponse{
			AccessToken:  tokenStr,
			TokenType:    oauth2.TokenTypeBearer,
			ExpiresIn:    accessJwtClaims.ExpiresAt,
			RefreshToken: refreshTokenStr,
			Scope:        scope,
		}
		return
	}
}

// NewRefreshAccessToken 创建默认刷新AccessToken方法
func NewRefreshAccessToken(key interface{}) oauth2.RefreshAccessTokenFunc {
	return func(clientID, refreshToken string) (token *oauth2.TokenResponse, err error) {
		refreshTokenClaims := &oauth2.JwtClaims{}
		refreshTokenClaims, err = ParseAccessToken(refreshToken, key)
		if err != nil {
			return
		}
		if refreshTokenClaims.Subject != clientID {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		if refreshTokenClaims.Scope != oauth2.ScopeRefreshToken {
			err = oauth2.ErrInvalidScope
			return
		}
		refreshTokenClaims.ExpiresAt = time.Now().Add(oauth2.AccessTokenExpire).Unix()

		var tokenClaims *oauth2.JwtClaims
		tokenClaims, err = ParseAccessToken(refreshTokenClaims.Id, key)
		if err != nil {
			return
		}
		if tokenClaims.Subject != clientID {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		tokenClaims.ExpiresAt = time.Now().Add(oauth2.AccessTokenExpire).Unix()

		var refreshTokenStr string
		refreshTokenStr, err = NewAccessToken(refreshTokenClaims, key)
		if err != nil {
			return
		}
		var tokenStr string
		tokenStr, err = NewAccessToken(tokenClaims, key)
		token = &oauth2.TokenResponse{
			AccessToken:  tokenStr,
			RefreshToken: refreshTokenStr,
			TokenType:    oauth2.TokenTypeBearer,
			ExpiresIn:    refreshTokenClaims.ExpiresAt,
			Scope:        tokenClaims.Scope,
		}
		return
	}
}

// NewParseAccessToken 创建默认解析AccessToken方法
func NewParseAccessToken(key interface{}) oauth2.ParseAccessTokenFunc {
	return func(accessToken string) (claims *oauth2.JwtClaims, err error) {
		return ParseAccessToken(accessToken, key)
	}
}
