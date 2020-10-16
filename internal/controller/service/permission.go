package service

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spf13/viper"

	"github.com/nilorg/oauth2"
	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
)

// PermissionService 许可服务
type PermissionService struct {
}

func checkReSource(ctx context.Context, resource *proto.Resource) (err error) {
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
	rs, rsErr = service.Resource.Get(ctx, model.ConvertStringToID(resource.Id))
	if rsErr != nil {
		logrus.Errorf("service.Resource.GetClient Error: %s", rsErr)
		err = status.Error(codes.Unavailable, oauth2.ErrUnauthorizedClient.Error())
		return
	}
	if convert.ToString(rs.ID) != resource.Id || rs.Secret != resource.Secret {
		err = status.Error(codes.Unavailable, "resource id or secret is not correct")
	}
	return
}

// VerifyHttpRoute 验证Http路由
func (ctl *PermissionService) VerifyHttpRoute(ctx context.Context, req *proto.VerifyHttpRouteRequest) (res *proto.VerifyHttpRouteResponse, err error) {
	ctx = contexts.WithContext(ctx)
	err = checkReSource(ctx, req.Resource)
	if err != nil {
		return
	}
	var openID model.ID
	openID, err = ctl.verifyToken(ctx, req.Token, req.Oauth2ClientId)
	if err != nil {
		return
	}
	res = new(proto.VerifyHttpRouteResponse)
	roles, _ := service.Role.GetAllRoleByUserID(ctx, openID)
	for _, role := range roles {
		sub := fmt.Sprintf("role:%s", role.RoleCode)                 // 希望访问资源的用户
		dom := fmt.Sprintf("resource:%s:web_route", req.Resource.Id) // 域/域租户,这里以资源为单位
		obj := req.Path                                              // 要访问的资源
		act := req.Method                                            // 用户对资源执行的操作
		check, checkErr := casbin.Enforcer.Enforce(sub, dom, obj, act)
		if checkErr != nil {
			err = status.Error(codes.Unavailable, checkErr.Error())
			return
		}
		if check {
			res.Allow = true
		}
	}
	// 返回用户信息
	if res.Allow && req.ReturnUserInfo {
		user, userErr := service.User.GetOneByID(ctx, openID)
		if userErr != nil {
			err = status.Error(codes.Unavailable, userErr.Error())
			return
		}
		res.UserInfo = &proto.VerifyHttpRouteResponse_UserInfo{
			OpenId:   convert.ToString(user.ID),
			Username: user.Username,
		}
		userInfo, userInfoErr := service.User.GetInfoOneByUserID(ctx, model.ConvertStringToID(res.UserInfo.OpenId))
		if userInfoErr == nil && userInfo != nil {
			res.UserInfo.NickName = userInfo.Nickname
			res.UserInfo.AvatarUrl = userInfo.Picture
			res.UserInfo.Gender = uint32(userInfo.Gender)
		}
	}
	return
}

// verifyToken 验证Token
func (ctl *PermissionService) verifyToken(ctx context.Context, token, oauth2ClientID string) (openID model.ID, err error) {
	if token == "" {
		err = status.Error(codes.InvalidArgument, "request token is empty")
		return
	}
	if oauth2ClientID == "" {
		err = status.Error(codes.InvalidArgument, "request oauth2_client_id is empty")
		return
	}
	var (
		exsit     bool
		claims    *oauth2.JwtClaims
		claimsErr error
	)
	rdsKey := fmt.Sprintf("oauth2_token_revocation:%s:access_token", oauth2ClientID)
	exsit, err = store.RedisClient.HExists(ctx, rdsKey, token).Result()
	if err != nil {
		err = status.Error(codes.Internal, fmt.Sprintf("check token revocation error: %s", err))
		return
	}
	if exsit {
		err = status.Error(codes.PermissionDenied, "token revocation")
		return
	}
	claims, claimsErr = oauth2.ParseJwtClaimsToken(token, []byte(viper.GetString("jwt.secret")))
	if claimsErr != nil {
		err = status.Error(codes.Internal, fmt.Sprintf("token is denied error: %s", claimsErr))
		return
	}
	if claims.VerifyAudience([]string{oauth2ClientID}, false) {
		err = status.Error(codes.PermissionDenied, fmt.Sprintf("token claims audience not equal to %s", oauth2ClientID))
		return
	}
	openID = model.ConvertStringToID(claims.Subject)
	return
}

// VerifyToken 验证Token
func (ctl *PermissionService) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (res *proto.VerifyTokenResponse, err error) {
	ctx = contexts.WithContext(ctx)
	err = checkReSource(ctx, req.Resource)
	if err != nil {
		return
	}
	var openID model.ID
	openID, err = ctl.verifyToken(ctx, req.Token, req.Oauth2ClientId)
	if err != nil {
		return
	}
	res = new(proto.VerifyTokenResponse)
	if req.ReturnUserInfo {
		user, userErr := service.User.GetOneByID(ctx, openID)
		if userErr != nil {
			err = status.Error(codes.Unavailable, userErr.Error())
			return
		}
		res.UserInfo = &proto.VerifyTokenResponse_UserInfo{
			OpenId:   convert.ToString(user.ID),
			Username: user.Username,
		}
		userInfo, userInfoErr := service.User.GetInfoOneByUserID(ctx, model.ConvertStringToID(res.UserInfo.OpenId))
		if userInfoErr == nil && userInfo != nil {
			res.UserInfo.NickName = userInfo.Nickname
			res.UserInfo.AvatarUrl = userInfo.Picture
			res.UserInfo.Gender = uint32(userInfo.Gender)
		}
	}
	return
}
