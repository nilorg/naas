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

// OAuth2Scoper ...
type OAuth2Scoper interface {
	Insert(ctx context.Context, m *model.OAuth2Scope) (err error)
	Delete(ctx context.Context, code model.Code) (err error)
	Select(ctx context.Context, code model.Code) (m *model.OAuth2Scope, err error)
	SelectAll(ctx context.Context) (m []*model.OAuth2Scope, err error)
	SelectByAllBasic(ctx context.Context) (m []*model.OAuth2Scope, err error)
	ListPaged(ctx context.Context, start, limit int) (scopes []*model.OAuth2Scope, total uint64, err error)
	Update(ctx context.Context, m *model.OAuth2Scope) (err error)
}

type oauth2Scope struct {
	cache cache.Cacher
}

func (*oauth2Scope) formatOneKey(code model.Code) string {
	return fmt.Sprintf("code:%s", code)
}

func (s *oauth2Scope) formatAllListKey() string {
	return "list:all"
}
func (s *oauth2Scope) formatAllListBasicKey() string {
	return "list:basic"
}

func (s *oauth2Scope) Insert(ctx context.Context, m *model.OAuth2Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Create(m).Error
	if err != nil {
		return
	}
	err = s.cache.RemoveMatch(ctx, "list:*")
	return
}

func (s *oauth2Scope) Delete(ctx context.Context, code model.Code) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.OAuth2Scope{}, "code = ?", code).Error
	if err != nil {
		return
	}
	err = s.cache.Remove(ctx, s.formatOneKey(code))
	return
}

func (s *oauth2Scope) selectAll(ctx context.Context) (m []*model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Scope{}).Find(&m).Error
	return
}

func (s *oauth2Scope) SelectAll(ctx context.Context) (m []*model.OAuth2Scope, err error) {
	if store.FromSkipCacheContext(ctx) {
		return s.SelectAll(ctx)
	}
	return s.selectAllFromCache(ctx)
}

func (s *oauth2Scope) selectAllFromCache(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	key := s.formatAllListKey()
	var items []*model.CacheCodePrimaryKey
	items, err = store.ScanByCacheCode(store.NewCacheContext(ctx, s.cache), key, model.OAuth2Scope{}, "1 = 1")
	if err != nil {
		return
	}
	return s.scanCacheCode(ctx, items)
}

func (s *oauth2Scope) SelectByAllBasic(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	if store.FromSkipCacheContext(ctx) {
		return s.selectByAllBasic(ctx)
	}
	return s.selectByAllBasicFromCache(ctx)
}

func (s *oauth2Scope) selectByAllBasic(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.OAuth2Scope{}).Where("type = ?", model.OAuth2ScopeTypeBasic).Find(&scopes).Error
	return
}

func (s *oauth2Scope) selectByAllBasicFromCache(ctx context.Context) (scopes []*model.OAuth2Scope, err error) {
	key := s.formatAllListKey()
	var items []*model.CacheCodePrimaryKey
	items, err = store.ScanByCacheCode(store.NewCacheContext(ctx, s.cache), key, model.OAuth2Scope{}, "type = ?", model.OAuth2ScopeTypeBasic)
	if err != nil {
		return
	}
	return s.scanCacheCode(ctx, items)
}

func (s *oauth2Scope) selectOne(ctx context.Context, code model.Code) (m *model.OAuth2Scope, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.OAuth2Scope
	err = gdb.First(&dbResult, "code = ?", code).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (s *oauth2Scope) Select(ctx context.Context, code model.Code) (m *model.OAuth2Scope, err error) {
	if store.FromSkipCacheContext(ctx) {
		return s.selectOne(ctx, code)
	}
	return s.selectFromCache(ctx, code)
}

func (s *oauth2Scope) selectFromCache(ctx context.Context, code model.Code) (m *model.OAuth2Scope, err error) {
	m = new(model.OAuth2Scope)
	key := s.formatOneKey(code)
	err = s.cache.Get(ctx, key, m)
	if err != nil {
		m = nil
		if err == redis.Nil {
			m, err = s.selectOne(ctx, code)
			if err != nil {
				return
			}
			err = s.cache.Set(ctx, key, m, random.TimeDuration(300, 600))
		}
	}
	return
}

func (s *oauth2Scope) Update(ctx context.Context, m *model.OAuth2Scope) (err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(m).Update(m).Error
	if err != nil {
		return
	}
	err = s.cache.Remove(ctx, s.formatOneKey(m.Code))
	return
}

func (s *oauth2Scope) scanCacheCode(ctx context.Context, items []*model.CacheCodePrimaryKey) (scopes []*model.OAuth2Scope, err error) {
	for _, item := range items {
		i, ierr := s.selectFromCache(ctx, item.Code)
		if ierr != nil {
			err = ierr
			return
		}
		scopes = append(scopes, i)
	}
	return
}

func (s *oauth2Scope) ListPaged(ctx context.Context, start, limit int) (scopes []*model.OAuth2Scope, total uint64, err error) {
	var gdb *gorm.DB
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	expression := gdb.Model(&model.OAuth2Scope{})
	expression.Count(&total)
	err = expression.Offset(start).Limit(limit).Find(&scopes).Error
	return
}
