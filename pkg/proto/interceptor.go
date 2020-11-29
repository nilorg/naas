package proto

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// GrpcContextHandler ...
type GrpcContextHandler func(ctx context.Context) context.Context

// UnaryServerInterceptor ...
func UnaryServerInterceptor(f GrpcContextHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = f(ctx)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor ...
func StreamServerInterceptor(f GrpcContextHandler) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = f(stream.Context())
		return handler(srv, wrapped)
	}
}
