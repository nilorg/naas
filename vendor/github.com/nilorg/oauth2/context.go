package oauth2

import (
	"context"
	"errors"
)

type openIDKey struct{}

var (
	// ErrContextNotFoundOpenID 上下文不存在OpenID
	ErrContextNotFoundOpenID = errors.New("OAuth2上下文不存在OpenID")
)

// OpenIDFromContext ...
func OpenIDFromContext(ctx context.Context) (string, error) {
	openID, ok := ctx.Value(openIDKey{}).(string)
	if !ok {
		return "", ErrContextNotFoundOpenID
	}
	return openID, nil
}

// NewOpenIDContext 创建OpenID上下文
func NewOpenIDContext(ctx context.Context, openID string) context.Context {
	return context.WithValue(ctx, openIDKey{}, openID)
}
