package store

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/go-redis/redis/v8"
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
var (
	// RedisClient redis 客户端
	RedisClient *redis.Client
	// DB ...
	DB *gorm.DB
	// Picture 头像
	Picture storage.Storager
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
func NewDBContext(dbs ...*gorm.DB) context.Context {
	if len(dbs) > 0 {
		return db.NewContext(context.Background(), dbs[0])
	}
	return db.NewContext(context.Background(), DB)
}

// NewCacheContext ...
func NewCacheContext(ctx context.Context, cache sdkCache.Cacher) context.Context {
	return sdkCache.NewCacheContext(ctx, cache)
}

// FromCacheContext ...
func FromCacheContext(ctx context.Context) (sdkCache.Cacher, error) {
	cache, ok := sdkCache.FromCacheContext(ctx)
	if !ok {
		return nil, ErrContextNotFoundCache
	}
	return cache, nil
}

func initPicture() {
	if viper.GetString("storage.type") == "default" {
		Picture = &storage.DefaultStorage{
			BasePath: filepath.Join(viper.GetString("storage.default.base_path"), "picture"),
		}
	}
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
