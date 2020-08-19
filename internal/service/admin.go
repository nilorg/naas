package service

import (
	"context"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/pkg/logger"
)

type admin struct {
}

// GetUserByUsername 根据用户名获取管理员
func (a *admin) GetUserByUsername(ctx context.Context, username string) (ma *model.Admin, err error) {
	return dao.Admin.SelectByUsername(ctx, username)
}

// Login 登录 ...
func (a *admin) Login(ctx context.Context, username, password string) (su *model.SessionAccount, err error) {
	var ma *model.Admin
	ma, err = a.GetUserByUsername(ctx, username)
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
