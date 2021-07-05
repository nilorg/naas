package weixin

import (
	wechat "github.com/nilorg/go-wechat"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// WechatClient 微信客户端
	WechatClient wechat.Clienter
	// WechatClientConfig 微信客户端配置文件
	WechatClientConfig wechat.Configer = &cf{}
)

type cf struct{}

func (c *cf) AppID() string {
	return viper.GetString("weixin.kfpt.app_id")
}
func (c *cf) AppSecret() string {
	return viper.GetString("weixin.kfpt.app_secret")
}

// Init 初始化全局变量
func Init() {
	WechatClient = wechat.NewClientFromRedis(
		wechat.ClientFromRedisOptionRedisClient(store.RedisClient),
		wechat.ClientFromRedisOptionAccessTokenKey("github.com/nilorg/go-wechat/access_token"),
		wechat.ClientFromRedisOptionJsAPITicketKey("github.com/nilorg/go-wechat/js_api_ticket"),
	)
	logrus.Debugf("微信Token初始化：%s", WechatClient.GetAccessToken())
}
