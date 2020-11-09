package storage

import (
	"context"
)

type contentTypeKey struct{}

// NewContentTypeContext ...
func NewContentTypeContext(ctx context.Context, contentType string) context.Context {
	return context.WithValue(ctx, contentTypeKey{}, contentType)
}

// FromContentTypeContext ...
func FromContentTypeContext(ctx context.Context) (contentType string, ok bool) {
	contentType, ok = ctx.Value(contentTypeKey{}).(string)
	return
}

// ---

type bucketNameKey struct{}

// NewBucketNameContext ...
func NewBucketNameContext(ctx context.Context, bucketName string) context.Context {
	return context.WithValue(ctx, bucketNameKey{}, bucketName)
}

// FromBucketNameContext ...
func FromBucketNameContext(ctx context.Context) (bucketName string, ok bool) {
	bucketName, ok = ctx.Value(bucketNameKey{}).(string)
	return
}
