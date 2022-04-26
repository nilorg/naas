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

// Userer ...
type Userer interface {
	SelectByUsername(ctx context.Context, username string) (mu *model.User, err error)
	Insert(ctx context.Context, mu *model.User) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids []model.ID) (err error)
	Select(ctx context.Context, id model.ID) (mu *model.User, err error)
	Update(ctx context.Context, mu *model.User) (err error)
	ListPaged(ctx context.Context, start, limit int) (user []*model.User, total int64, err error)
	ExistByUsername(ctx context.Context, username string) (exist bool, err error)
	ExistByWxUnionID(ctx context.Context, wxUnionID string) (exist bool, err error)
	ExistByID(ctx context.Context, id model.ID) (exist bool, err error)
}

type user struct {
	cache cache.Cacher
}

func (*user) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (u *user) formatOneKeys(ids ...model.ID) (keys []string) {
	for _, id := range ids {
		keys = append(keys, u.formatOneKey(id))
	}
	return
}

func (*user) SelectByUsername(ctx context.Context, username string) (mu *model.User, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.User
	err = gdb.Where("username = ?", username).First(&dbResult).Error
	if err != nil {
		return
	}
	mu = &dbResult
	return
}

func (*user) Insert(ctx context.Context, mu *model.User) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mu).Error
	return
}

func (u *user) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.User{}, id).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(id))
	return
}

func (u *user) DeleteInIDs(ctx context.Context, ids []model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in (?)", ids).Delete(model.User{}).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKeys(ids...)...)
	return
}

func (*user) selectOne(ctx context.Context, id model.ID) (mu *model.User, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	mu = new(model.User)
	err = gdb.Model(mu).Where("id = ?", id).Take(mu).Error
	if err != nil {
		mu = nil
		return
	}
	return
}

func (u *user) Select(ctx context.Context, id model.ID) (mu *model.User, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectOne(ctx, id)
	}
	return u.selectFromCache(ctx, id)
}

func (u *user) selectFromCache(ctx context.Context, id model.ID) (m *model.User, err error) {
	m = new(model.User)
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

func (*user) Update(ctx context.Context, mu *model.User) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mu).Save(mu).Error
	if err != nil {
		return
	}
	return
}

func (*user) ListPaged(ctx context.Context, start, limit int) (user []*model.User, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.User{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&user).Error
	return
}

func (*user) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.User{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByPhone 判断用户名是否存在
func (u *user) ExistByUsername(ctx context.Context, username string) (exist bool, err error) {
	return u.exist(ctx, "username = ?", username)
}

// ExistByWxUnionID 判断微信UnionID是否存在
func (u *user) ExistByWxUnionID(ctx context.Context, wxUnionID string) (exist bool, err error) {
	return u.exist(ctx, "wx_union_id = ?", wxUnionID)
}

// ExistByID 判断ID是否存在
func (u *user) ExistByID(ctx context.Context, id model.ID) (exist bool, err error) {
	return u.exist(ctx, "id = ?", id)
}
