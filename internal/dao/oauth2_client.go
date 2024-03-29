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

// OAuth2Clienter oauth2 client 接口
type OAuth2Clienter interface {
	Insert(ctx context.Context, mc *model.OAuth2Client) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteInIDs(ctx context.Context, ids []model.ID) (err error)
	Update(ctx context.Context, mc *model.OAuth2Client) (err error)
	UpdateRedirectURI(ctx context.Context, id model.ID, redirectURI string) (err error)
	SelectByID(ctx context.Context, clientID model.ID) (mc *model.OAuth2Client, err error)
	ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total int64, err error)
	ExistByID(ctx context.Context, clientID model.ID) (exist bool, err error)
}

type oauth2Client struct {
	cache cache.Cacher
}

func (*oauth2Client) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}
func (o *oauth2Client) formatOneKeys(ids ...model.ID) (keys []string) {
	for _, id := range ids {
		keys = append(keys, o.formatOneKey(id))
	}
	return
}

func (*oauth2Client) Insert(ctx context.Context, mc *model.OAuth2Client) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(mc).Error
	return
}

func (o *oauth2Client) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
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

func (o *oauth2Client) DeleteInIDs(ctx context.Context, ids []model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
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
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(mc).Save(mc).Error
	if err != nil {
		return
	}
	err = o.cache.Remove(ctx, o.formatOneKey(mc.ClientID))
	return
}

func (o *oauth2Client) UpdateRedirectURI(ctx context.Context, id model.ID, redirectURI string) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
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

func (*oauth2Client) selectByID(ctx context.Context, clientID model.ID) (mc *model.OAuth2Client, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
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

func (o *oauth2Client) SelectByID(ctx context.Context, clientID model.ID) (mc *model.OAuth2Client, err error) {
	if store.FromSkipCacheContext(ctx) {
		return o.selectByID(ctx, clientID)
	}
	return o.selectByIDFromCache(ctx, clientID)
}

func (o *oauth2Client) selectByIDFromCache(ctx context.Context, id model.ID) (m *model.OAuth2Client, err error) {
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

func (*oauth2Client) ListPaged(ctx context.Context, start, limit int) (clientList []*model.OAuth2Client, total int64, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.OAuth2Client{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&clientList).Error
	return
}

func (*oauth2Client) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.OAuth2Client{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// ExistByID 判断ID是否存在
func (o *oauth2Client) ExistByID(ctx context.Context, id model.ID) (exist bool, err error) {
	return o.exist(ctx, "client_id = ?", id)
}
