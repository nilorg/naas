package service

import (
	"context"

	"github.com/nilorg/naas/pkg/proto"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterGrpcGateway 注册Grpc网关
func RegisterGrpcGateway(mux *runtime.ServeMux) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err = proto.RegisterPermissionHandlerServer(ctx, mux, new(PermissionServer)); err != nil {
		return
	}
	if err = proto.RegisterCasbinAdapterHandlerServer(ctx, mux, new(CasbinAdapterServer)); err != nil {
		return
	}
	return nil
}

// RegisterGrpc 注册Grpc
func RegisterGrpc(server *grpc.Server) {
	proto.RegisterPermissionServer(server, new(PermissionServer))
	proto.RegisterCasbinAdapterServer(server, new(CasbinAdapterServer))
	proto.RegisterResourceServer(server, new(ResourceServer))
}
