package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/naas/internal/dao"
	"github.com/nilorg/naas/internal/module"
	"github.com/nilorg/naas/internal/server"
	"github.com/nilorg/naas/internal/service"
	"github.com/spf13/viper"
)

func init() {
	// 初始化线程数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	configFilename := "configs/config.yaml"
	if v := os.Getenv("NAAS_CONFIG"); v != "" {
		configFilename = v
	}
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s\n", err)
	}
	viper.WatchConfig()

	module.Init()
	dao.Init()
	oauth2.Init()

	if viper.GetBool("casbin.init.enabled") {
		service.Casbin.InitLoadAllPolicy()
	}
}

func main() {

	// 监控系统信号和创建 Context 现在一步搞定
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// 在收到信号的时候，会自动触发 ctx 的 Done ，这个 stop 是不再捕获注册的信号的意思，算是一种释放资源。
	defer stop()
	grpcEnable := os.Getenv("GRPC_ENABLE")
	grpcGatewayEnable := os.Getenv("GRPC_GATEWAY_ENABLE")
	httpEnable := os.Getenv("HTTP_ENABLE")
	if strings.EqualFold(grpcEnable, "true") {
		server.RunGRpc()
	}
	if strings.EqualFold(grpcGatewayEnable, "true") {
		server.RunGRpcGateway()
	}
	if strings.EqualFold(httpEnable, "true") {
		server.RunHTTP()
	}
	<-ctx.Done()
}
