package storage

import (
	"context"
)

type mdIncomingKey struct{}

// Metadata 元数据
type Metadata map[string]string

// Get obtains the value for a given key.
func (md Metadata) Get(k string) string {
	return md[k]
}

// Set sets the value of a given key with a slice of value.
func (md Metadata) Set(k string, val string) {
	md[k] = val
}

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

// FromIncomingContext returns the incoming metadata in ctx if it exists.
func FromIncomingContext(ctx context.Context) (md Metadata, ok bool) {
	md, ok = ctx.Value(mdIncomingKey{}).(Metadata)
	return
}
