package wechat

import (
	"context"
	"errors"
)

var (
	// ErrContextNotFoundClient 上下文不存在客户端错误
	ErrContextNotFoundClient = errors.New("上下文中没有获取到微信客户端")
)

type wechatKey struct{}

// FromContext 从上下文中获取微信客户端
func FromContext(ctx context.Context) (Clienter, error) {
	c, ok := ctx.Value(wechatKey{}).(Clienter)
	if !ok {
		return nil, ErrContextNotFoundClient
	}
	return c, nil
}

// NewContext 创建微信客户端上下文
func NewContext(ctx context.Context, c Clienter) context.Context {
	return context.WithValue(ctx, wechatKey{}, c)
}
