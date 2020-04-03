package gateway

import (
	"context"

	"github.com/nilorg/naas/pkg/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func Service(mux *runtime.ServeMux) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err = proto.RegisterAccountHandlerServer(ctx, mux, new(AccountService)); err != nil {
		return
	}
	return nil
}
