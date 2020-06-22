package service

import (
	"context"

	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/pkg/logger"
)

type role struct {
}

func (r *role) Recursive() []*model.Role {
	var (
		rootRoles []*model.Role
		err       error
	)
	ctx := store.NewDBContext()
	rootRoles, err = dao.Role.SelectByRoot(ctx)
	if err != nil {
		logger.Errorf("dao.Role.SelectByRoot: %s", err)
	}
	r.recursive(ctx, rootRoles)
	return rootRoles
}

func (r *role) recursive(ctx context.Context, roles []*model.Role) {
	if len(roles) == 0 {
		return
	}
	var (
		childRoles []*model.Role
		err        error
	)
	for _, role := range roles {
		childRoles, err = dao.Role.SelectByParentCode(ctx, role.Code)
		if err != nil {
			logger.Errorf("dao.Role.SelectByRoot: %s", err)
		}
		r.recursive(ctx, childRoles)
		role.ChildRoles = childRoles
	}
}
