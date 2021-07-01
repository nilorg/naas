package wechat

import (
	"github.com/go-redis/redis/v8"
)

// ClientFromRedisOptions 可选参数列表
type ClientFromRedisOptions struct {
	RedisClient    *redis.Client
	AccessTokenKey string
	JsAPITicketKey string
}

// ClientFromRedisOption 为可选参数赋值的函数
type ClientFromRedisOption func(*ClientFromRedisOptions)

const (
	// RedisAccessTokenKey ...
	RedisAccessTokenKey = "nilorg:go-wechat:access_token"
	// RedisJsAPITicketKey ...
	RedisJsAPITicketKey = "nilorg:go-wechat:js_api_ticket"
)

// NewClientFromRedisOptions 创建可选参数
func NewClientFromRedisOptions(opts ...ClientFromRedisOption) ClientFromRedisOptions {
	opt := ClientFromRedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
		AccessTokenKey: RedisAccessTokenKey,
		JsAPITicketKey: RedisJsAPITicketKey,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// ClientFromRedisOptionRedisClient ...
func ClientFromRedisOptionRedisClient(client *redis.Client) ClientFromRedisOption {
	return func(o *ClientFromRedisOptions) {
		o.RedisClient = client
	}
}

// ClientFromRedisOptionAccessTokenKey ...
func ClientFromRedisOptionAccessTokenKey(accessTokenKey string) ClientFromRedisOption {
	return func(o *ClientFromRedisOptions) {
		o.AccessTokenKey = accessTokenKey
	}
}

// ClientFromRedisOptionJsAPITicketKey ...
func ClientFromRedisOptionJsAPITicketKey(jsAPITicketKey string) ClientFromRedisOption {
	return func(o *ClientFromRedisOptions) {
		o.JsAPITicketKey = jsAPITicketKey
	}
}
