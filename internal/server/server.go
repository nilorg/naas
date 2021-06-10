package server

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	daprd "github.com/dapr/go-sdk/service/grpc"
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
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/global"
	internalService "github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	// swagger doc ...
	_ "github.com/nilorg/naas/docs"
	"github.com/nilorg/naas/internal/server/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"

	//ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/nilorg/naas/internal/pkg/contexts"
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
		SameSite: http.SameSiteStrictMode,
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
		oidcGroup := r.Group("/oidc")
		{
			if viper.GetBool("server.oidc.userinfo_endpoint_enabled") {
				oidcGroup.GET("/userinfo", middleware.OAuth2AuthUserinfoRequired(global.JwtPublicKey), oidc.GetUserinfo)
			}
		}
	}

	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/login", oauth2.LoginPage)
		oauth2Group.POST("/login", oauth2.Login)
		oauth2Group.GET("/authorize", middleware.OAuth2AuthRequired, oauth2.AuthorizePage)
		oauth2Group.POST("/authorize", middleware.OAuth2AuthRequired, oauth2.Authorize)
		oauth2Group.POST("/token", oauth2.Token)
		if viper.GetBool("server.oauth2.device_authorization_endpoint_enabled") {
			oauth2Group.GET("/device/activate", oauth2.DeviceActivatePage)
			oauth2Group.POST("/device/activate", oauth2.DeviceActivate)
			oauth2Group.GET("/device/confirmation", middleware.OAuth2AuthDeviceRequired, oauth2.DeviceConfirmationPage)
			oauth2Group.POST("/device/confirmation", middleware.OAuth2AuthDeviceRequired, oauth2.DeviceConfirmation)
			oauth2Group.GET("/device/success", oauth2.DeviceSuccessPage)
			oauth2Group.GET("/device/error", oauth2.DeviceErrorPage)

			oauth2Group.POST("/device/code", oauth2.DeviceCode)
		}
		if viper.GetBool("server.oauth2.introspection_endpoint_enabled") {
			oauth2Group.POST("/introspect", oauth2.TokenIntrospection)
		}
		if viper.GetBool("server.oauth2.revocation_endpoint_enabled") {
			oauth2Group.POST("/revoke", oauth2.TokenRevoke)
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
			apiGroup.POST("/roles", api.Role.Create)
			apiGroup.PUT("/roles", api.Role.Update)
			apiGroup.DELETE("/roles", api.Role.Delete)

			apiGroup.GET("/resource/servers", api.Resource.ListServerByPaged)
			apiGroup.GET("/resource/servers/:resource_server_id", api.Resource.GetServerOne)
			apiGroup.POST("/resource/servers", api.Resource.CreateServer)
			apiGroup.PUT("/resource/servers/:resource_server_id", api.Resource.UpdateServer)
			apiGroup.DELETE("/resource/servers", api.Resource.DeleteServer)

			apiGroup.GET("/resource/routes", api.Resource.ListRoutePaged)
			apiGroup.POST("/resource/routes", api.Resource.AddRoute)
			apiGroup.PUT("/resource/routes/:resource_route_id", api.Resource.UpdateRoute)
			apiGroup.DELETE("/resource/routes", api.Resource.DeleteRoute)
			apiGroup.GET("/resource/routes/:resource_route_id", api.Resource.GetRouteOne)
			apiGroup.GET("/resource/menus", api.Resource.ListMenuPaged)
			apiGroup.POST("/resource/menus", api.Resource.AddMenu)
			apiGroup.PUT("/resource/menus/:resource_menu_id", api.Resource.UpdateMenu)
			apiGroup.DELETE("/resource/menus", api.Resource.DeleteMenu)
			apiGroup.GET("/resource/menus/:resource_menu_id", api.Resource.GetMenuOne)

			apiGroup.GET("/resource/actions", api.Resource.ListActionPaged)
			apiGroup.POST("/resource/actions", api.Resource.AddAction)
			apiGroup.PUT("/resource/actions/:resource_action_id", api.Resource.UpdateAction)
			apiGroup.DELETE("/resource/actions", api.Resource.DeleteAction)
			apiGroup.GET("/resource/actions/:resource_action_id", api.Resource.GetActionOne)

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
			apiGroup.GET("/organizations/:organization_id", api.Organization.GetOne)
			apiGroup.POST("/organizations", api.Organization.Create)
			apiGroup.PUT("/organizations/:organization_id", api.Organization.Update)
			apiGroup.DELETE("/organizations/:organization_id", api.Organization.Delete)

			apiGroup.GET("/casbin/resource/:resource_server_id/routes", api.Casbin.ListResourceRoutes)
			apiGroup.GET("/casbin/role/:role_code/resource/:resource_server_id/routes", api.Casbin.ListRoleResourceRoutes)
			apiGroup.PUT("/casbin/role/:role_code/resource_routes", api.Casbin.AddResourceRoute)

			apiGroup.GET("/casbin/resource/:resource_server_id/menus", api.Casbin.ListResourceMenus)
			apiGroup.GET("/casbin/role/:role_code/resource/:resource_server_id/menus", api.Casbin.ListRoleResourceMenus)
			apiGroup.PUT("/casbin/role/:role_code/resource_menus", api.Casbin.AddResourceMenu)

			apiGroup.GET("/casbin/resource/:resource_server_id/actions", api.Casbin.ListResourceActions)
			apiGroup.GET("/casbin/role/:role_code/resource/:resource_server_id/actions", api.Casbin.ListRoleResourceActions)
			apiGroup.PUT("/casbin/role/:role_code/resource_actions", api.Casbin.AddResourceAction)

			apiGroup.GET("/common/select", api.Common.SelectQueryChildren())
			apiGroup.GET("/common/tree", api.Common.TreeQueryChildren())
		}
	}
	if viper.GetBool("geetest.enabled") {
		geetestGroup := r.Group("/geetest")
		{
			geetestGroup.GET("/register", oauth2.GeetestRegister)
		}
	}
	go func() {
		addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.oauth2.port"))
		if runErr := r.Run(addr); runErr != nil {
			serverName := viper.GetString("server.name")
			logrus.Fatalf("%s http server listen %s: %v\n", serverName, addr, runErr)
		}
	}()
}

func grpcResourceAuth(inCtx context.Context) (outCtx context.Context, err error) {
	var basic string
	basic, err = grpc_auth.AuthFromMD(inCtx, "basic")
	if err != nil {
		return
	}
	var basicBytes []byte
	basicBytes, err = base64.StdEncoding.DecodeString(basic)
	if err != nil {
		return
	}
	ss := strings.SplitN(string(basicBytes), ":", 2)
	if len(ss) != 2 {
		err = errors.New("Invalid basic format")
		return
	}
	logrus.Debugf("resourceID: %s, resourceSecret: %s", ss[0], ss[1])
	err = internalService.Resource.AuthServer(inCtx, model.ConvertStringToID(ss[0]), ss[1])
	if err != nil {
		return
	}
	md, ok := metadata.FromIncomingContext(inCtx)
	if !ok {
		return
	}
	md.Set("resource_id", ss[0])
	md.Set("resource_secret", ss[1])
	outCtx = metadata.NewOutgoingContext(inCtx, md)
	return
}

// RunGRpc 运行Grpc
func RunGRpc() {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	gRPCServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_auth.StreamServerInterceptor(grpcResourceAuth),
				grpc_logrus.StreamServerInterceptor(logrusEntry),
				proto.StreamServerInterceptor(contexts.WithContext),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(grpcResourceAuth),
				grpc_logrus.UnaryServerInterceptor(logrusEntry),
				proto.UnaryServerInterceptor(contexts.WithContext),
			),
		),
	)
	service.RegisterGrpc(gRPCServer)
	// 在gRPC服务器上注册反射服务。
	reflection.Register(gRPCServer)
	addr := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("server.grpc.port"))
	serverName := viper.GetString("server.name")
	logrus.Infof("%s grpc server listen: %s", serverName, addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Errorf("net.Listen Error: %s", err)
		return
	}
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			logrus.Fatalf("%s grpc server failed to serve: %v", serverName, err)
		}
	}()
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
	serverName := viper.GetString("server.name")
	logrus.Infof("%s grpc-gateway server listen: %s", serverName, addr)
	go func() {
		if srvErr := srv.ListenAndServe(); srvErr != nil {
			logrus.Fatalf("%s gateway server listen: %v\n", serverName, srvErr)
		}
	}()
}

// RunDapr 运行Dapr
func RunDapr() {
	addr := getEnvVar("DAPR_ADDRESS", ":5001")
	// create serving server
	daprdService, err := daprd.NewService(addr)
	if err != nil {
		log.Fatalf("dapr new service error: %v", err)
	}
	if err := service.RegisterDapr(daprdService); err != nil {
		log.Fatalf("dapr register method error: %v", err)
	}
	serverName := viper.GetString("server.name")
	logrus.Infof("%s dapr server listen: %s", serverName, addr)
	go func() {
		// start the server to handle incoming events
		if err := daprdService.Start(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	address := flag.String("address", fallbackValue, "service address")
	flag.Parse()
	return *address
}
