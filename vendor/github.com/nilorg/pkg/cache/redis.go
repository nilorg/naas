package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache redis缓存
type RedisCache struct {
	redisClient *redis.Client
	Prefix      string
}

// NewRedisCache ...
func NewRedisCache(redisClient *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		redisClient: redisClient,
		Prefix:      prefix,
	}
}

// Get ...
func (r *RedisCache) Get(ctx context.Context, key string, v interface{}) (err error) {
	key = r.formatKey(key)
	var bytes []byte
	bytes, err = r.getBytes(ctx, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, v)
	return
}

// GetString ...
func (r *RedisCache) GetString(ctx context.Context, key string) (value string, err error) {
	key = r.formatKey(key)
	var bytes []byte
	bytes, err = r.getBytes(ctx, key)
	if err != nil {
		return
	}
	value = string(bytes)
	return
}

// GetBytes ...
func (r *RedisCache) GetBytes(ctx context.Context, key string) (bytes []byte, err error) {
	key = r.formatKey(key)
	return r.getBytes(ctx, key)
}

// Set ...
func (r *RedisCache) Set(ctx context.Context, key string, v interface{}, expiration ...time.Duration) (err error) {
	var bytes []byte
	bytes, err = json.Marshal(v)
	if err != nil {
		return
	}
	key = r.formatKey(key)
	err = r.setBytes(ctx, key, bytes, expiration...)
	return
}

// SetString ...
func (r *RedisCache) SetString(ctx context.Context, key string, value string, expiration ...time.Duration) (err error) {
	key = r.formatKey(key)
	return r.setBytes(ctx, key, []byte(value), expiration...)
}

// SetBytes ...
func (r *RedisCache) SetBytes(ctx context.Context, key string, bytes []byte, expiration ...time.Duration) (err error) {
	key = r.formatKey(key)
	return r.setBytes(ctx, key, bytes, expiration...)
}

// Remove ...
func (r *RedisCache) Remove(ctx context.Context, keys ...string) (err error) {
	l := len(keys)
	for i := 0; i < l; i++ {
		keys[i] = r.formatKey(keys[i])
	}
	return r.remove(ctx, keys...)
}

func (r *RedisCache) getBytes(ctx context.Context, key string) (bytes []byte, err error) {
	return r.redisClient.Get(ctx, key).Bytes()
}
func (r *RedisCache) setBytes(ctx context.Context, key string, bytes []byte, expiration ...time.Duration) (err error) {
	if len(expiration) > 0 {
		err = r.redisClient.Set(ctx, key, bytes, expiration[0]).Err()
	} else {
		err = r.redisClient.Set(ctx, key, bytes, 0).Err()
	}
	return
}
func (r *RedisCache) remove(ctx context.Context, keys ...string) (err error) {
	err = r.redisClient.Del(ctx, keys...).Err()
	return
}

// RemoveMatch ...
func (r *RedisCache) RemoveMatch(ctx context.Context, match string) (err error) {
	match = r.formatKey(match)
	return r.removeMatch(ctx, match)
}

func (r *RedisCache) removeMatch(ctx context.Context, match string) (err error) {
	iter := r.redisClient.Scan(ctx, 0, match, 0).Iterator()
	for iter.Next(ctx) {
		err = r.redisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			return
		}
	}
	return
}

func (r *RedisCache) formatKey(key string) string {
	return r.Prefix + key
}
