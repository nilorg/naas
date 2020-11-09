package dao

import (
	"context"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// RoleResourceWebRouter ...
type RoleResourceWebRouter interface {
	Insert(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (roleResourceWebRoute *model.RoleResourceWebRoute, err error)
	SelectAll(ctx context.Context) (roleResourceWebRoutes []*model.RoleResourceWebRoute, err error)
	Update(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	ExistByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode model.Code, resourceWebRouteID model.ID) (exist bool, err error)
}

type roleResourceWebRoute struct {
}

func (*roleResourceWebRoute) Insert(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(roleResourceWebRoute).Error
	return
}
func (*roleResourceWebRoute) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.RoleResourceWebRoute{}, id).Error
	return
}
func (*roleResourceWebRoute) Select(ctx context.Context, id model.ID) (roleResourceWebRoute *model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	roleResourceWebRoute = new(model.RoleResourceWebRoute)
	err = gdb.Model(roleResourceWebRoute).Where("id = ?", id).Take(roleResourceWebRoute).Error
	if err != nil {
		roleResourceWebRoute = nil
		return
	}
	return
}
func (*roleResourceWebRoute) Update(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(roleResourceWebRoute).Save(roleResourceWebRoute).Error
	return
}

func (*roleResourceWebRoute) SelectAll(ctx context.Context) (roleResourceWebRoutes []*model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.RoleResourceWebRoute{}).Find(&roleResourceWebRoutes).Error
	return
}

func (*roleResourceWebRoute) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.RoleResourceWebRoute{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByRoleCodeAndResourceWebRouteID 判断根据RoleCode和资源web路由ID
func (r *roleResourceWebRoute) ExistByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode model.Code, resourceWebRouteID model.ID) (exist bool, err error) {
	return r.exist(ctx, "role_code = ? and resource_web_route_id = ?", roleCode, resourceWebRouteID)
}
