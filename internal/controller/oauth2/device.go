package oauth2

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/geetest/gt3"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/geetest"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DeviceActivatePage ...
func DeviceActivatePage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	geetestEnabled := viper.GetBool("geetest.enabled")
	userCode := ctx.Query("user_code")
	ctx.HTML(http.StatusOK, "device_activate.tmpl", gin.H{
		"error":           errMsg,
		"geetest_enabled": geetestEnabled,
		"user_code":       userCode,
	})
}

// DeviceActivate device activate post
func DeviceActivate(ctx *gin.Context) {
	userCode := ctx.PostForm("userCode")
	if userCode == "" || strings.TrimSpace(userCode) == "" {
		SetErrorMessage(ctx, "请输入您的一次性代码")
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	session := sessions.Default(ctx)
	var err error
	if viper.GetBool("geetest.enabled") {
		challenge := ctx.PostForm(gt3.GeetestChallenge)
		seccode := ctx.PostForm(gt3.GeetestSeccode)
		status := session.Get("gt_server_status_for_device")
		if status == nil {
			SetErrorMessage(ctx, "未找到极验验证授权信息")
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		} else if status.(int) == 1 {
			var res *gt3.ValidateResponse
			res, err = geetest.GeetestClient.Validate(challenge, seccode)
			if err != nil {
				SetErrorMessage(ctx, err.Error())
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
				return
			}
			if res.Seccode == "false" {
				SetErrorMessage(ctx, "验证码未通过")
				ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			}
		} else {
			err = SetErrorMessage(ctx, "极验验证授权状态错误")
			ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
			return
		}
	}
	// var userValue map[string]string
	// userValue, err = store.RedisClient.HGetAll(context.Background(), key.WrapOAuth2UserCode(userCode)).Result()
	// if err != nil {
	// 	logrus.Errorln(err)
	// 	SetErrorMessage(ctx, "用户代码无效或过期")
	// 	ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
	// 	return
	// }
	session.Set(key.SessionDeviceUserCode, userCode)
	err = session.Save()
	if err != nil {
		logrus.Errorf("Login-Success-session.Save: %s", err)
		err = SetErrorMessage(ctx, err.Error())
		if err != nil {
			logrus.Errorf("SetErrorMessage: %s", err)
		}
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
		return
	}
	ctx.Redirect(http.StatusFound, "/oauth2/device/confirmation")
}

// DeviceConfirmationPage ...
func DeviceConfirmationPage(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userCode := session.Get(key.SessionDeviceUserCode).(string)
	// var (
	// 	err       error
	// 	userValue map[string]string
	// )
	errMsg := GetErrorMessage(ctx)
	if errMsg != "" {
		goto END
	}
	// userValue, err = store.RedisClient.HGetAll(context.Background(), key.WrapOAuth2UserCode(userCode)).Result()
	// if err != nil {
	// 	logrus.Errorln(err)
	// 	return
	// }
END:
	// TODO：获取客户端信息
	ctx.HTML(http.StatusOK, "device_confirmation.tmpl", gin.H{
		"error":     errMsg,
		"user_code": userCode,
	})
}

// DeviceConfirmation ...
func DeviceConfirmation(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userCode := session.Get(key.SessionDeviceUserCode).(string)
	userCodeKey := key.WrapOAuth2UserCode(userCode)
	var (
		userValue      map[string]string
		err            error
		currentAccount *model.SessionAccount
	)
	defer func() {
		store.RedisClient.Del(ctx, userCodeKey)
	}()
	userValue, err = store.RedisClient.HGetAll(context.Background(), userCodeKey).Result()
	if err != nil {
		logrus.Errorln(err)
		goto END
	}
	currentAccount = session.Get(key.SessionAccount).(*model.SessionAccount)
	err = store.RedisClient.HSet(
		context.Background(), key.WrapOAuth2DeviceCode(userValue["device_code"]),
		"status", "1",
		"open_id", model.ConvertIDToString(currentAccount.UserID),
	).Err()
	if err != nil {
		logrus.Errorln(err)
	}
END:
	if err != nil {
		SetErrorMessage(ctx, "用户代码无效或过期")
		ctx.Redirect(http.StatusFound, ctx.Request.RequestURI)
	} else {
		session.Delete(key.SessionDeviceUserCode)
		session.Save()
		ctx.Redirect(http.StatusFound, "/oauth2/device/success")
	}
}

// DeviceSuccessPage ...
func DeviceSuccessPage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	ctx.HTML(http.StatusOK, "device_success.tmpl", gin.H{
		"error": errMsg,
	})
}

// DeviceErrorPage ...
func DeviceErrorPage(ctx *gin.Context) {
	errMsg := GetErrorMessage(ctx)
	ctx.HTML(http.StatusBadRequest, "device_error.tmpl", gin.H{
		"error": errMsg,
	})
}

// func redirectLogin(ctx *gin.Context, clientID, rURI string) {
// 	redirectURI, _ := url.Parse("/oauth2/login")
// 	redirectURIQuery := url.Values{}
// 	redirectURIQuery.Set("client_id", clientID)
// 	redirectURIQuery.Set("login_redirect_uri", rURI)
// 	redirectURI.RawQuery = redirectURIQuery.Encode()
// 	ctx.Redirect(302, redirectURI.String())
// 	ctx.Abort()
// }
