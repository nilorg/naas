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

// UserThirder ...
type UserThirder interface {
	Insert(ctx context.Context, m *model.UserThird) (err error)
	Delete(ctx context.Context, id model.ID) (err error)
	DeleteByUserIDAndThirdID(ctx context.Context, userID model.ID, thirdID string) (err error)
	Select(ctx context.Context, id model.ID) (m *model.UserThird, err error)
	SelectByThirdIDAndThirdType(ctx context.Context, thirdID string, thirdType model.UserThirdType) (m *model.UserThird, err error)
	Update(ctx context.Context, m *model.UserThird) (err error)
	SelectAllByUserID(ctx context.Context, userID model.ID) (m []*model.UserThird, err error)
	ExistByThirdIDAndThirdType(ctx context.Context, thirdID string, thirdType model.UserThirdType) (exist bool, err error)
	ExistByUserIDAndThirdType(ctx context.Context, userID model.ID, thirdType model.UserThirdType) (exist bool, err error)
}

type userThird struct {
	cache cache.Cacher
}

func (*userThird) formatOneKey(id model.ID) string {
	return fmt.Sprintf("id:%d", id)
}
func (*userThird) formatListKey() string {
	return "list"
}
func (*userThird) formatUserListKey(userID model.ID) string {
	return fmt.Sprintf("list:user:%d", userID)
}
func (*userThird) formatOneByThirdIDAndThirdTypeKey(thirdID string, thirdType model.UserThirdType) string {
	return fmt.Sprintf("third:%s:type:%s", thirdID, thirdType)
}
func (u *userThird) Insert(ctx context.Context, m *model.UserThird) (err error) {
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

func (u *userThird) delete(ctx context.Context, query interface{}, args ...interface{}) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Where(query, args...).Delete(model.UserThird{}).Error
	return
}

func (u *userThird) Delete(ctx context.Context, id model.ID) (err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Delete(&model.UserThird{}, id).Error
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatOneKey(id))
	return
}

func (u *userThird) DeleteByUserIDAndThirdID(ctx context.Context, userID model.ID, thirdID string) (err error) {
	err = u.delete(ctx, "user_id = ? and third_id = ?", userID, thirdID)
	if err != nil {
		return
	}
	err = u.cache.Remove(ctx, u.formatListKey(), u.formatUserListKey(userID))
	return
}

func (u *userThird) selectOne(ctx context.Context, id model.ID) (m *model.UserThird, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	m = new(model.UserThird)
	err = gdb.Model(m).Where("id = ?", id).Take(m).Error
	if err != nil {
		m = nil
		return
	}
	return
}

func (u *userThird) Select(ctx context.Context, id model.ID) (m *model.UserThird, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectOne(ctx, id)
	}
	return u.selectFromCache(ctx, id)
}

func (*userThird) selectByThirdIDAndThirdType(ctx context.Context, thirdID string, thirdType model.UserThirdType) (m *model.UserThird, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var dbResult model.UserThird
	err = gdb.Where("third_id = ? and third_type = ?", thirdID, thirdType).First(&dbResult).Error
	if err != nil {
		return
	}
	m = &dbResult
	return
}

func (u *userThird) SelectByThirdIDAndThirdType(ctx context.Context, thirdID string, thirdType model.UserThirdType) (m *model.UserThird, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectByThirdIDAndThirdType(ctx, thirdID, thirdType)
	}
	return u.SelectByThirdIDAndThirdTypeFromCache(ctx, thirdID, thirdType)
}

func (u *userThird) SelectByThirdIDAndThirdTypeFromCache(ctx context.Context, thirdID string, thirdType model.UserThirdType) (m *model.UserThird, err error) {
	m = new(model.UserThird)
	key := u.formatOneByThirdIDAndThirdTypeKey(thirdID, thirdType)
	err = u.cache.Get(ctx, key, m)
	if err != nil {
		m = nil
		if err == redis.Nil {
			m, err = u.selectByThirdIDAndThirdType(ctx, thirdID, thirdType)
			if err != nil {
				return
			}
			err = u.cache.Set(ctx, key, m, random.TimeDuration(300, 600))
		}
	}
	return
}

func (u *userThird) selectFromCache(ctx context.Context, id model.ID) (m *model.UserThird, err error) {
	m = new(model.UserThird)
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

func (u *userThird) Update(ctx context.Context, m *model.UserThird) (err error) {
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

func (u *userThird) selectAllByUserID(ctx context.Context, userID model.ID) (thirds []*model.UserThird, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	err = gdb.Model(model.UserThird{}).Where("user_id = ?", userID).Find(&thirds).Error
	return
}

func (u *userThird) SelectAllByUserID(ctx context.Context, userID model.ID) (thirds []*model.UserThird, err error) {
	if store.FromSkipCacheContext(ctx) {
		return u.selectAllByUserID(ctx, userID)
	}
	return u.SelectAllByUserIDFromCache(ctx, userID)
}

func (u *userThird) scanCacheID(ctx context.Context, items []*model.CacheIDPrimaryKey) (thirds []*model.UserThird, err error) {
	for _, item := range items {
		i, ierr := u.selectFromCache(ctx, item.ID)
		if ierr != nil {
			err = ierr
			return
		}
		thirds = append(thirds, i)
	}
	return
}

func (u *userThird) SelectAllByUserIDFromCache(ctx context.Context, userID model.ID) (thirds []*model.UserThird, err error) {
	key := u.formatUserListKey(userID)
	var items []*model.CacheIDPrimaryKey
	items, err = store.ScanByCacheID(store.NewCacheContext(ctx, u.cache), key, model.UserThird{}, "user_id = ?", userID)
	if err != nil {
		return
	}
	return u.scanCacheID(ctx, items)
}

func (u *userThird) exist(ctx context.Context, query interface{}, args ...interface{}) (exist bool, err error) {
	var gdb *gorm.DB
	gdb, err = contexts.FromGormContext(ctx)
	if err != nil {
		return
	}
	var count int64
	err = gdb.Model(&model.UserThird{}).Where(query, args...).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func (u *userThird) ExistByThirdIDAndThirdType(ctx context.Context, thirdID string, thirdType model.UserThirdType) (exist bool, err error) {
	return u.exist(ctx, "third_id = ? and third_type = ?", thirdID, thirdType)
}

func (u *userThird) ExistByUserIDAndThirdType(ctx context.Context, userID model.ID, thirdType model.UserThirdType) (exist bool, err error) {
	return u.exist(ctx, "user_id = ? and third_type = ?", userID, thirdType)
}
