package logger

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisHook ...
type RedisHook struct {
	Channel   string
	redis     *redis.Client
	formatter logrus.Formatter
	LogLevels []logrus.Level
}

// NewRedisHook ...
func NewRedisHook(redis *redis.Client, fields logrus.Fields, channel string) logrus.Hook {
	return &RedisHook{
		redis:     redis,
		formatter: DefaultFormatter(fields),
		Channel:   channel,
		LogLevels: logrus.AllLevels,
	}
}

// Fire ...
func (h *RedisHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	if e.Context == nil {
		e.Context = context.Background()
	}
	err = h.redis.Publish(e.Context, h.Channel, dataBytes).Err()
	return err
}

// Levels returns all logrus levels.
func (h *RedisHook) Levels() []logrus.Level {
	return h.LogLevels
}
