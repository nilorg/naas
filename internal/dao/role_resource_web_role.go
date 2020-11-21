package dao

import (
	"context"
	"fmt"

	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"gorm.io/gorm"
)

// RoleResourceWebRouter ...
type RoleResourceWebRouter interface {
	Insert(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByRoleCode(ctx context.Context, roleCode model.Code) (err error)
	DeleteByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode model.Code, resourceWebRouteID model.ID) (err error)
	Select(ctx context.Context, id model.ID) (roleResourceWebRoute *model.RoleResourceWebRoute, err error)
	SelectAll(ctx context.Context) (roleResourceWebRoutes []*model.RoleResourceWebRoute, err error)
	Update(ctx context.Context, roleResourceWebRoute *model.RoleResourceWebRoute) (err error)
	ExistByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode model.Code, resourceWebRouteID model.ID) (exist bool, err error)
	ExistByResourceWebRouteID(ctx context.Context, resourceWebRouteID model.ID) (exist bool, err error)
	ListByRoleCode(ctx context.Context, roleCode model.Code, limit int) (list []*model.RoleResourceWebRoute, err error)
	ListPagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (list []*model.RoleResourceWebRoute, total int64, err error)
	ListPagedByRoleCodeAndResourceServerID(ctx context.Context, start, limit int, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceWebRoute, total int64, err error)
	ListByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceWebRoute, err error)
	ListForResourceWebRouteIDByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (ids []model.ID, err error)
}

type roleResourceWebRoute struct {
}

func (*roleResourceWebRoute) formatRoleListKey(roleCode model.Code) string {
	return fmt.Sprintf("list:role:%s", roleCode)
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

// ExistByResourceWebRouteID 判断根据资源web路由ID
func (r *roleResourceWebRoute) ExistByResourceWebRouteID(ctx context.Context, resourceWebRouteID model.ID) (exist bool, err error) {
	return r.exist(ctx, "resource_web_route_id = ?", resourceWebRouteID)
}

func (r *roleResourceWebRoute) ListByRoleCode(ctx context.Context, roleCode model.Code, limit int) (list []*model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	exp := gdb.Model(&model.RoleResourceWebRoute{}).Where("role_code = ?", roleCode)
	if limit > 0 {
		exp = exp.Offset(0).Limit(limit)
	}
	err = exp.Find(&list).Error
	return
}

func (r *roleResourceWebRoute) ListPagedByRoleCode(ctx context.Context, start, limit int, roleCode model.Code) (list []*model.RoleResourceWebRoute, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.RoleResourceWebRoute{}).Where("role_code = ?", roleCode)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *roleResourceWebRoute) ListPagedByRoleCodeAndResourceServerID(ctx context.Context, start, limit int, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceWebRoute, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.RoleResourceWebRoute{}).Where("role_code = ? and resource_server_id = ?", roleCode, resourceServerID)
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&list).Error
	return
}

func (r *roleResourceWebRoute) ListByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (list []*model.RoleResourceWebRoute, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(&model.RoleResourceWebRoute{}).Where("role_code = ? and resource_server_id = ?", roleCode, resourceServerID).Find(&list).Error
	return
}

func (r *roleResourceWebRoute) ListForResourceWebRouteIDByRoleCodeAndResourceServerID(ctx context.Context, roleCode model.Code, resourceServerID model.ID) (ids []model.ID, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	// resultResourceWebRouteID ...
	type resultResourceWebRouteID struct {
		ResourceWebRouteID model.ID `gorm:"column:resource_web_route_id"`
	}

	var items []*resultResourceWebRouteID
	err = gdb.Model(&model.RoleResourceWebRoute{}).Where("role_code = ? and resource_server_id = ?", roleCode, resourceServerID).Find(&items).Error
	if err != nil {
		return
	}
	for _, i := range items {
		ids = append(ids, i.ResourceWebRouteID)
	}
	return
}

func (r *roleResourceWebRoute) delete(ctx context.Context, query interface{}, args ...interface{}) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where(query, args...).Delete(model.RoleResourceWebRoute{}).Error
	return
}

func (r *roleResourceWebRoute) DeleteByRoleCode(ctx context.Context, roleCode model.Code) (err error) {
	err = r.delete(ctx, "role_code = ?", roleCode)
	if err != nil {
		return
	}
	return
}

func (r *roleResourceWebRoute) DeleteByRoleCodeAndResourceWebRouteID(ctx context.Context, roleCode model.Code, resourceWebRouteID model.ID) (err error) {
	err = r.delete(ctx, "role_code = ? and resource_web_route_id = ?", roleCode, resourceWebRouteID)
	if err != nil {
		return
	}
	return
}
