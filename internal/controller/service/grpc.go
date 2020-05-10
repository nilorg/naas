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
	if err = proto.RegisterAccountHandlerServer(ctx, mux, new(AccountService)); err != nil {
		return
	}
	return nil
}

// RegisterGrpc 注册Grpc
func RegisterGrpc(server *grpc.Server) {
	proto.RegisterAccountServer(server, new(AccountService))
}
