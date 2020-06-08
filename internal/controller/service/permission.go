package service

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/protobuf/errors"
	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// PermissionService 许可服务
type PermissionService struct {
}

func (ctl *PermissionService) checkOAuth2(oauth2Client *proto.OAuth2Client) (err *errors.BusinessError) {
	if oauth2Client == nil {
		err = &errors.BusinessError{
			Code: 0,
			Msg:  "request oauth2_client is nil",
		}
		return
	}

	if oauth2Client.GetId() == "" {
		err = &errors.BusinessError{
			Code: 0,
			Msg:  "request oauth2_client_id is empty",
		}
		return
	}

	if oauth2Client.GetSecret() == "" {
		err = &errors.BusinessError{
			Code: 0,
			Msg:  "request oauth2_client_secret is empty",
		}
		return
	}

	var (
		client    *model.OAuth2Client
		clientErr error
	)
	client, clientErr = service.OAuth2.GetClient(oauth2Client.Id)
	if clientErr != nil {
		logger.Errorf("service.OAuth2.GetClient Error: %s", clientErr)
		err = &errors.BusinessError{
			Code: 0,
			Msg:  oauth2.ErrUnauthorizedClient.Error(),
		}
		return
	}
	if convert.ToString(client.ClientID) != oauth2Client.Id || client.ClientSecret != oauth2Client.Secret {
		err = &errors.BusinessError{
			Code: 0,
			Msg:  oauth2.ErrUnauthorizedClient.Error(),
		}
		return
	}
	return
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
	if res.Err = ctl.checkOAuth2(req.Client); res.Err != nil {
		return
	}
	var (
		claims    *oauth2.JwtClaims
		claimsErr error
	)
	claims, claimsErr = oauth2.ParseJwtClaimsToken(req.GetToken(), []byte(viper.GetString("jwt.secret")))
	if claimsErr != nil {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  fmt.Sprintf("token is denied error: %s", claimsErr),
		}
		return
	}
	if claims.VerifyAudience([]string{req.GetClient().Id}, false) {
		res.Err = &errors.BusinessError{
			Code: 0,
			Msg:  fmt.Sprintf("token claims audience not equal to %s", req.GetClient().Id),
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
		userInfo, userInfoErr := service.User.GetInfoOneByUserID(convert.ToUint64(res.UserInfo.OpenId))
		if userInfoErr == nil && userInfo != nil {
			res.UserInfo.NickName = userInfo.Nickname
			res.UserInfo.AvatarUrl = userInfo.Picture
			res.UserInfo.Gender = uint32(userInfo.Gender)
		}
	}
	return
}
