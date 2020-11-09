package cache

import (
	"context"
	"time"
)

// Cacher 缓存
type Cacher interface {
	Geter
	Seter
	Remover
}

// Geter ...
type Geter interface {
	Get(ctx context.Context, key string, v interface{}) (err error)
	GetString(ctx context.Context, key string) (value string, err error)
	GetBytes(ctx context.Context, key string) (bytes []byte, err error)
}

// Seter ...
type Seter interface {
	Set(ctx context.Context, key string, v interface{}, expiration ...time.Duration) (err error)
	SetString(ctx context.Context, key string, value string, expiration ...time.Duration) (err error)
	SetBytes(ctx context.Context, key string, bytes []byte, expiration ...time.Duration) (err error)
}

// Remover ...
type Remover interface {
	Remove(ctx context.Context, keys ...string) (err error)
	RemoveMatch(ctx context.Context, match string) (err error)
}
