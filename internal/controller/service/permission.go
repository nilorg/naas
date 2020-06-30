package service

import (
	"context"
	"fmt"
	"github.com/nilorg/naas/internal/module/casbin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spf13/viper"

	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// PermissionService 许可服务
type PermissionService struct {
}

func (ctl *PermissionService) checkReSource(resource *proto.Resource) (err error) {
	if resource == nil {
		err = status.Error(codes.InvalidArgument, "request resource is nil")
		return
	}

	if resource.GetId() == "" {
		err = status.Error(codes.InvalidArgument, "request resource_id is empty")
		return
	}

	if resource.GetSecret() == "" {
		err = status.Error(codes.InvalidArgument, "request resource_secret is empty")
		return
	}

	var (
		rs    *model.Resource
		rsErr error
	)
	rs, rsErr = service.Resource.Get(convert.ToUint64(resource.Id))
	if rsErr != nil {
		logger.Errorf("service.Resource.GetClient Error: %s", rsErr)
		err = status.Error(codes.Unavailable, oauth2.ErrUnauthorizedClient.Error())
		return
	}
	if convert.ToString(rs.ID) != resource.Id || rs.Secret != resource.Secret {
		err = status.Error(codes.Unavailable, oauth2.ErrUnauthorizedClient.Error())
	}
	return
}

// VerificationHttpRouter 验证Http路由
func (ctl *PermissionService) VerificationHttpRouter(ctx context.Context, req *proto.VerificationHttpRouterRequest) (res *proto.VerificationHttpRouterResponse, err error) {
	var resToken *proto.VerificationTokenResponse
	resToken, err = ctl.VerificationToken(ctx, &proto.VerificationTokenRequest{
		Resource:       req.Resource,
		Token:          req.Token,
		ReturnUserInfo: false,
	})
	if err != nil {
		return
	}
	res = new(proto.VerificationHttpRouterResponse)
	if !resToken.Allow {
		res.Allow = false
	}
	res.UserInfo = &proto.VerificationHttpRouterResponse_UserInfo{
		OpenId:    res.UserInfo.OpenId,
		Username:  res.UserInfo.Username,
		NickName:  res.UserInfo.NickName,
		AvatarUrl: res.UserInfo.AvatarUrl,
		Gender:    res.UserInfo.Gender,
	}

	roles, _ := service.Role.GetAllRoleByUserID(convert.ToUint64(res.UserInfo.OpenId))
	for _, role := range roles {
		sub := fmt.Sprintf("role:%s", role.RoleCode)                 // 希望访问资源的用户
		dom := fmt.Sprintf("resource:%s:web_route", req.Resource.Id) // 域/域租户,这里以资源为单位
		obj := req.Path                                              // 要访问的资源
		act := req.Method                                            // 用户对资源执行的操作
		if check, checkErr := casbin.Enforcer.Enforce(sub, dom, obj, act); checkErr != nil && check {
			res.Allow = true
		}
	}
	return
}

// VerificationToken 验证Token
func (ctl *PermissionService) VerificationToken(_ context.Context, req *proto.VerificationTokenRequest) (res *proto.VerificationTokenResponse, err error) {
	res = new(proto.VerificationTokenResponse)
	if req.GetToken() == "" {
		err = status.Error(codes.InvalidArgument, "request token is empty")
		return
	}
	err = ctl.checkReSource(req.Resource)
	if err != nil {
		return
	}
	var (
		claims    *oauth2.JwtClaims
		claimsErr error
	)
	claims, claimsErr = oauth2.ParseJwtClaimsToken(req.GetToken(), []byte(viper.GetString("jwt.secret")))
	if claimsErr != nil {
		err = status.Error(codes.Internal, fmt.Sprintf("token is denied error: %s", claimsErr))
		return
	}
	if claims.VerifyAudience([]string{req.Resource.Id}, false) {
		err = status.Error(codes.PermissionDenied, fmt.Sprintf("token claims audience not equal to %s", req.Resource.Id))
		return
	}
	openID := claims.Subject

	if req.GetReturnUserInfo() {
		user, userErr := service.User.GetOneByID(openID)
		if userErr != nil {
			err = status.Error(codes.Unavailable, userErr.Error())
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
