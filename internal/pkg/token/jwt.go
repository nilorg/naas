package token

import (
	"time"

	"github.com/nilorg/oauth2"
	"github.com/nilorg/sdk/strings"
)

// NewGenerateAccessToken 创建默认生成AccessToken方法
func NewGenerateAccessToken(key interface{}, idTokenEnabled bool) oauth2.GenerateAccessTokenFunc {
	return func(issuer, clientID, scope, openID string) (token *oauth2.TokenResponse, err error) {
		accessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, scope, openID)
		var tokenStr string
		tokenStr, err = oauth2.NewJwtToken(accessJwtClaims, "RS256", key)
		if err != nil {
			err = oauth2.ErrServerError
			return
		}

		refreshAccessJwtClaims := oauth2.NewJwtClaims(issuer, clientID, oauth2.ScopeRefreshToken, "")
		refreshAccessJwtClaims.ID = tokenStr
		var refreshTokenStr string
		refreshTokenStr, err = oauth2.NewJwtToken(accessJwtClaims, "RS256", key)
		if err != nil {
			err = oauth2.ErrServerError
			return
		}
		currTime := time.Now()
		token = &oauth2.TokenResponse{
			AccessToken:  tokenStr,
			TokenType:    oauth2.TokenTypeBearer,
			ExpiresIn:    accessJwtClaims.ExpiresAt,
			RefreshToken: refreshTokenStr,
			Scope:        scope,
		}
		idTokenFlag := false
		for _, s := range strings.Split(scope, " ") {
			if s == "openid" {
				idTokenFlag = true
				break
			}
		}
		if idTokenFlag && idTokenEnabled {
			idTokenJwtClaims := oauth2.JwtClaims{
				JwtStandardClaims: oauth2.JwtStandardClaims{
					Issuer:    issuer,
					Subject:   openID,
					IssuedAt:  currTime.Unix(),
					ExpiresAt: currTime.Add(oauth2.AccessTokenExpire).Unix(),
					Audience:  []string{clientID},
				},
				Scope: scope,
			}
			var idToken string
			idToken, err = oauth2.NewJwtClaimsToken(&idTokenJwtClaims, "RS256", key)
			if err != nil {
				err = oauth2.ErrServerError
				return
			}
			token.IDToken = idToken
		}
		return
	}
}

// NewRefreshAccessToken 创建默认刷新AccessToken方法
func NewRefreshAccessToken(key interface{}) oauth2.RefreshAccessTokenFunc {
	return func(clientID, refreshToken string) (token *oauth2.TokenResponse, err error) {
		refreshTokenClaims := &oauth2.JwtClaims{}
		refreshTokenClaims, err = oauth2.ParseJwtClaimsToken(refreshToken, key)
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
		tokenClaims, err = oauth2.ParseJwtClaimsToken(refreshTokenClaims.ID, key)
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
		return oauth2.ParseJwtClaimsToken(accessToken, key)
	}
}
