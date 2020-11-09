package oauth2

import "time"

const (
	contentTypeJSON = "application/json"
	// AccessTokenExpire ...
	AccessTokenExpire = time.Second * 3600
	// RefreshTokenExpire ...
	RefreshTokenExpire = AccessTokenExpire / 2
	// TokenTypeBearer ...
	TokenTypeBearer = "Bearer"
	// ScopeRefreshToken ...
	ScopeRefreshToken = "refresh_token"
	// DefaultJwtIssuer ...
	DefaultJwtIssuer = "github.com/nilorg/oauth2"
)
