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

// OAuth2ClientScoper oauth2 client 范围 接口
type OAuth2ClientScoper interface {
	Select(ctx context.Context, id uint64) (mc *model.OAuth2ClientScope, err error)
	SelectFromCache(ctx context.Context, id uint64) (mc *model.OAuth2ClientScope, err error)
	SelectByOAuth2ClientID(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error)
	SelectByOAuth2ClientIDFromCache(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error)
	Insert(ctx context.Context, mc *model.OAuth2ClientScope) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	Update(ctx context.Context, mc *model.OAuth2ClientScope) (err error)
}

type oauth2ClientScope struct {
	cache cache.Cacher
}

func (*oauth2ClientScope) formatOneKey(id uint64) string {
	return fmt.Sprintf("id:%d", id)
}
func (*oauth2ClientScope) formatClientListKey(clientID uint64) string {
	return fmt.Sprintf("list:clientid:%d", clientID)
}

func (*oauth2ClientScope) Select(ctx context.Context, id uint64) (mc *model.OAuth2ClientScope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2ClientScope
	err = gdb.First(&dbResult, id).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}

func (o *oauth2ClientScope) SelectFromCache(ctx context.Context, id uint64) (mc *model.OAuth2ClientScope, err error) {
	mc = new(model.OAuth2ClientScope)
	key := o.formatOneKey(id)
	err = o.cache.Get(ctx, key, mc)
	if err != nil {
		mc = nil
		if err == redis.Nil {
			mc, err = o.Select(ctx, id)
			if err != nil {
				return
			}
			err = o.cache.Set(ctx, key, mc, random.TimeDuration(300, 600))
		}
	}
	return
}

func (*oauth2ClientScope) SelectByOAuth2ClientID(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("oauth2_client_id = ?", clientID).Find(&scopes).Error
	return
}

func (o *oauth2ClientScope) SelectByOAuth2ClientIDFromCache(ctx context.Context, clientID uint64) (scopes []*model.OAuth2ClientScope, err error) {
	key := o.formatClientListKey(clientID)
	var items []*model.CacheIDPrimaryKey
	items, err = store.ScanByCacheID(store.NewCacheContext(ctx, o.cache), key, model.OAuth2ClientScope{}, "oauth2_client_id = ?", clientID)
	if err != nil {
		return
	}
	return o.scanCacheID(ctx, items)
}

func (o *oauth2ClientScope) Insert(ctx context.Context, mc *model.OAuth2ClientScope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	if err != nil {
		return
	}
	err = o.cache.RemoveMatch(ctx, "list:*")
	return
}

func (o *oauth2ClientScope) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(model.OAuth2ClientScope{}, id).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(id))
	return
}

func (o *oauth2ClientScope) Update(ctx context.Context, mc *model.OAuth2ClientScope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Update(mc).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(mc.ID))
	return
}

func (o *oauth2ClientScope) scanCacheID(ctx context.Context, items []*model.CacheIDPrimaryKey) (scopes []*model.OAuth2ClientScope, err error) {
	for _, item := range items {
		i, ierr := o.SelectFromCache(ctx, item.ID)
		if ierr != nil {
			err = ierr
			return
		}
		scopes = append(scopes, i)
	}
	return
}
