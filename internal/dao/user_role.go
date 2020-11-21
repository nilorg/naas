package dao

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/pkg/random"
	"github.com/nilorg/naas/internal/pkg/store"
	"github.com/nilorg/sdk/cache"
	"gorm.io/gorm"
)

// UserRoleer ...
type UserRoleer interface {
	Insert(ctx context.Context, m *model.UserRole) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByUserID(ctx context.Context, userID model.ID) (err error)
	DeleteByUserIDAndOrganizationID(ctx context.Context, userID, organizationID model.ID) (err error)
	DeleteByRoleCodeAndUserIDAndOrganizationID(ctx context.Context, roleCode model.Code, userID, organizationID model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.UserRole, err error)
	SelectAll(ctx context.Context) (roles []*model.UserRole, err error)
	Update(ctx context.Context, m *model.UserRole) (err error)
	SelectAllByUserID(ctx context.Context, userID model.ID) (m []*model.UserRole, err error)
	ListForRoleCodeByUserID(ctx context.Context, userID model.ID) (codes []model.Code, err error)
	ListForRoleCodeByUserIDAndOrganizationID(ctx context.Context, userID, organizationID model.ID) (codes []model.Code, err error)
	ListByUserIDAndOrganizationID(ctx context.Context, userID model.ID, organizationID model.ID) (list []*model.UserRole, err error)
	ExistByUserIDAndRoleCode(ctx context.Context, userID model.ID, roleCode model.Code) (exist bool, err error)
}

type userRole struct {
	cache cache.Cacher
}

func (*userRole) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (*userRole) formatListKey() string {
	return "list"
}

func (*userRole) formatUserListKey(userID model.ID) string {
	return fmt.Sprintf("list:user:%d", userID)
}

func (*userRole) formatUserAndOrgListKey(userID model.ID, orgID model.ID) string {
	return fmt.Sprintf("list:user:%d:org:%d", userID, orgID)
}

func (u *userRole) Insert(ctx context.Context, m *model.UserRole) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	if err != nil {
		return
	}
	err = u.cache.RemoveMatch(ctx, "list:*")
	return
}

func (u *userRole) delete(ctx context.Context, query interface{}, args ...interface{}) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where(query, args...).Delete(model.UserRole{}).Error
	return
}

func (u *userRole) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserRole{}, id).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(id))
	return
}

func (u *userRole) DeleteByUserID(ctx context.Context, userID model.ID) (err error) {
	err = u.delete(ctx, "user_id = ?", userID)
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatListKey(), u.formatUserListKey(userID))
	return
}

func (u *userRole) DeleteByUserIDAndOrganizationID(ctx context.Context, userID, organizationID model.ID) (err error) {
	err = u.delete(ctx, "user_id = ? and organization_id = ?", userID, organizationID)
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatListKey(), u.formatUserListKey(userID))
	return
}

func (u *userRole) DeleteByRoleCodeAndUserIDAndOrganizationID(ctx context.Context, roleCode model.Code, userID, organizationID model.ID) (err error) {
	err = u.delete(ctx, "role_code = ? and user_id = ? and organization_id = ?", roleCode, userID, organizationID)
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatListKey(), u.formatUserListKey(userID))
	return
}

func (u *userRole) selectOne(ctx context.Context, id model.ID) (m *model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	m = new(model.UserRole)
	err = gdb.Model(m).Where("id = ?", id).Take(m).Error
	if err != nil {
		m = nil
		return
	}
	return
}

func (u *userRole) Select(ctx context.Context, id model.ID) (m *model.UserRole, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectOne(ctx, id)
	}
	return u.selectFromCache(ctx, id)
}

func (u *userRole) selectFromCache(ctx context.Context, id model.ID) (m *model.UserRole, err error) {
	m = new(model.UserRole)
	key := u.formatOneKey(id)
	err = u.cache.Get(ctx, key, m)
	if err != nil {
		m = nil
		if err == redis.Nil {
			m, err = u.selectOne(ctx, id)
			if err != nil {
				return
			}
			err = u.cache.Set(ctx, key, m, random.TimeDuration(300, 600))
		}
	}
	return
}

func (u *userRole) Update(ctx context.Context, m *model.UserRole) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Save(m).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(m.ID))
	return
}

func (u *userRole) selectAllByUserID(ctx context.Context, userID model.ID) (roles []*model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.UserRole{}).Where("user_id = ?", userID).Find(&roles).Error
	return
}

func (u *userRole) SelectAllByUserID(ctx context.Context, userID model.ID) (roles []*model.UserRole, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectAllByUserID(ctx, userID)
	}
	return u.SelectAllByUserIDFromCache(ctx, userID)
}

func (u *userRole) listForRoleCodeByUserID(ctx context.Context, userID model.ID) (codes []model.Code, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	// resultResourceWebRouteID ...
	type resultResourceWebRouteID struct {
		RouteCode model.Code `gorm:"column:route_code"`
	}
	var items []*resultResourceWebRouteID
	err = gdb.Model(model.UserRole{}).Where("user_id = ?", userID).Find(&items).Error
	if err != nil {
		return
	}
	for _, item := range items {
		codes = append(codes, item.RouteCode)
	}
	return
}

func (u *userRole) ListForRoleCodeByUserID(ctx context.Context, userID model.ID) (codes []model.Code, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.listForRoleCodeByUserID(ctx, userID)
	}
	return u.listForRoleCodeByUserID(ctx, userID)
}

func (u *userRole) listForRoleCodeByUserIDAndOrganizationID(ctx context.Context, userID, organizationID model.ID) (codes []model.Code, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	// result ...
	type result struct {
		RouteCode model.Code `gorm:"column:role_code"`
	}
	var items []*result
	err = gdb.Model(model.UserRole{}).Where("user_id = ? and organization_id = ?", userID, organizationID).Find(&items).Error
	if err != nil {
		return
	}
	for _, item := range items {
		codes = append(codes, item.RouteCode)
	}
	return
}

func (u *userRole) ListForRoleCodeByUserIDAndOrganizationID(ctx context.Context, userID, organizationID model.ID) (codes []model.Code, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.listForRoleCodeByUserIDAndOrganizationID(ctx, userID, organizationID)
	}
	return u.listForRoleCodeByUserIDAndOrganizationID(ctx, userID, organizationID)
}

func (u *userRole) scanCacheID(ctx context.Context, items []*model.CacheIDPrimaryKey) (roles []*model.UserRole, err error) {
	for _, item := range items {
		i, ierr := u.selectFromCache(ctx, item.ID)
		if ierr != nil {
			err = ierr
			return
		}
		roles = append(roles, i)
	}
	return
}

func (u *userRole) SelectAllByUserIDFromCache(ctx context.Context, userID model.ID) (roles []*model.UserRole, err error) {
	key := u.formatUserListKey(userID)
	var items []*model.CacheIDPrimaryKey
	items, err = store.ScanByCacheID(store.NewCacheContext(ctx, u.cache), key, model.UserRole{}, "user_id = ?", userID)
	if err != nil {
		return
	}
	return u.scanCacheID(ctx, items)
}

func (u *userRole) ListByUserIDAndOrganizationID(ctx context.Context, userID model.ID, organizationID model.ID) (list []*model.UserRole, err error) {
	var (
		gdb *gorm.DB
	)
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	// expression := gdb.Table("user_role").Select("user_role.id").Joins("left join role on user_role.role_code = role.code").Where("user_role.user_id = ? and role.organization_id = ?", userID, organizationID)
	expression := gdb.Model(model.UserRole{}).Select("id").Where("user_id = ? and organization_id = ?", userID, organizationID)
	key := u.formatUserAndOrgListKey(userID, organizationID)
	var items []*model.CacheIDPrimaryKey
	items, err = store.ExpressionByCacheID(store.NewCacheContext(ctx, u.cache), key, expression)
	if err != nil {
		return
	}
	return u.scanCacheID(ctx, items)
}

func (u *userRole) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.UserRole{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func (u *userRole) ExistByUserIDAndRoleCode(ctx context.Context, userID model.ID, roleCode model.Code) (exist bool, err error) {
	return u.exist(ctx, "user_id = ? and role_code = ?", userID, roleCode)
}

func (u *userRole) selectAll(ctx context.Context) (roles []*model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.UserRole{}).Find(&roles).Error
	return
}

func (u *userRole) SelectAllFromCache(ctx context.Context) (roles []*model.UserRole, err error) {
	key := u.formatListKey()
	var items []*model.CacheIDPrimaryKey
	items, err = store.ScanByCacheID(store.NewCacheContext(ctx, u.cache), key, model.UserRole{}, "")
	if err != nil {
		return
	}
	return u.scanCacheID(ctx, items)
}

func (u *userRole) SelectAll(ctx context.Context) (roles []*model.UserRole, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectAll(ctx)
	}
	return u.SelectAllFromCache(ctx)
}
