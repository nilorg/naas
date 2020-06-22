package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nilorg/naas/internal/controller/api"
	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/naas/internal/controller/oidc"
	"github.com/nilorg/naas/internal/controller/service"
	"github.com/nilorg/naas/internal/controller/wellknown"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	store.Options(sessions.Options{
		Path:     viper.GetString("session.options.path"),
		Domain:   viper.GetString("session.options.domain"),
		MaxAge:   viper.GetInt("session.options.max_age"),
		Secure:   viper.GetBool("session.options.secure"),
		HttpOnly: viper.GetBool("session.options.http_only"),
	})
	r := gin.Default()
	r.Use(middleware.Header())
	r.Use(sessions.Sessions(viper.GetString("session.name"), store))
	// use ginSwagger middleware to
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	r.Static("/static", "./web/static")
	storageType := viper.GetString("storage.type")
	if storageType != "" {
		r.Static("/storage", viper.GetString(fmt.Sprintf("storage.%s.base_path", storageType)))
	}
	r.LoadHTMLGlob("./web/templates/oauth2/*")
	if viper.GetBool("server.admin.enabled") {
		if viper.GetBool("server.admin.external") {
			r.GET("/", func(ctx *gin.Context) {
				ctx.Redirect(302, viper.GetString("server.admin.external_url"))
			})
		} else {
			r.Static("/admin", "./web/templates/admin")
			r.GET("/", func(ctx *gin.Context) {
				ctx.Redirect(302, "/admin/index.html")
			})
		}
	}

	if viper.GetBool("server.oidc.enabled") {
		r.GET("/.well-known/jwks.json", wellknown.GetJwks)
		r.GET("/.well-known/openid-configuration", wellknown.GetOpenIDProviderMetadata)
	}

	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", oauth2.LoginPage)
		oauth2Group.POST("/login", oauth2.Login)
		oauth2Group.GET("/authorize", middleware.OAuth2AuthRequired, oauth2.AuthorizePage)
		oauth2Group.POST("/authorize", middleware.OAuth2AuthRequired, oauth2.Authorize)
		oauth2Group.POST("/token", oauth2.Token)
		if viper.GetBool("server.oauth2.device_authorization_endpoint_enabled") {
			oauth2Group.POST("/device/code", oauth2.DeviceCode)
		}
		if viper.GetBool("server.oauth2.introspection_endpoint_enabled") {
			oauth2Group.POST("/introspect", oauth2.TokenIntrospection)
		}
		if viper.GetBool("server.oauth2.revocation_endpoint_enabled") {
			oauth2Group.POST("/revoke", oauth2.TokenRevoke)
		}
	}
	if viper.GetBool("server.oidc.enabled") {
		oidcGroup := r.Group("/oidc")
		{
			if viper.GetBool("server.oidc.userinfo_endpoint_enabled") {
				oidcGroup.GET("/userinfo", middleware.OAuth2AuthUserinfoRequired(global.JwtPublicKey), oidc.GetUserinfo)
			}
		}
	}
	if viper.GetBool("server.admin.enabled") {
		r.GET("/roles", api.Role.Recursive)
		apiGroup := r.Group("api/v1", middleware.AdminAuthRequired(global.JwtPublicKey), middleware.AdminAuthSuperUserRequired())
		{
			apiGroup.GET("/users", api.User.ListByPaged)
			apiGroup.GET("/users/:user_id", api.User.GetOne)
			apiGroup.POST("/users", api.User.Create)
			apiGroup.PUT("/users/:user_id", api.User.Update)
			apiGroup.DELETE("/users/:user_id", api.User.Delete)

			apiGroup.GET("/oauth2/scopes", api.OAuth2.GetScopes)
			apiGroup.GET("/oauth2/clients", api.OAuth2.ClientListByPaged)
			apiGroup.POST("/oauth2/clients", api.OAuth2.CreateClient)
			apiGroup.PUT("/oauth2/clients/:client_id", api.OAuth2.UpdateClient)
			apiGroup.GET("/oauth2/clients/:client_id/scopes", api.OAuth2.GetClientScopes)
		}
	}
	r.Run(fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.oauth2.port"))) // listen and serve on 0.0.0.0:8080
}

// RunGRpc 运行Grpc
func RunGRpc() {
	gRPCServer := grpc.NewServer()
	service.RegisterGrpc(gRPCServer)
	// 在gRPC服务器上注册反射服务。
	reflection.Register(gRPCServer)
	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.grpc.port"))
	logger.Infof("%s grpc server listen: %s", "naas", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorf("net.Listen Error: %s", err)
		return
	}
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			logger.Infof("%s grpc server failed to serve: %v", "naas", err)
		}
	}()
	return
}

// RunGRpcGateway 运行Grpc网关
func RunGRpcGateway() {
	var (
		err error
	)
	gatewayMux := runtime.NewServeMux()
	err = service.RegisterGrpcGateway(gatewayMux)
	if err != nil {
		return
	}
	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.grpc.gateway.port"))
	srv := &http.Server{
		Addr:    addr,
		Handler: gatewayMux,
	}
	logger.Infof("启动GRpcGateway: %s", addr)
	go func() {
		if srvErr := srv.ListenAndServe(); srvErr != nil {
			log.Printf("%s gateway server listen: %v\n", viper.GetString("server.name"), srvErr)
		}
	}()
}
