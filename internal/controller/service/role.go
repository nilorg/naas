package service

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// RoleServer can be embedded to have forward compatible implementations.
type RoleServer struct {
}

// ListActionByRoles ...
func (*RoleServer) ListRoleByOpenID(ctx context.Context, req *proto.ListRoleByOpenIDRequest) (resp *proto.ListRoleByOpenIDResponse, err error) {
	var roles []*model.UserRole
	roles, err = service.Role.GetAllRoleByUserID(ctx, model.ConvertStringToID(req.OpenId))
	if err != nil {
		return
	}
	resp = new(proto.ListRoleByOpenIDResponse)
	for _, role := range roles {
		resp.Roles = append(resp.Roles, &proto.UserRole{
			Id:             model.ConvertIDToUint64(role.ID),
			UserId:         uint64(role.UserID),
			RoleCode:       string(role.RoleCode),
			OrganizationId: uint64(role.OrganizationID),
		})
	}
	return
}
