package store

import (
	"context"

	"github.com/go-redis/redis/v7"
	"github.com/nilorg/pkg/db"
	"github.com/nilorg/pkg/logger"

	"github.com/jinzhu/gorm"
	// use db mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nilorg/naas/internal/model"
	"github.com/spf13/viper"
)

var (
	// RedisClient redis 客户端
	RedisClient *redis.Client
	// DB ...
	DB *gorm.DB
)

// Init 初始化
func Init() {
	initRedis()
	initMySQL()
}

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	pong, err := RedisClient.Ping().Result()
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
		&model.OAuth2Client{},
		&model.OAuth2ClientInfo{},
		&model.OAuth2Scope{},
		&model.Organization{},
		&model.OrganizationRole{},
		&model.Role{},
		&model.RoleWebFunction{},
		&model.User{},
		&model.UserRole{},
		&model.WebComponent{},
		&model.WebFunction{},
		&model.WebFunctionComponent{},
	)
}

// NewDBContext ...
func NewDBContext() context.Context {
	return db.NewContext(context.Background(), DB)
}
