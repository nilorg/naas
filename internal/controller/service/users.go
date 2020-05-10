package service

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/protobuf/errors"

	"github.com/nilorg/naas/internal/service"

	"github.com/nilorg/naas/pkg/proto"
)

// AccountService 账户服务
type AccountService struct {
}

// GetUserInfo 根据OpenID获取用户信息
func (ctl *AccountService) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (res *proto.GetCurrentResponse, err error) {
	res = new(proto.GetCurrentResponse)
	if req.OpenId == "" {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  "metadata OpenId is empty",
		}
		return
	}
	user, userErr := service.User.GetOneByID(req.OpenId)
	if userErr != nil {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  userErr.Error(),
		}
		return
	}
	var deletedAt int64 = 0
	if user.DeletedAt != nil {
		deletedAt = user.DeletedAt.Unix()
	}
	res.UserDetail = &proto.UserDetail{
		Id:        convert.ToString(user.ID),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.CreatedAt.Unix(),
		DeletedAt: deletedAt,
		Username:  user.Username,
	}
	return
}

// GetCurrent 获取当前用户
func (ctl *AccountService) GetCurrent(ctx context.Context, _ *proto.GetCurrentRequest) (res *proto.GetCurrentResponse, err error) {
	res = new(proto.GetCurrentResponse)
	openID := metautils.ExtractIncoming(ctx).Get("OpenId")
	if openID == "" {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  "metadata OpenId is empty",
		}
		return
	}
	res, err = ctl.GetUserInfo(ctx, &proto.GetUserInfoRequest{
		OpenId: openID,
	})
	return
}
