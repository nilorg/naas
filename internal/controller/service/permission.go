package service

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/nilorg/oauth2"
	"github.com/nilorg/protobuf/errors"
	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// PermissionService 许可服务
type PermissionService struct {
}

// VerificationHttpRouter 验证Http路由
func (ctl *PermissionService) VerificationHttpRouter(ctx context.Context, req *proto.VerificationHttpRouterRequest) (res *proto.VerificationHttpRouterResponse, err error) {
	res = new(proto.VerificationHttpRouterResponse)
	res.Err = &errors.BusinessError{
		Code: 0,
		Msg:  "接口未开发",
	}
	return
}

// VerificationToken 验证Token
func (ctl *PermissionService) VerificationToken(ctx context.Context, req *proto.VerificationTokenRequest) (res *proto.VerificationTokenResponse, err error) {
	res = new(proto.VerificationTokenResponse)
	if req.GetToken() == "" {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  "request tokne is empty",
		}
		return
	}
	if req.GetClientId() == "" {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  "request client_id is empty",
		}
		return
	}
	var (
		claims    *oauth2.JwtClaims
		claimsErr error
	)
	claims, claimsErr = oauth2.ParseAccessToken(req.GetToken(), []byte(viper.GetString("jwt.secret")))
	if claimsErr != nil {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  fmt.Sprintf("token is denied error: %s", claimsErr),
		}
		return
	}
	if claims.Audience != req.GetClientId() {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  fmt.Sprintf("token claims audience not equal to %s", req.GetClientId()),
		}
		return
	}
	openID := claims.Subject

	if req.GetReturnUserInfo() {

		user, userErr := service.User.GetOneByID(openID)
		if userErr != nil {
			res.Err = &errors.BusinessError{
				Code: 0,
				Msg:  userErr.Error(),
			}
			return
		}
		res.UserInfo = &proto.VerificationTokenResponse_UserInfo{
			OpenId:   convert.ToString(user.ID),
			Username: user.Username,
		}
		userInfo, userInfoErr := service.User.GetInfoOneByUserID(res.UserInfo.OpenId)
		if userInfoErr == nil && userInfo != nil {
			res.UserInfo.NickName = userInfo.NickName
			res.UserInfo.AvatarUrl = userInfo.AvatarURL
			res.UserInfo.Gender = uint32(userInfo.Gender)
		}
	}
	return
}
