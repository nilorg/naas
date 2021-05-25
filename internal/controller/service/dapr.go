package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/dapr/go-sdk/service/common"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// RegisterDapr 注册Dapr
func RegisterDapr(service common.Service) (err error) {
	service.AddServiceInvocationHandler("/roles/by/open_id", Roles)
	return
}

func Roles(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if strings.EqualFold(in.Verb, "get") {
		openID := model.ConvertStringToID(string(in.Data))
		roles, _ := service.Role.GetAllRoleByUserID(ctx, openID)
		out, err = proto.EncodeValue(roles)
	} else {
		err = fmt.Errorf("未实现：%s", in.Verb)
	}
	return
}
