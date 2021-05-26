package store

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	pkgStorage "github.com/nilorg/pkg/storage"
	"github.com/nilorg/sdk/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/nilorg/naas/internal/model"
	"github.com/spf13/viper"
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
	logrus.Println(pong)
}

func initMySQL() {
	var err error
	gormLogger := logger.Discard
	if viper.GetBool("mysql.log") {
		std := *logrus.StandardLogger()
		std.SetReportCaller(false)
		gormLogger = logger.New(&std, logger.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	}
	DB, err = gorm.Open(
		mysql.Open(viper.GetString("mysql.address")),
		&gorm.Config{
			Logger:                 gormLogger,
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	)
	if err != nil {
		panic("failed to connect database.")
	}
	// Migrate the schema
	DB.AutoMigrate(
		&model.OAuth2Scope{},
		&model.OAuth2Client{},
		&model.OAuth2ClientScope{},
		&model.OAuth2ClientInfo{},
		&model.Organization{},
		&model.Role{},
		&model.RoleResourceRelation{},
		&model.User{},
		&model.UserInfo{},
		&model.UserRole{},
		&model.UserOrganization{},
		&model.Resource{},
		&model.ResourceRoute{},
		&model.ResourceMenu{},
		&model.ResourceAction{},
	)
}

func initPicture() {
	typ := viper.GetString("storage.type")
	if typ == "default" {
		basePath := viper.GetString("storage.default.base_path")
		if strings.TrimSpace(basePath) == "" {
			basePath = "./web/storage"
		}
		basePath = filepath.Join(basePath, "picture")
		Picture = NewDefaultStorage(
			&storage.DefaultStorage{
				BasePath: basePath,
			},
		)
	} else if typ == "oss" {
		client, err := oss.New(viper.GetString("storage.oss.endpoint"), viper.GetString("storage.oss.access.key_id"), viper.GetString("storage.oss.access.key_secret"))
		if err != nil {
			logrus.Fatalln(err)
			return
		}
		bucket := viper.GetString("storage.oss.bucket")
		var ossStorage *pkgStorage.AliyunOssStorage
		ossStorage, err = pkgStorage.NewAliyunOssStorage(client, false, []string{bucket})
		if err != nil {
			logrus.Fatalln(err)
			return
		}
		Picture = NewDefaultStorage(ossStorage)
	} else {
		panic("type error")
	}
}

// FileRename 重命名文件名
func FileRename(filename string) string {
	suffix := filepath.Ext(filename)
	return uuid.New().String() + suffix
}
