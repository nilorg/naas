package geetest

import (
	"github.com/nilorg/geetest/gt3"
	"github.com/spf13/viper"
)

var (
	// GeetestClient ...
	GeetestClient *gt3.Client
)

// Init ...
func Init() {
	if viper.GetBool("geetest.enabled") {
		initGeetest()
	}
}

func initGeetest() {
	GeetestClient = gt3.NewClient(viper.GetString("geetest.id"), viper.GetString("geetest.key"))
}
