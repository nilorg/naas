package service

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
	"google.golang.org/grpc/metadata"
)

// ResourceServer can be embedded to have forward compatible implementations.
type ResourceServer struct {
}

// ListActionByRoles ...
func (*ResourceServer) ListActionByRoles(ctx context.Context, req *proto.ListActionByRolesRequest) (resp *proto.ListActionByRolesResponse, err error) {
	if len(req.Roles) == 0 {
		return
	}
	md, _ := metadata.FromIncomingContext(ctx)
	resourceID := md.Get("resource_id")[0]
	var list []*model.ResourceAction
	list, err = service.Resource.ListActionByResourceIDAndRoleCodes(ctx, model.ConvertStringToID(resourceID), req.Roles...)
	if err != nil {
		return
	}
	resp = new(proto.ListActionByRolesResponse)
	for _, action := range list {
		resp.Actions = append(resp.Actions, &proto.ResourceAction{
			Id:          model.ConvertIDToString(action.ID),
			Name:        action.Name,
			Group:       action.Group,
			Description: action.Description,
		})
	}
	return
}

// ListMenuByRoles ...
func (*ResourceServer) ListMenuByRoles(ctx context.Context, req *proto.ListMenuByRolesRequest) (resp *proto.ListMenuByRolesResponse, err error) {
	if len(req.Roles) == 0 {
		return
	}
	md, _ := metadata.FromIncomingContext(ctx)
	resourceID := md.Get("resource_id")[0]
	var list []*model.ResourceMenu
	list, err = service.Resource.ListMenuByResourceIDAndRoleCodes(ctx, model.ConvertStringToID(resourceID), req.Roles...)
	if err != nil {
		return
	}
	resp = new(proto.ListMenuByRolesResponse)
	for _, menu := range list {
		resp.Menus = append(resp.Menus, &proto.ResourceMenu{
			Id:           model.ConvertIDToString(menu.ID),
			Name:         menu.Name,
			Icon:         menu.Icon,
			Level:        int32(menu.Level),
			SerialNumber: int32(menu.SerialNumber),
			Leaf:         menu.Leaf,
			ParentId:     model.ConvertIDToString(menu.ParentID),
		})
	}
	return
}
