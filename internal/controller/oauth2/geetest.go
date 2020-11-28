package oauth2

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/geetest/gt3"
	"github.com/nilorg/naas/internal/module/geetest"
	"github.com/nilorg/oauth2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GeetestRegister 验证初始化接口，GET请求
func GeetestRegister(ctx *gin.Context) {
	var (
		registerRes *gt3.RegisterResponse
		err         error
		res         gt3.RegisterResponseForWeb
	)
	registerRes, err = geetest.GeetestClient.Register("md5")
	if err != nil {
		res.Success = 0
		logrus.Errorln(err)
	} else {
		res.Success = 1
		res.Gt = viper.GetString("geetest.id")
		res.Challenge = geetest.GeetestClient.BuildChallenge(registerRes.Challenge, "md5")
		res.NewCaptcha = "true"
	}
	q := ctx.Query("q")
	session := sessions.Default(ctx)
	if q == oauth2.DeviceCodeKey {
		session.Set("gt_server_status_for_device", res.Success)
	} else {
		session.Set(gt3.GeetestServerStatusSessionKey, res.Success)
	}
	session.Save()

	ctx.JSON(http.StatusOK, res)
}
