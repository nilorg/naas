package gateway

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
	user, userErr := service.User.GetOneByID(openID)
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
