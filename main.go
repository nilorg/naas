package main

import (
	"fmt"
	"github.com/nilorg/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/controller/oauth2"
	"github.com/nilorg/naas/middleware"
	"github.com/nilorg/naas/module"
	"github.com/spf13/viper"
)

func init() {
	logger.Init()
	viper.SetConfigType("toml") // or viper.SetConfigType("YAML")
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	module.Init()
}

func main() {
	store, _ := redis.NewStore(viper.GetInt("session.redis.db"), "tcp", viper.GetString("session.redis.address"), "", []byte(viper.GetString("session.secret")))
	r := gin.Default()
	r.Use(sessions.Sessions(viper.GetString("session.name"), store))
	r.Static("/assets", "./assets")
	r.Static("/www", "./templates/www")
	r.LoadHTMLGlob("./templates/oauth2/*")
	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", oauth2.LoginPage)
		oauth2Group.POST("/login", oauth2.Login)
		oauth2Group.GET("/authorize", middleware.AuthRequired, oauth2.AuthorizePage)
		oauth2Group.POST("/authorize", middleware.AuthRequired, oauth2.Authorize)
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
