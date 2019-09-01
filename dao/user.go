package dao

import (
	"github.com/nilorg/naas/module/store"
	"github.com/nilorg/naas/model"
)

type user struct {
}

func (*user) SelectByUsername(username string) (mu *model.User, err error) {
	var dbResult model.User
	err = store.DB.Where("username = ?", username).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}
