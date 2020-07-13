package service

import (
	"bytes"
	"fmt"
	"image"
	"net/url"
	"path"

	"image/png"

	"context"

	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/convert"
	"github.com/o1egl/govatar"
	"github.com/spf13/viper"
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
	u, err = url.Parse(viper.GetString("storage.picture.public_path"))
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, filename)
	bs = u.String()
	return
}

// Create 创建用户
func (u *user) Create(username, password string) (err error) {
	tran := store.DB.Begin()
	ctx := store.NewDBContext(tran)
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

// UserUpdateModel ...
type UserUpdateModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Update 修改用户
func (u *user) Update(id uint64, update *UserUpdateModel) (err error) {
	ctx := store.NewDBContext()
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
func (u *user) GetUserByUsername(username string) (usr *model.User, err error) {
	return dao.User.SelectByUsername(store.NewDBContext(), username)
}

// GetOneByID 根据ID获取用户
func (u *user) GetOneByID(id uint64) (usr *model.User, err error) {
	return dao.User.Select(store.NewDBContext(), id)
}

// GetInfoOneByUserID 根据用户ID获取信息
func (u *user) GetInfoOneByUserID(userID uint64) (usr *model.UserInfo, err error) {
	return dao.UserInfo.SelectByUserID(store.NewDBContext(), userID)
}

// GetInfoOneByCache 根据用户ID获取信息
func (u *user) GetInfoOneByCache(userID uint64) (usr *model.User, usrInfo *model.UserInfo, err error) {
	usr, err = dao.User.Select(store.NewDBContext(), userID)
	if err != nil {
		return
	}
	usrInfo, err = dao.UserInfo.SelectByUserID(store.NewDBContext(), userID)
	return
}

// Login 登录 ...
func (u *user) Login(username, password string) (su *model.SessionAccount, err error) {
	var usr *model.User
	usr, err = u.GetUserByUsername(username)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if usr.Username == username && usr.Password == password {
		su = &model.SessionAccount{
			UserID:   usr.ID,
			UserName: usr.Username,
		}
		var userInfo *model.UserInfo
		userInfo, err = u.GetInfoOneByUserID(usr.ID)
		if err == nil {
			su.Nickname = userInfo.Nickname
			su.Picture = userInfo.Picture
		}
	} else {
		err = errors.ErrUsernameOrPassword
	}
	return
}

func (u *user) ListPaged(start, limit int) (result []*model.ResultUserInfo, total uint64, err error) {
	var (
		userList []*model.User
	)
	userList, total, err = dao.User.ListPaged(store.NewDBContext(), start, limit)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = nil
		}
		return
	}
	for _, user := range userList {
		userInfo, userInfoErr := u.GetInfoOneByUserID(user.ID)
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
func (u *user) DeleteByIDs(ids ...uint64) (err error) {
	tran := store.DB.Begin()
	ctx := store.NewDBContext(tran)
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
