package service

import (
	"bytes"
	"fmt"
	"image"
	"net/url"
	"path"
	"time"

	"image/png"

	"context"

	"github.com/nilorg/go-wechat/oauth"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/module/weixin"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/sdk/convert"
	"github.com/o1egl/govatar"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type user struct {
}

// createPicture 创建头像
func createPicture(typ, id string) (bs string, err error) {
	var img image.Image
	img, err = govatar.GenerateForUsername(govatar.MALE, id)
	if err != nil {
		return
	}
	buff := bytes.NewBuffer(nil)
	err = png.Encode(buff, img)
	if err != nil {
		return
	}
	ctx := context.Background()
	filename := fmt.Sprintf("%s-%s.png", typ, id)
	_, err = store.Picture.Upload(ctx, buff, filename)
	if err != nil {
		return
	}
	var u *url.URL
	u, err = url.Parse(viper.GetString("storage.public_path"))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, "/picture/", filename)
	bs = u.String()
	return
}

const (
	createUserTypePassword  = "password"
	createUserTypeWxOpenID  = "wxopenid"
	createUserTypeWxUnionID = "wxunionid"
)

func (u *user) create(ctx context.Context, username, password, openID, typ string) (err error) {
	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	var exist bool
	exist, err = dao.User.ExistByUsername(ctx, username)
	if err != nil {
		tran.Rollback()
		return
	}
	if exist {
		tran.Rollback()
		err = errors.ErrUsernameExist
		return
	}
	user := &model.User{
		Username: username,
		Password: password, // TODO: 后期需要使用加密，或者前台加密
	}
	err = dao.User.Insert(ctx, user)
	if err != nil {
		tran.Rollback()
		return
	}
	if typ != createUserTypePassword {
		var userThirdType model.UserThirdType
		if typ == createUserTypeWxUnionID {
			userThirdType = model.UserThirdTypeWxUnionID
		} else if typ == createUserTypeWxOpenID {
			userThirdType = model.UserThirdTypeWxOpenIDForKfpt
		} else {
			err = fmt.Errorf("创建类型错误")
			tran.Rollback()
			return
		}
		exist, err = dao.UserThird.ExistByThirdIDAndThirdType(ctx, openID, userThirdType)
		if err != nil {
			tran.Rollback()
			return
		}
		if exist {
			tran.Rollback()
			err = errors.ErrThirdExistUser
			return
		}
		userThird := &model.UserThird{
			UserID:    user.ID,
			ThirdType: userThirdType,
			ThirdID:   openID,
		}
		err = dao.UserThird.Insert(ctx, userThird)
		if err != nil {
			tran.Rollback()
			return
		}
	}
	userInfo := &model.UserInfo{
		UserID:   user.ID,
		Nickname: fmt.Sprintf("用户%d", user.ID),
	}
	userInfo.Picture, err = createPicture("user", convert.ToString(user.ID))
	if err != nil {
		tran.Rollback()
		return
	}
	err = dao.UserInfo.Insert(ctx, userInfo)
	if err != nil {
		tran.Rollback()
		return
	}
	err = tran.Commit().Error
	return
}

// Create 创建用户
func (u *user) Create(ctx context.Context, username, password string) (err error) {
	return u.create(ctx, username, password, "", createUserTypePassword)
}

// CreateFromWeixin 从微信角度创建用户
func (u *user) CreateFromWeixin(ctx context.Context, wxUnionID string) (err error) {
	return u.create(ctx, wxUnionID, wxUnionID, wxUnionID, createUserTypeWxUnionID)
}

// InitFromWeixinKfptOpenID 使用微信OpenID初始化账户
func (u *user) InitFromWeixinKfptOpenID(ctx context.Context, wxOpenID string) (su *model.SessionAccount, err error) {
	err = u.create(ctx, wxOpenID, wxOpenID, wxOpenID, createUserTypeWxOpenID)
	if err != nil {
		return
	}
	su, err = u.loginForWxKfptOpenID(ctx, wxOpenID)
	return
}

// UserUpdateModel ...
type UserUpdateModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Update 修改用户
func (u *user) Update(ctx context.Context, id model.ID, update *UserUpdateModel) (err error) {
	var (
		user          *model.User
		usernameExist bool
	)
	user, err = dao.User.Select(ctx, id)
	if err != nil {
		return
	}
	if user.Username != update.Username {
		usernameExist, err = dao.User.ExistByUsername(ctx, update.Username)
		if err != nil {
			return
		}
		if usernameExist {
			err = errors.ErrUsernameExist
			return
		}
	}
	user.Username = update.Username
	// TODO: 后期需要使用加密，或者前台加密
	user.Password = update.Password
	err = dao.User.Update(ctx, user)
	if err != nil {
		return
	}
	return
}

// GetUserByUsername 根据用户名获取用户
func (u *user) GetUserByUsername(ctx context.Context, username string) (usr *model.User, err error) {
	return dao.User.SelectByUsername(ctx, username)
}

// GetOneByID 根据ID获取用户
func (u *user) GetOneByID(ctx context.Context, id model.ID) (usr *model.User, err error) {
	return dao.User.Select(ctx, id)
}

// GetInfoOneByUserID 根据用户ID获取信息
func (u *user) GetInfoOneByUserID(ctx context.Context, userID model.ID) (usr *model.UserInfo, err error) {
	return dao.UserInfo.SelectByUserID(ctx, userID)
}

// GetInfoOneByCache 根据用户ID获取信息
func (u *user) GetInfoOneByCache(ctx context.Context, userID model.ID) (usr *model.User, usrInfo *model.UserInfo, err error) {
	usr, err = dao.User.Select(ctx, userID)
	if err != nil {
		return
	}
	usrInfo, err = dao.UserInfo.SelectByUserID(ctx, userID)
	return
}

// Login 登录 ...
func (u *user) Login(ctx context.Context, username, password string) (su *model.SessionAccount, err error) {
	var usr *model.User
	usr, err = u.GetUserByUsername(ctx, username)
	if err != nil {
		logrus.Errorln(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.ErrUsernameOrPassword
		}
		return
	}
	redisCountKey := fmt.Sprintf("user:login:%d:errcount", usr.ID)
	redisLockKey := fmt.Sprintf("user:login:%d:lock", usr.ID)
	const lock = "lock"
	lockValue := store.RedisClient.Get(ctx, redisLockKey).String()
	if lockValue == lock {
		err = fmt.Errorf("账户已被锁定，请24小时候再试")
		return
	}
	if usr.Username == username && usr.Password == password {
		su = &model.SessionAccount{
			UserID:   usr.ID,
			UserName: usr.Username,
		}
		var userInfo *model.UserInfo
		userInfo, err = u.GetInfoOneByUserID(ctx, usr.ID)
		if err == nil {
			su.Nickname = userInfo.Nickname
			su.Picture = userInfo.Picture
		}
		err = store.RedisClient.Del(ctx, redisCountKey, redisLockKey).Err()
		if err != nil {
			logrus.WithContext(ctx).Errorf("用户(%d)登录成功后删除相关的锁：%v", usr.ID, err)
			err = errors.ErrBadRequest
			return
		}
	} else {
		var count int
		count, err = store.RedisClient.Get(ctx, redisCountKey).Int()
		if err != nil {
			logrus.WithContext(ctx).Errorf("判断登录次数错误：", err)
			err = errors.ErrBadRequest
			return
		}
		timeout := 24 * time.Hour
		if count > 4 {
			err = store.RedisClient.Set(ctx, redisLockKey, lock, timeout).Err()
			if err != nil {
				logrus.WithContext(ctx).Errorf("用户(%d)密码错误登录次数过多加锁：%v", usr.ID, err)
				err = errors.ErrBadRequest
				return
			}
		}
		err = store.RedisClient.Set(ctx, redisCountKey, count+1, timeout).Err()
		if err != nil {
			logrus.WithContext(ctx).Errorf("添加用户(%d)判断登录次数错误：%v", usr.ID, err)
			err = errors.ErrBadRequest
			return
		}
		err = errors.ErrUsernameOrPassword
	}
	return
}

// loginForUserID 登录 ...
func (u *user) loginForUserID(ctx context.Context, userID model.ID) (su *model.SessionAccount, err error) {
	var usr *model.User
	usr, err = u.GetOneByID(ctx, userID)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	su = &model.SessionAccount{
		UserID:   usr.ID,
		UserName: usr.Username,
	}
	var userInfo *model.UserInfo
	userInfo, err = u.GetInfoOneByUserID(ctx, usr.ID)
	if err == nil {
		su.Nickname = userInfo.Nickname
		su.Picture = userInfo.Picture
	}
	return
}

// loginForWxKfptOpenID 微信OpenID登录
func (u *user) loginForWxKfptOpenID(ctx context.Context, wxOpenID string) (su *model.SessionAccount, err error) {
	var userThird *model.UserThird
	userThird, err = dao.UserThird.SelectByThirdIDAndThirdType(ctx, wxOpenID, model.UserThirdTypeWxOpenIDForKfpt)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	su, err = u.loginForUserID(ctx, userThird.UserID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

// LoginForWeixinKfptCode 根据微信Code进行登录
func (u *user) LoginForWeixinKfptCode(ctx context.Context, code string) (su *model.SessionAccount, err error) {
	xoauth := oauth.NewOAuth(weixin.KfptWechatClientConfig)
	var reply *oauth.AccessTokenReply
	reply, err = xoauth.GetAccessToken(code)
	if err != nil {
		logrus.WithContext(ctx).Errorf("xoauth.GetAccessToken:%s", err)
		return
	}
	var exist bool
	exist, err = dao.UserThird.ExistByThirdIDAndThirdType(ctx, reply.OpenID, model.UserThirdTypeWxOpenIDForKfpt)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if exist {
		var userThird *model.UserThird
		userThird, err = dao.UserThird.SelectByThirdIDAndThirdType(ctx, reply.OpenID, model.UserThirdTypeWxOpenIDForKfpt)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		su, err = u.loginForUserID(ctx, userThird.UserID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
		}
	} else {
		su = &model.SessionAccount{
			WxOpenID: reply.OpenID,
			Action:   model.SessionAccountActionBindWx,
		}
		err = errors.ErrThirdUserNotFound
	}
	return
}

// LoginForWeixinFwhCode 根据微信服务号Code进行登录
func (u *user) LoginForWeixinFwhCode(ctx context.Context, code string) (su *model.SessionAccount, st *model.SessionThirdBind, err error) {
	xoauth := oauth.NewOAuth(weixin.FwhWechatClientConfig)
	var reply *oauth.AccessTokenReply
	reply, err = xoauth.GetAccessToken(code)
	if err != nil {
		logrus.WithContext(ctx).Errorf("xoauth.GetAccessToken:%s", err)
		return
	}
	var exist bool
	exist, err = dao.UserThird.ExistByThirdIDAndThirdType(ctx, reply.OpenID, model.UserThirdTypeWxOpenIDForFwh)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
		return
	}
	if exist {
		var userThird *model.UserThird
		userThird, err = dao.UserThird.SelectByThirdIDAndThirdType(ctx, reply.OpenID, model.UserThirdTypeWxOpenIDForFwh)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
			return
		}
		su, err = u.loginForUserID(ctx, userThird.UserID)
		if err != nil {
			logrus.WithContext(ctx).Errorln(err)
		}
	} else {
		st = &model.SessionThirdBind{
			ThirdID: reply.OpenID,
			Type:    model.UserThirdTypeWxOpenIDForFwh,
		}
		err = errors.ErrThirdUserNotFound
	}
	return
}

// LoginForUserID 根据用户ID登录
func (u *user) LoginForUserID(ctx context.Context, userID model.ID) (su *model.SessionAccount, err error) {
	su, err = u.loginForUserID(ctx, userID)
	if err != nil {
		logrus.WithContext(ctx).Errorln(err)
	}
	return
}

func (u *user) ListPaged(ctx context.Context, start, limit int) (result []*model.ResultUserInfo, total int64, err error) {
	var (
		userList []*model.User
	)
	userList, total, err = dao.User.ListPaged(ctx, start, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	for _, user := range userList {
		userInfo, userInfoErr := u.GetInfoOneByUserID(ctx, user.ID)
		resultInfo := &model.ResultUserInfo{
			User: user,
		}
		if userInfoErr == nil {
			resultInfo.UserInfo = userInfo
		}
		result = append(result, resultInfo)
	}
	return
}

// DeleteByID 根据ID删除用户
func (u *user) DeleteByIDs(ctx context.Context, ids ...model.ID) (err error) {
	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	err = dao.User.DeleteInIDs(ctx, ids)
	if err != nil {
		tran.Rollback()
		return
	}
	err = dao.UserInfo.DeleteInUserIDs(ctx, ids)
	if err != nil {
		tran.Rollback()
		return
	}
	err = tran.Commit().Error
	return
}

// UserUpdateRoleModel ...
type UserUpdateRoleModel struct {
	Roles          []model.Code `json:"roles"`
	OrganizationID model.ID     `json:"organization_id"`
}

// UpdateRole 修改用户角色
func (u *user) UpdateRole(ctx context.Context, userID model.ID, update *UserUpdateRoleModel) (err error) {
	var (
		exist bool
	)
	exist, err = dao.User.ExistByID(ctx, userID)
	if err != nil {
		return
	}
	if !exist {
		err = errors.ErrUserNotFound
		return
	}
	var routeCodes []model.Code
	routeCodes, err = dao.UserRole.ListForRoleCodeByUserIDAndOrganizationID(ctx, userID, update.OrganizationID)
	if err != nil {
		return
	}
	added, deleted := model.DiffCodeSlice(routeCodes, update.Roles)

	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	defer func() {
		if err != nil {
			tran.Rollback()
		}
	}()

	for _, r := range deleted {
		err = dao.UserRole.DeleteByRoleCodeAndUserIDAndOrganizationID(ctx, r, userID, update.OrganizationID)
		if err != nil {
			return
		}
		// user, role, domain := formatRoleForUserInDomain(userID, update.OrganizationID, r)
		// _, err = casbin.Enforcer.DeleteRoleForUserInDomain(user, role, domain)
		// if err != nil {
		// 	return
		// }
	}

	for _, r := range added {
		err = dao.UserRole.Insert(ctx, &model.UserRole{
			UserID:         userID,
			RoleCode:       model.Code(r),
			OrganizationID: update.OrganizationID,
		})
		if err != nil {
			return
		}
		// user, role, domain := formatRoleForUserInDomain(userID, update.OrganizationID, r)
		// _, err = casbin.Enforcer.AddRoleForUserInDomain(user, role, domain)
		// if err != nil {
		// 	return
		// }
	}
	err = tran.Commit().Error
	return
}

// UserUpdateOrganizationModel ...
type UserUpdateOrganizationModel struct {
	Organizations []model.ID `json:"organizations"`
}

// UpdateOrganization 修改用户组织
func (u *user) UpdateOrganization(ctx context.Context, userID model.ID, update *UserUpdateOrganizationModel) (err error) {
	var (
		exist bool
	)
	exist, err = dao.User.ExistByID(ctx, userID)
	if err != nil {
		return
	}
	if !exist {
		err = errors.ErrUserNotFound
		return
	}
	// TODO: 这地方有待优化
	tran := store.DB.Begin()
	ctx = contexts.NewGormTranContext(ctx, tran)
	defer func() {
		if err != nil {
			tran.Rollback()
		}
	}()
	err = dao.UserOrganization.DeleteByUserID(ctx, userID)
	if err != nil {
		return
	}
	for _, orgID := range update.Organizations {
		err = dao.UserOrganization.Insert(ctx, &model.UserOrganization{
			UserID:         userID,
			OrganizationID: orgID,
		})
		if err != nil {
			return
		}
	}
	err = tran.Commit().Error
	return
}

// BindThird 绑定第三方
func (u *user) BindThird(ctx context.Context, userID model.ID, thirdID string, thirdType model.UserThirdType) (err error) {
	var exist bool
	exist, err = dao.User.ExistByID(ctx, userID)
	if err != nil {
		return
	}
	if !exist {
		err = errors.ErrUserNotFound
		return
	}
	if !model.UserThirdTypeVerify(thirdType) {
		err = fmt.Errorf("绑定第三方类型错误：%v", thirdType)
		return
	}
	// 判断用户有没有绑定第三方
	exist, err = dao.UserThird.ExistByUserIDAndThirdType(ctx, userID, thirdType)
	if err != nil {
		return
	}
	if exist {
		err = errors.ErrUserExistThird
		return
	}
	// 判断第三方有没有绑定用户
	exist, err = dao.UserThird.ExistByThirdIDAndThirdType(ctx, thirdID, thirdType)
	if err != nil {
		return
	}
	if exist {
		err = errors.ErrThirdExistUser
		return
	}
	userThird := &model.UserThird{
		UserID:    userID,
		ThirdID:   thirdID,
		ThirdType: thirdType,
	}
	err = dao.UserThird.Insert(ctx, userThird)
	return
}
