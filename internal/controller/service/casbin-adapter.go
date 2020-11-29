package service

import (
	"context"
	"strings"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CasbinAdapterServer casbin适配器服务端
type CasbinAdapterServer struct {
}

// LoadPolicy 加载规则
func (ctl *CasbinAdapterServer) LoadPolicy(ctx context.Context, req *proto.LoadPolicyRequest) (resp *proto.LoadPolicyResponse, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	resourceID := md.Get("resource_id")[0]

	resp = new(proto.LoadPolicyResponse)
	results, resultErr := service.Resource.LoadPolicy(ctx, model.ConvertStringToID(resourceID))
	if resultErr != nil {
		err = status.Error(codes.Unavailable, resultErr.Error())
		return
	}
	for _, result := range results {
		resp.Policys = append(resp.Policys, policyLine(result))
	}
	return
}

func policyLine(line *gormadapter.CasbinRule) (lineText string) {
	var p = []string{line.PType,
		line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}

	if line.V5 != "" {
		lineText = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lineText = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lineText = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lineText = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lineText = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lineText = strings.Join(p[:2], ", ")
	}
	return
}
