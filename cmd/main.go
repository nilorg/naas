package main

import (
	"os"
	"runtime"

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
	viper.WatchConfig()

	module.Init()
	oauth2.Init()
}

func main() {
	// server.RunGRpc()
	// server.RunGRpcGateway()
	server.RunHTTP()
}
