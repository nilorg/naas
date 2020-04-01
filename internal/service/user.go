package service

import (
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/pkg/logger"
)

type user struct {
}

// GetUserByUsername 根据用户名获取用户
func (u *user) GetUserByUsername(username string) (usr *model.User, err error) {
	return dao.User.SelectByUsername(store.NewDBContext(), username)
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
