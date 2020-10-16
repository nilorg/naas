package store

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/nilorg/naas/internal/model"

	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/pkg/random"
	"github.com/nilorg/sdk/cache"
	sdkCache "github.com/nilorg/sdk/cache"
	"gorm.io/gorm"
)

var (
	// ErrContextNotFoundCache 上下文不存在Cache错误
	ErrContextNotFoundCache = errors.New("上下文中没有获取到Cache")
)

// NewCacheContext 创建缓存对象到上下文中
func NewCacheContext(ctx context.Context, cache sdkCache.Cacher) context.Context {
	return sdkCache.NewCacheContext(ctx, cache)
}

// FromCacheContext 从缓存上下文中获取缓存对象
func FromCacheContext(ctx context.Context) (sdkCache.Cacher, error) {
	cache, ok := sdkCache.FromCacheContext(ctx)
	if !ok {
		return nil, ErrContextNotFoundCache
	}
	return cache, nil
}

type skipCache struct{}

// NewSkipCacheContext 创建跳过缓存到上下文
func NewSkipCacheContext(ctx context.Context, skip ...bool) context.Context {
	s := true
	if len(skip) > 0 {
		s = skip[0]
	}
	return context.WithValue(ctx, skipCache{}, s)
}

// FromSkipCacheContext 从上下文中获取跳过缓存变量
func FromSkipCacheContext(ctx context.Context) (skip bool) {
	var ok bool
	skip, ok = ctx.Value(skipCache{}).(bool)
	if !ok {
		skip = false
	}
	return
}

// ScanByCacheID ...
func ScanByCacheID(ctx context.Context, cacheKey string, table interface{}, query interface{}, args ...interface{}) (items []*model.CacheIDPrimaryKey, err error) {
	err = scanByCache(ctx, cacheKey, table, &items, query, args...)
	return
}

// ScanByCacheCode ...
func ScanByCacheCode(ctx context.Context, cacheKey string, table interface{}, query interface{}, args ...interface{}) (items []*model.CacheCodePrimaryKey, err error) {
	err = scanByCache(ctx, cacheKey, table, &items, query, args...)
	return
}

// scanByCache ...
func scanByCache(ctx context.Context, cacheKey string, table interface{}, values interface{}, query interface{}, args ...interface{}) (err error) {
	var (
		gdb   *gorm.DB
		cache cache.Cacher
	)
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	cache, err = FromCacheContext(ctx)
	if err != nil {
		return
	}
	err = cache.Get(ctx, cacheKey, values)
	if err != nil {
		if err == redis.Nil {
			if err = gdb.Model(table).Where(query, args...).Scan(values).Error; err != nil {
				return
			}
			if err = cache.Set(ctx, cacheKey, values, random.TimeDuration(300, 600)); err != nil {
				return
			}
		} else {
			return
		}
	}
	return
}
