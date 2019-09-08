package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/controller/oauth2"
	"github.com/nilorg/pkg/logger"

	// swagger doc ...
	_ "github.com/nilorg/naas/docs"
	"github.com/nilorg/naas/server/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RunHTTP ...
func RunHTTP() {
	store, _ := redis.NewStore(viper.GetInt("session.redis.db"), "tcp", viper.GetString("session.redis.address"), "", []byte(viper.GetString("session.secret")))
	r := gin.Default()
	r.Use(middleware.Header())
	r.Use(sessions.Sessions(viper.GetString("session.name"), store))
	// use ginSwagger middleware to
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	r.Static("/static", "./web/static")
	r.Static("/www", "./web/templates/www")
	r.LoadHTMLGlob("./web/templates/oauth2/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/www/index.html")
	})

	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", oauth2.LoginPage)
		oauth2Group.POST("/login", oauth2.Login)
		oauth2Group.GET("/authorize", middleware.AuthRequired, oauth2.AuthorizePage)
		oauth2Group.POST("/authorize", middleware.AuthRequired, oauth2.Authorize)
		oauth2Group.POST("/token", oauth2.Token)
	}
	// the jwt middleware
	jwtMiddleware, err := middleware.NewJwtMiddleware()
	if err != nil {
		logger.Fatalf("JWT Error:%s", err)
	}
	apiGroup := r.Group("/v1")
	{
		apiGroup.POST("/auth/login", jwtMiddleware.LoginHandler)
		authorized := apiGroup.Group("/")
		authorized.Use(jwtMiddleware.MiddlewareFunc())
		{
			// Refresh time can be longer than token timeout
			authorized.GET("/auth/refresh_token", jwtMiddleware.RefreshHandler)
		}
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
