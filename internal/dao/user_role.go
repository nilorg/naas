package dao

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/random"
	"github.com/nilorg/pkg/db"
	"github.com/nilorg/sdk/cache"
)

// UserRoleer ...
type UserRoleer interface {
	Insert(ctx context.Context, m *model.UserRole) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (m *model.UserRole, err error)
	Update(ctx context.Context, m *model.UserRole) (err error)
	SelectAllByUserID(ctx context.Context, userID model.ID) (m []*model.UserRole, err error)
}

type userRole struct {
	cache cache.Cacher
}

func (*userRole) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (*userRole) formatUserListKey(userID model.ID) string {
	return fmt.Sprintf("list:user:%d", userID)
}

func (u *userRole) Insert(ctx context.Context, m *model.UserRole) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
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

func (u *userRole) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
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

func (u *userRole) selectOne(ctx context.Context, id model.ID) (m *model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	m = new(model.UserRole)
	err = gdb.Model(m).Where("id = ?", id).Scan(m).Error
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
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Update(m).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(m.ID))
	return
}

func (u *userRole) selectAllByUserID(ctx context.Context, userID model.ID) (roles []*model.UserRole, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
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
