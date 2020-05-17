package token

import (
	"time"

	"github.com/nilorg/oauth2"
)

// NewGenerateAccessToken 创建默认生成AccessToken方法
func NewGenerateAccessToken(key interface{}) oauth2.GenerateAccessTokenFunc {
	return func(issuer, clientID, scope, openID string) (token *oauth2.TokenResponse, err error) {
		accessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, scope, openID)
		var tokenStr string
		tokenStr, err = oauth2.NewJwtToken(accessJwtClaims, "RS256", key)
		if err != nil {
			err = oauth2.ErrServerError
		}

		refreshAccessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, oauth2.ScopeRefreshToken, "")
		refreshAccessJwtClaims.ID = tokenStr
		var refreshTokenStr string
		refreshTokenStr, err = oauth2.NewJwtToken(accessJwtClaims, "RS256", key)
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
		refreshTokenClaims, err = oauth2.ParseJwtToken(refreshToken, key)
		if err != nil {
			return
		}
		if refreshTokenClaims.Subject != clientID {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		if refreshTokenClaims.VerifyScope(oauth2.ScopeRefreshToken, false) {
			err = oauth2.ErrInvalidScope
			return
		}
		refreshTokenClaims.ExpiresAt = time.Now().Add(oauth2.AccessTokenExpire).Unix()

		var tokenClaims *oauth2.JwtClaims
		tokenClaims, err = oauth2.ParseJwtToken(refreshTokenClaims.ID, key)
		if err != nil {
			return
		}
		if tokenClaims.Subject != clientID {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		tokenClaims.ExpiresAt = time.Now().Add(oauth2.AccessTokenExpire).Unix()

		var refreshTokenStr string
		refreshTokenStr, err = oauth2.NewJwtToken(refreshTokenClaims, "RS256", key)
		if err != nil {
			return
		}
		var tokenStr string
		tokenStr, err = oauth2.NewJwtToken(tokenClaims, "RS256", key)
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
		return oauth2.ParseJwtToken(accessToken, key)
	}
}
