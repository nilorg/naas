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

// UserOrganizationer ...
type UserOrganizationer interface {
	Insert(ctx context.Context, m *model.UserOrganization) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.UserOrganization, err error)
	Update(ctx context.Context, m *model.UserOrganization) (err error)
	SelectAllByUserID(ctx context.Context, userID model.ID) (m []*model.UserOrganization, err error)
}

type userOrganization struct {
	cache cache.Cacher
}

func (*userOrganization) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (*userOrganization) formatUserListKey(userID model.ID) string {
	return fmt.Sprintf("list:user:%d", userID)
}

func (u *userOrganization) Insert(ctx context.Context, m *model.UserOrganization) (err error) {
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

func (u *userOrganization) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserOrganization{}, id).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(id))
	return
}

func (u *userOrganization) selectOne(ctx context.Context, id model.ID) (m *model.UserOrganization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	m = new(model.UserOrganization)
	err = gdb.Model(m).Where("id = ?", id).Take(m).Error
	if err != nil {
		m = nil
		return
	}
	return
}

func (u *userOrganization) Select(ctx context.Context, id model.ID) (m *model.UserOrganization, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectOne(ctx, id)
	}
	return u.selectFromCache(ctx, id)
}

func (u *userOrganization) selectFromCache(ctx context.Context, id model.ID) (m *model.UserOrganization, err error) {
	m = new(model.UserOrganization)
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

func (u *userOrganization) Update(ctx context.Context, m *model.UserOrganization) (err error) {
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

func (u *userOrganization) selectAllByUserID(ctx context.Context, userID model.ID) (roles []*model.UserOrganization, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.UserOrganization{}).Where("user_id = ?", userID).Find(&roles).Error
	return
}

func (u *userOrganization) SelectAllByUserID(ctx context.Context, userID model.ID) (roles []*model.UserOrganization, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectAllByUserID(ctx, userID)
	}
	return u.SelectAllByUserIDFromCache(ctx, userID)
}

func (u *userOrganization) scanCacheID(ctx context.Context, items []*model.CacheIDPrimaryKey) (roles []*model.UserOrganization, err error) {
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

func (u *userOrganization) SelectAllByUserIDFromCache(ctx context.Context, userID model.ID) (roles []*model.UserOrganization, err error) {
	key := u.formatUserListKey(userID)
	var items []*model.CacheIDPrimaryKey
	items, err = store.ScanByCacheID(store.NewCacheContext(ctx, u.cache), key, model.UserOrganization{}, "user_id = ?", userID)
	if err != nil {
		return
	}
	return u.scanCacheID(ctx, items)
}
