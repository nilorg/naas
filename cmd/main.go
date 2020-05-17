package main

import (
	"os"
	"runtime"

	_ "github.com/nilorg/naas/docs"
	"github.com/nilorg/naas/internal/controller/oauth2"
	"github.com/nilorg/naas/internal/module"
	"github.com/nilorg/naas/internal/server"
	"github.com/nilorg/pkg/logger"
	"github.com/spf13/viper"
)

func init() {
	// 初始化线程数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Init()
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	configFilename := "configs/config.yaml"
	if v := os.Getenv("NAAS_CONFIG"); v != "" {
		configFilename = v
	}
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logger.Fatalf("Fatal error config file: %s ", err)
	}
	module.Init()
	oauth2.Init()
}

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
// @BasePath /v1
func main() {
	// server.RunGRpc()
	// server.RunGRpcGateway()
	server.RunHTTP()
}
