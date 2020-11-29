package proto

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// SetResourceAuth 设置资源服务器身份验证
func SetResourceAuth(ctx context.Context, resourceID, resourceSecret string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}
	buffer := bytes.NewBufferString("basic ")
	buffer.WriteString(
		base64.StdEncoding.EncodeToString(
			[]byte(
				fmt.Sprintf("%s:%s", resourceID, resourceSecret),
			),
		),
	)
	md.Set("authorization", buffer.String())
	return metadata.NewOutgoingContext(ctx, md)
}
