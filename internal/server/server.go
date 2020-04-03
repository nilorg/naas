package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nilorg/naas/internal/controller/gateway"
	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/pkg/logger"

	// swagger doc ...
	_ "github.com/nilorg/naas/docs"
	"github.com/nilorg/naas/internal/server/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RunHTTP ...
func RunHTTP() {
	store, err := redis.NewStore(10, "tcp", viper.GetString("session.redis.address"), viper.GetString("session.redis.password"), []byte(viper.GetString("session.secret")))
	if err != nil {
		logger.Errorf("redis.NewStore Error:", err)
		return
	}
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
	apiGroup := r.Group("api/v1")
	{
		apiGroup.POST("/auth/login", jwtMiddleware.LoginHandler)
		authorized := apiGroup.Group("/")
		authorized.Use(jwtMiddleware.MiddlewareFunc())
		{
			// Refresh time can be longer than token timeout
			authorized.GET("/auth/refresh_token", jwtMiddleware.RefreshHandler)
		}
	}
	r.Run(fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.oauth2.port"))) // listen and serve on 0.0.0.0:8080
}

func RunGRpcGateway() {
	var (
		err error
	)
	gatewayMux := runtime.NewServeMux()
	err = gateway.Service(gatewayMux)
	if err != nil {
		return
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.grpc.gateway.port")),
		Handler: gatewayMux,
	}
	go func() {
		if srvErr := srv.ListenAndServe(); srvErr != nil {
			log.Printf("%s gateway server listen: %v\n", viper.GetString("server.name"), srvErr)
		}
	}()
}
