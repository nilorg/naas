package cache

import (
	"context"
)

type cacheKey struct{}

// NewCacheContext ...
func NewCacheContext(ctx context.Context, cache Cacher) context.Context {
	return context.WithValue(ctx, cacheKey{}, cache)
}

// FromCacheContext ...
func FromCacheContext(ctx context.Context) (cache Cacher, ok bool) {
	cache, ok = ctx.Value(cacheKey{}).(Cacher)
	return
}
