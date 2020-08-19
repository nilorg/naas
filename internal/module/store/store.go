package store

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/nilorg/pkg/db"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/cache"
	sdkCache "github.com/nilorg/sdk/cache"
	"github.com/nilorg/sdk/storage"

	// use db mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/pkg/random"
	"github.com/spf13/viper"
)

var (
	// ErrContextNotFoundCache 上下文不存在Cache错误
	ErrContextNotFoundCache = errors.New("上下文中没有获取到Cache")
)

// DefaultStorage 默认存储
type DefaultStorage struct {
	storage.Storager
}

// NewDefaultStorage 创建默认存储
func NewDefaultStorage(storage storage.Storager) *DefaultStorage {
	return &DefaultStorage{
		Storager: storage,
	}
}

// MaxMemory 最大上传大小
func (*DefaultStorage) MaxMemory() int64 {
	maxMemory := viper.GetInt64("storage.max_memory")
	if maxMemory <= 0 {
		maxMemory = 20 // 20 MB
	}
	return maxMemory << 20
}

var (
	// RedisClient redis 客户端
	RedisClient *redis.Client
	// DB ...
	DB *gorm.DB
	// Picture 头像
	Picture *DefaultStorage
)

// Init 初始化
func Init() {
	initRedis()
	initMySQL()
	initPicture()
}

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	logger.Println(pong)
}

func initMySQL() {
	var err error
	DB, err = gorm.Open("mysql", viper.GetString("mysql.address"))
	if err != nil {
		panic("failed to connect database")
	}
	DB.LogMode(viper.GetBool("mysql.log"))
	DB.DB().SetMaxOpenConns(viper.GetInt("mysql.max_open"))
	DB.DB().SetMaxIdleConns(viper.GetInt("mysql.max_idle"))
	// 关闭复数表名，如果设置为true，`User`表的表名就会是`user`，而不是`users`
	DB.SingularTable(true)
	// Migrate the schema
	DB.AutoMigrate(
		&model.Admin{},
		&model.OAuth2Scope{},
		&model.OAuth2Client{},
		&model.OAuth2ClientScope{},
		&model.OAuth2ClientInfo{},
		&model.Organization{},
		&model.OrganizationRole{},
		&model.Role{},
		&model.RoleResourceWebRoute{},
		&model.RoleResourceWebFunction{},
		&model.User{},
		&model.UserInfo{},
		&model.UserRole{},
		&model.Resource{},
		&model.ResourceWebRoute{},
		&model.ResourceWebFunction{},
		&model.ResourceWebComponent{},
		&model.ResourceWebFunctionComponent{},
	)
}

// NewDBContext ...
func NewDBContext(ctx context.Context, dbs ...*gorm.DB) context.Context {
	if len(dbs) > 0 {
		return db.NewContext(ctx, dbs[0])
	}
	return db.NewContext(ctx, DB)
}

// NewCacheContext 创建缓存对象到上下文中
func NewCacheContext(ctx context.Context, cache sdkCache.Cacher) context.Context {
	return sdkCache.NewCacheContext(ctx, cache)
}

// FromCacheContext 从缓存上下文中获取缓存对象
func FromCacheContext(ctx context.Context) (sdkCache.Cacher, error) {
	cache, ok := sdkCache.FromCacheContext(ctx)
	if !ok {
		return nil, ErrContextNotFoundCache
	}
	return cache, nil
}

type skipCache struct{}

// NewSkipCacheContext 创建跳过缓存到上下文
func NewSkipCacheContext(ctx context.Context, skip ...bool) context.Context {
	s := true
	if len(skip) > 0 {
		s = skip[0]
	}
	return context.WithValue(ctx, skipCache{}, s)
}

// FromSkipCacheContext 从上下文中获取跳过缓存变量
func FromSkipCacheContext(ctx context.Context) (skip bool) {
	var ok bool
	skip, ok = ctx.Value(skipCache{}).(bool)
	if !ok {
		skip = false
	}
	return
}

func initPicture() {
	if viper.GetString("storage.type") == "default" {
		Picture = NewDefaultStorage(
			&storage.DefaultStorage{
				BasePath: filepath.Join(viper.GetString("storage.default.base_path"), "picture"),
			},
		)
	}
}

// FileRename 重命名文件名
func FileRename(filename string) string {
	suffix := filepath.Ext(filename)
	return uuid.New().String() + suffix
}

// ScanByCacheID ...
func ScanByCacheID(ctx context.Context, cacheKey string, table interface{}, query interface{}, args ...interface{}) (items []*model.CacheIDPrimaryKey, err error) {
	err = scanByCache(ctx, cacheKey, table, &items, query, args...)
	return
}

// ScanByCacheCode ...
func ScanByCacheCode(ctx context.Context, cacheKey string, table interface{}, query interface{}, args ...interface{}) (items []*model.CacheCodePrimaryKey, err error) {
	err = scanByCache(ctx, cacheKey, table, &items, query, args...)
	return
}

// scanByCache ...
func scanByCache(ctx context.Context, cacheKey string, table interface{}, values interface{}, query interface{}, args ...interface{}) (err error) {
	var (
		gdb   *gorm.DB
		cache cache.Cacher
	)
	gdb, err = db.FromContext(ctx)
	if err != nil {
		return
	}
	cache, err = FromCacheContext(ctx)
	if err != nil {
		return
	}
	err = cache.Get(ctx, cacheKey, values)
	if err != nil {
		if err == redis.Nil {
			if err = gdb.Model(table).Where(query, args...).Scan(values).Error; err != nil {
				return
			}
			if err = cache.Set(ctx, cacheKey, values, random.TimeDuration(300, 600)); err != nil {
				return
			}
		} else {
			return
		}
	}
	return
}
