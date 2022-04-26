package weixin

import (
	wechat "github.com/nilorg/go-wechat"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// KfptWechatClient 开放平台微信客户端
	KfptWechatClient wechat.Clienter
	// KfptWechatClientConfig 开放平台微信客户端配置文件
	KfptWechatClientConfig wechat.Configer = &kfptCf{}
	// FwhWechatClient 服务号微信客户端
	FwhWechatClient wechat.Clienter
	// FwhWechatClientConfig 服务号微信客户端配置文件
	FwhWechatClientConfig wechat.Configer = &fwhCf{}
)

type kfptCf struct{}

func (c *kfptCf) AppID() string {
	return viper.GetString("weixin.kfpt.app_id")
}
func (c *kfptCf) AppSecret() string {
	return viper.GetString("weixin.kfpt.app_secret")
}

type fwhCf struct{}

func (c *fwhCf) AppID() string {
	return viper.GetString("weixin.fwh.app_id")
}
func (c *fwhCf) AppSecret() string {
	return viper.GetString("weixin.fwh.app_secret")
}

// Init 初始化全局变量
func Init() {
	KfptWechatClient = wechat.NewClientFromRedis(
		wechat.ClientFromRedisOptionRedisClient(store.RedisClient),
		wechat.ClientFromRedisOptionAccessTokenKey("naas:wx:kfpt:access_token"),
		wechat.ClientFromRedisOptionJsAPITicketKey("naas:wx:kfpt:js_api_ticket"),
	)
	logrus.Debugf("微信开放平台Token初始化：%s", KfptWechatClient.GetAccessToken())
	FwhWechatClient = wechat.NewClientFromRedis(
		wechat.ClientFromRedisOptionRedisClient(store.RedisClient),
		wechat.ClientFromRedisOptionAccessTokenKey("github.com/nilorg/go-wechat/access_token"),
		wechat.ClientFromRedisOptionJsAPITicketKey("github.com/nilorg/go-wechat/js_api_ticket"),
	)
	logrus.Debugf("微信服务号Token初始化：%s", FwhWechatClient.GetAccessToken())
}
