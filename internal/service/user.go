package service

import (
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/convert"
)

type user struct {
}

// GetUserByUsername 根据用户名获取用户
func (u *user) GetUserByUsername(username string) (usr *model.User, err error) {
	return dao.User.SelectByUsername(store.NewDBContext(), username)
}

// GetOneByID 根据ID获取用户
func (u *user) GetOneByID(id string) (usr *model.User, err error) {
	return dao.User.Select(store.NewDBContext(), convert.ToUint64(id))
}

// GetInfoOneByUserID 根据用户ID获取信息
func (u *user) GetInfoOneByUserID(userID string) (usr *model.UserInfo, err error) {
	return dao.UserInfo.Select(store.NewDBContext(), convert.ToUint64(userID))
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
	} else {
		err = errors.ErrUsernameOrPassword
	}
	return
}
