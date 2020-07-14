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

// OAuth2Clienter oauth2 client 接口
type OAuth2Clienter interface {
	Insert(ctx context.Context, mc *model.OAuth2Client) (err error)
	Delete(ctx context.Context, id uint64) (err error)
	DeleteInIDs(ctx context.Context, ids []uint64) (err error)
	Update(ctx context.Context, mc *model.OAuth2Client) (err error)
	UpdateRedirectURI(ctx context.Context, id uint64, redirectURI string) (err error)
	SelectByID(ctx context.Context, clientID uint64) (mc *model.OAuth2Client, err error)
	ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total uint64, err error)
}

type oauth2Client struct {
	cache cache.Cacher
}

func (*oauth2Client) formatOneKey(id uint64) string {
	return fmt.Sprintf("id:%d", id)
}
func (o *oauth2Client) formatOneKeys(ids ...uint64) (keys []string) {
	for _, id := range ids {
		keys = append(keys, o.formatOneKey(id))
	}
	return
}

func (*oauth2Client) Insert(ctx context.Context, mc *model.OAuth2Client) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	return
}

func (o *oauth2Client) Delete(ctx context.Context, id uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(model.OAuth2Client{}, id).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(id))
	return
}

func (o *oauth2Client) DeleteInIDs(ctx context.Context, ids []uint64) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where("id in (?)", ids).Delete(model.OAuth2Client{}).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKeys(ids...)...)
	return
}

func (o *oauth2Client) Update(ctx context.Context, mc *model.OAuth2Client) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Update(mc).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(mc.ClientID))
	return
}

func (o *oauth2Client) UpdateRedirectURI(ctx context.Context, id uint64, redirectURI string) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Client{}).Where("client_id = ?", id).UpdateColumn("redirect_uri", redirectURI).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(id))
	return
}

func (*oauth2Client) selectByID(ctx context.Context, clientID uint64) (mc *model.OAuth2Client, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2Client
	err = gdb.Where("client_id = ?", clientID).First(&dbResult).Error
	if err != nil {
		return
	}
	mc = &dbResult
	return
}

func (o *oauth2Client) SelectByID(ctx context.Context, clientID uint64) (mc *model.OAuth2Client, err error) {
	if store.FromSkipCacheContext(ctx) {
		return o.selectByID(ctx, clientID)
	}
	return o.selectByIDFromCache(ctx, clientID)
}

func (o *oauth2Client) selectByIDFromCache(ctx context.Context, id uint64) (m *model.OAuth2Client, err error) {
	m = new(model.OAuth2Client)
	key := o.formatOneKey(id)
	err = o.cache.Get(ctx, key, m)
	if err != nil {
		m = nil
		if err == redis.Nil {
			m, err = o.selectByID(ctx, id)
			if err != nil {
				return
			}
			err = o.cache.Set(ctx, key, m, random.TimeDuration(300, 600))
		}
	}
	return
}

func (*oauth2Client) ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total uint64, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.OAuth2Client{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&clientList).Error
	return
}
