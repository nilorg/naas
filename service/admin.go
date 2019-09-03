package service

import (
	"github.com/nilorg/naas/dao"
	"github.com/nilorg/naas/errors"
	"github.com/nilorg/naas/model"
	"github.com/nilorg/naas/module/store"
	"github.com/nilorg/pkg/logger"
)

type admin struct {
}

// GetUserByUsername 根据用户名获取管理员
func (a *admin) GetUserByUsername(username string) (ma *model.Admin, err error) {
	return dao.Admin.SelectByUsername(store.NewDBContext(), username)
}

// Login 登录 ...
func (a *admin) Login(username, password string) (su *model.SessionAccount, err error) {
	var ma *model.Admin
	ma, err = a.GetUserByUsername(username)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if ma.Username == username && ma.Password == password {
		su = &model.SessionAccount{
			UserID:   ma.ID,
			UserName: ma.Username,
		}
	} else {
		err = errors.ErrUsernameOrPassword
	}
	return
}
