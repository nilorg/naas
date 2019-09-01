package store

import (
	"github.com/go-redis/redis"
	"github.com/nilorg/pkg/logger"

	// use db mysql
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nilorg/naas/model"
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
	// Migrate the schema
	DB.AutoMigrate(&model.OAuth2Client{}, &model.OAuth2Scope{}, &model.User{})
}
