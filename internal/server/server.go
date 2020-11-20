package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nilorg/naas/internal/controller/api"
	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/naas/internal/controller/oidc"
	"github.com/nilorg/naas/internal/controller/open"
	"github.com/nilorg/naas/internal/controller/service"
	"github.com/nilorg/naas/internal/controller/wellknown"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// swagger doc ...
	_ "github.com/nilorg/naas/docs"
	"github.com/nilorg/naas/internal/server/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"

	//ginSwagger "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/nilorg/naas/internal/pkg/gin-swagger"
)

// @title NilOrg认证授权服务
// @version 1.0
// @description NilOrg认证授权服务Api详情.
// @termsOfService https://github.com/nilorg

// @contact.name API Support
// @contact.url https://github.com/nilorg/naas
// @contact.email 862860000@qq.com

// @license.name GNU General Public License v3.0
// @license.url https://github.com/nilorg/naas/blob/master/LICENSE

// @host localhost:8080
// @BasePath /api/v1

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl http://localhost:8080/oauth2/token
// @authorizationurl http://localhost:8080/oauth2/authorize
// @scope.openid 用户openid
// @scope.profile 用户资料
// @scope.email 用户emial
// @scope.phone 用户手机号

// RunHTTP ...
func RunHTTP() {
	store, err := redis.NewStore(10, "tcp", viper.GetString("session.redis.address"), viper.GetString("session.redis.password"), []byte(viper.GetString("session.secret")))
	if err != nil {
		logrus.Errorf("redis.NewStore Error:", err)
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
	if viper.GetBool("swagger.enabled") {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.OAuth(&ginSwagger.OAuthConfig{
				ClientId:     viper.GetString("swagger.oauth2.client_id"),
				ClientSecret: viper.GetString("swagger.oauth2.client_secret"),
				Realm:        viper.GetString("swagger.oauth2.realm"),
				AppName:      viper.GetString("swagger.oauth2.app_name"),
			}),
			ginSwagger.OAuth2RedirectURL(viper.GetString("swagger.oauth2.redirect_url")),
		))
	}

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
	if viper.GetBool("server.open.enabled") {
		oidcGroup := r.Group("/open")
		{
			oidcGroup.POST("/users/wx", middleware.OAuth2AuthScopeRequired("wx_create_user"), open.User.CreateUserFromWeixin)
		}
	}
	if viper.GetBool("server.admin.enabled") {
		apiGroup := r.Group("api/v1", middleware.JWTAuthRequired(global.JwtPublicKey, viper.GetString("server.admin.oauth2.client_id")), middleware.CasbinAuthRequired(casbin.Enforcer))
		{
			apiGroup.GET("/users", api.User.ListByPaged)
			apiGroup.GET("/users/:user_id", api.User.GetOne)
			apiGroup.POST("/users", api.User.Create)
			apiGroup.PUT("/users/:user_id", api.User.Update)
			apiGroup.DELETE("/users/:user_id", api.User.Delete)
			apiGroup.GET("/users/:user_id/organizations", api.User.GetOrganizationList)
			apiGroup.PUT("/users/:user_id/organizations", api.User.UpdateOrganization)
			apiGroup.GET("/users/:user_id/organizations/:organization_id/roles", api.User.GetRoleList)
			apiGroup.PUT("/users/:user_id/roles", api.User.UpdateRole)

			apiGroup.GET("/roles", api.Role.QueryChildren())
			apiGroup.GET("/roles/:role_code", api.Role.GetOne)
			apiGroup.POST("/roles/:role_code/resource_web_route/:resource_web_route_id", api.Role.AddResourceWebRoute)
			apiGroup.POST("/roles", api.Role.Create)
			apiGroup.PUT("/roles", api.Role.Update)
			apiGroup.DELETE("/roles", api.Role.Delete)

			apiGroup.GET("/resource/servers", api.Resource.ListServerByPaged)
			apiGroup.GET("/resource/servers/:resource_server_id", api.Resource.GetServerOne)
			apiGroup.POST("/resource/servers", api.Resource.CreateServer)
			apiGroup.PUT("/resource/servers/:resource_server_id", api.Resource.UpdateServer)
			apiGroup.DELETE("/resource/servers", api.Resource.DeleteServer)
			apiGroup.GET("/resource/web_routes", api.Resource.ListWebRoutePaged)
			apiGroup.POST("/resource/web_routes", api.Resource.AddWebRoute)
			apiGroup.PUT("/resource/web_routes/:resource_web_route_id", api.Resource.UpdateWebRoute)
			apiGroup.DELETE("/resource/web_routes", api.Resource.DeleteWebRoute)
			apiGroup.GET("/resource/web_routes/:resource_web_route_id", api.Resource.GetWebRouteOne)

			apiGroup.POST("/files", api.File.Upload)

			apiGroup.GET("/oauth2/clients", api.OAuth2.ClientListByPaged)
			apiGroup.GET("/oauth2/clients/:client_id", api.OAuth2.GetClient)
			apiGroup.POST("/oauth2/clients", api.OAuth2.CreateClient)
			apiGroup.PUT("/oauth2/clients/:client_id", api.OAuth2.UpdateClient)

			apiGroup.GET("/oauth2/clients/:client_id/scopes", api.OAuth2.GetClientScopes)
			apiGroup.PUT("/oauth2/clients/:client_id/scopes", api.OAuth2.UpdateClientScopes)
			apiGroup.GET("/oauth2/scopes", api.OAuth2.ScopeQueryChildren())
			apiGroup.GET("/oauth2/scopes/:scop_code", api.OAuth2.GetScopeOne)
			apiGroup.PUT("/oauth2/scopes/:scop_code", api.OAuth2.EditScope)

			apiGroup.GET("/organizations", api.Organization.QueryChildren())
			apiGroup.GET("/organizations/:org_id", api.Organization.GetOne)
			apiGroup.POST("/organizations", api.Organization.Create)
			apiGroup.PUT("/organizations/:org_id", api.Organization.Update)
			apiGroup.DELETE("/organizations/:org_id", api.Organization.Delete)

			apiGroup.GET("/casbin/resource/:resource_server_id/web_routes", api.Casbin.ListResourceWebRoutes)
			apiGroup.GET("/casbin/role/:role_code/resource/:resource_server_id/web_routes", api.Casbin.ListRoleResourceWebRoutes)
			apiGroup.GET("/common/select", api.Common.SelectQueryChildren())
		}
	}
	r.Run(fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.oauth2.port"))) // listen and serve on 0.0.0.0:8080
}

func grpcOAuth2(inCtx context.Context) (outCtx context.Context, err error) {
	var clientID string
	clientID, err = grpc_auth.AuthFromMD(inCtx, "client_id")
	if err != nil {
		return
	}
	var clientSecret string
	clientSecret, err = grpc_auth.AuthFromMD(inCtx, "client_secret")
	if err != nil {
		return
	}
	logrus.Debugf("clientID: %s, clientSecret: %s", clientID, clientSecret)
	return
}

// RunGRpc 运行Grpc
func RunGRpc() {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	gRPCServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_auth.StreamServerInterceptor(grpcOAuth2),
				grpc_logrus.StreamServerInterceptor(logrusEntry),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(grpcOAuth2),
				grpc_logrus.UnaryServerInterceptor(logrusEntry),
			),
		),
	)
	service.RegisterGrpc(gRPCServer)
	// 在gRPC服务器上注册反射服务。
	reflection.Register(gRPCServer)
	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.grpc.port"))
	logrus.Infof("%s grpc server listen: %s", "naas", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Errorf("net.Listen Error: %s", err)
		return
	}
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			logrus.Infof("%s grpc server failed to serve: %v", "naas", err)
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
	logrus.Infof("启动GRpcGateway: %s", addr)
	go func() {
		if srvErr := srv.ListenAndServe(); srvErr != nil {
			log.Printf("%s gateway server listen: %v\n", viper.GetString("server.name"), srvErr)
		}
	}()
}
