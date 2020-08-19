package dao

import (
	"context"
	"fmt"

	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/random"
	"github.com/nilorg/pkg/db"
	"github.com/nilorg/sdk/cache"
)

// Resourcer ...
type Resourcer interface {
	Insert(ctx context.Context, resource *model.Resource) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	Select(ctx context.Context, id model.ID) (resource *model.Resource, err error)
	Update(ctx context.Context, resource *model.Resource) (err error)
	LoadPolicy(ctx context.Context, resourceID model.ID) (results []*gormadapter.CasbinRule, err error)
}

type resource struct {
	cache cache.Cacher
}

func (*resource) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}

func (*resource) Insert(ctx context.Context, resource *model.Resource) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(resource).Error
	return
}
func (r *resource) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.Resource{}, id).Error
	if err != nil {
		return
	}
	err = r.cache.Remove(ctx, r.formatOneKey(id))
	return
}

func (*resource) selectOne(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	resource = new(model.Resource)
	err = gdb.Model(resource).Where("id = ?", id).Scan(resource).Error
	if err != nil {
		resource = nil
		return
	}
	return
}

func (r *resource) Select(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
	if store.FromSkipCacheContext(ctx) {
		return r.selectOne(ctx, id)
	}
	return r.selectFromCache(ctx, id)
}

func (r *resource) selectFromCache(ctx context.Context, id model.ID) (resource *model.Resource, err error) {
	resource = new(model.Resource)
	key := r.formatOneKey(id)
	err = r.cache.Get(ctx, key, resource)
	if err != nil {
		resource = nil
		if err == redis.Nil {
			resource, err = r.selectOne(ctx, id)
			if err != nil {
				return
			}
			err = r.cache.Set(ctx, key, resource, random.TimeDuration(300, 600))
		}
	}
	return
}

func (r *resource) Update(ctx context.Context, resource *model.Resource) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(resource).Update(resource).Error
	if err != nil {
		return
	}
	err = r.cache.Remove(ctx, r.formatOneKey(resource.ID))
	return
}

// LoadPolicy 加载规则
func (*resource) LoadPolicy(ctx context.Context, resourceID model.ID) (results []*gormadapter.CasbinRule, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("v1 like ?", fmt.Sprintf("resource:%d%%", resourceID)).Find(&results).Error
	return
}
