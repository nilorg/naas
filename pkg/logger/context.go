package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nilorg/sdk/log"
)

const (
	// TraceIDKey 跟踪ID
	TraceIDKey = "trace_id"
	// SpanIDKey 请求ID
	SpanIDKey = "span_id"
	// UserIDKey 用户ID
	UserIDKey = "user_id"
)

// WithGinContext ...
func WithGinContext(ctx *gin.Context) context.Context {
	parent := context.Background()
	if traceID := ctx.GetString("X-Trace-Id"); traceID != "" {
		parent = log.NewTraceIDContext(parent, traceID)
	} else {
		parent = log.NewTraceIDContext(parent, uuid.New().String())
	}
	if spanID := ctx.GetString("X-Span-Id"); spanID != "" {
		parent = log.NewSpanIDContext(parent, spanID)
	} else {
		parent = log.NewSpanIDContext(parent, "0")
	}
	if uid := ctx.GetString("X-User-Id"); uid != "" {
		parent = log.NewUserIDContext(parent, uid)
	}
	return parent
}
