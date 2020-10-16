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

// UserInfoer ...
type UserInfoer interface {
	SelectByUserID(ctx context.Context, userID model.ID) (mu *model.UserInfo, err error)
	Insert(ctx context.Context, mu *model.UserInfo) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByUserID(ctx context.Context, userID model.ID) (err error)
	DeleteInUserIDs(ctx context.Context, userIDs []model.ID) (err error)
	Select(ctx context.Context, id model.ID) (mu *model.UserInfo, err error)
	Update(ctx context.Context, mu *model.UserInfo) (err error)
}

type userInfo struct {
	cache cache.Cacher
}

func (*userInfo) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (*userInfo) formatOneUserIDKey(id model.ID) string {
	return fmt.Sprintf("user_id:%d", id)
}
func (u *userInfo) formatOneUserIDKeys(ids ...model.ID) (keys []string) {
	for _, id := range ids {
		keys = append(keys, u.formatOneUserIDKey(id))
	}
	return
}

func (*userInfo) selectByUserID(ctx context.Context, userID model.ID) (mu *model.UserInfo, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.UserInfo
	err = gdb.Where("user_id = ?", userID).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (u *userInfo) SelectByUserID(ctx context.Context, userID model.ID) (mu *model.UserInfo, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectByUserID(ctx, userID)
	}
	return u.selectByUserIDFromCache(ctx, userID)
}

func (u *userInfo) selectByUserIDFromCache(ctx context.Context, userID model.ID) (mu *model.UserInfo, err error) {
	mu = new(model.UserInfo)
	key := u.formatOneUserIDKey(userID)
	err = u.cache.Get(ctx, key, mu)
	if err != nil {
		mu = nil
		if err == redis.Nil {
			mu, err = u.selectByUserID(ctx, userID)
			if err != nil {
				return
			}
			err = u.cache.Set(ctx, key, mu, random.TimeDuration(300, 600))
		}
	}
	return
}

func (*userInfo) Insert(ctx context.Context, mu *model.UserInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mu).Create(mu).Error
	return
}

func (u *userInfo) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserInfo{}, id).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(id))
	return
}

func (u *userInfo) DeleteByUserID(ctx context.Context, userID model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("user_id = ?", userID).Delete(model.UserInfo{}).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneUserIDKey(userID))
	return
}

func (u *userInfo) DeleteInUserIDs(ctx context.Context, userIDs []model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("user_id in (?)", userIDs).Delete(model.UserInfo{}).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneUserIDKeys(userIDs...)...)
	return
}

func (*userInfo) selectOne(ctx context.Context, id model.ID) (mu *model.UserInfo, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	mu = new(model.UserInfo)
	err = gdb.Model(mu).Where("id = ?", id).Scan(mu).Error
	if err != nil {
		mu = nil
		return
	}
	return
}

func (u *userInfo) Select(ctx context.Context, id model.ID) (mu *model.UserInfo, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectOne(ctx, id)
	}
	return u.selectFromCache(ctx, id)
}

func (u *userInfo) selectFromCache(ctx context.Context, id model.ID) (m *model.UserInfo, err error) {
	m = new(model.UserInfo)
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

func (u *userInfo) Update(ctx context.Context, mu *model.UserInfo) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Save(mu).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(mu.ID))
	return
}
