package third

import (
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/controller/api"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/naas/internal/pkg/contexts"
	"github.com/nilorg/naas/internal/service"
	"github.com/nilorg/naas/pkg/errors"
	"github.com/nilorg/naas/pkg/tools/key"
	"github.com/nilorg/sdk/random"
	"github.com/sirupsen/logrus"
	skip2Qrcode "github.com/skip2/go-qrcode"
)

const (
	QrcodeLoginStatusPending = "pending"
	QrcodeLoginStatusExpired = "expired"
)

type qrcode struct {
}

// GenerateLoginQrCode 生成登录二维码
func (*qrcode) GenerateLoginQrCode(ctx *gin.Context) {
	loginCode := random.AZaz09(32)
	loginCodeKey := key.WrapQrCodeLoginCode(loginCode)
	userCode := random.AZaz09(32)
	userCodekey := key.WrapQrCodeLoginUserCode(userCode)
	expires := time.Now().Add(time.Minute)
	// 把登录Code写入Session
	session := sessions.Default(ctx)
	session.Set(key.SessionQrCodeLoginCode, loginCode)
	err := session.Save()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	// 用户Code和登录Code对应
	err = store.RedisClient.HSet(
		ctx, userCodekey,
		"login_code", loginCode,
		"open_id", "",
	).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	err = store.RedisClient.ExpireAt(ctx, userCodekey, expires).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	// 登录Code中写入状态数据
	err = store.RedisClient.HSet(
		ctx, loginCodeKey,
		"user_code", userCode,
		"status", "0",
	).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	err = store.RedisClient.ExpireAt(ctx, loginCodeKey, expires).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	qrcodeURL := fmt.Sprintf("%s://%s/third/qrcode/confirmation?user_code=%s", scheme, "naas.nilorg.com", userCode)
	var qrcodeBytes []byte
	qrcodeBytes, err = skip2Qrcode.Encode(qrcodeURL, skip2Qrcode.Medium, 256)
	if err != nil {
		api.Writer(ctx, errors.New("生成二维码错误"))
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(qrcodeBytes)
}

// CycleValidationLoginQrCode 循环验证，确认登录登录二维码
func (*qrcode) CycleValidationLoginQrCode(ctx *gin.Context) {
	// 把登录Code写入Session
	session := sessions.Default(ctx)
	loginCode, loginCodeOK := session.Get(key.SessionQrCodeLoginCode).(string)
	if !loginCodeOK || utf8.RuneCountInString(loginCode) != 32 {
		api.Writer(ctx, errors.New("登录Code格式不正确"))
		return
	}
	loginCodeKey := key.WrapQrCodeLoginCode(loginCode)
	parentCtx := contexts.WithGinContext(ctx)
	var err error
	if store.RedisClient.Exists(parentCtx, loginCodeKey).Val() != 1 {
		api.Writer(ctx, gin.H{
			"login": QrcodeLoginStatusExpired,
		})
		return
	}
	var loginValue map[string]string
	if loginValue, err = store.RedisClient.HGetAll(parentCtx, loginCodeKey).Result(); err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	// 检查用户是否同意授权
	if loginValue["status"] != "1" {
		api.Writer(ctx, gin.H{
			"login": QrcodeLoginStatusPending,
		})
		return
	}
	// 授权通过
	// 获取用户Code
	userCode := loginValue["user_code"]
	userCodekey := key.WrapQrCodeLoginUserCode(userCode)
	if store.RedisClient.Exists(parentCtx, userCodekey).Val() != 1 {
		api.Writer(ctx, gin.H{
			"login": QrcodeLoginStatusExpired,
		})
		return
	}
	// 通过用户Code获取用户信息和登录Code
	var userValue map[string]string
	if userValue, err = store.RedisClient.HGetAll(parentCtx, userCodekey).Result(); err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	if userValue["login_code"] != loginCode {
		api.Writer(ctx, errors.New("登录Code不匹配"))
		return
	}
	openID := userValue["open_id"]
	var suser *model.SessionAccount
	suser, err = service.User.LoginForUserID(parentCtx, model.ConvertStringToID(openID))
	if err != nil {
		api.Writer(ctx, err)
		return
	}
	session.Delete(key.SessionQrCodeLoginCode)
	session.Set(key.SessionAccount, suser)
	err = session.Save()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	// 删除验证过程中的信息
	err = store.RedisClient.Del(parentCtx, loginCodeKey, userCodekey).Err()
	if err != nil {
		logrus.Errorln(err)
	}
	api.Writer(ctx, gin.H{
		"login": "success",
	})
}

// ConfirmationLoginQrCode 确认登录二维码
func (*qrcode) ConfirmationLoginQrCode(ctx *gin.Context) {
	userCode := ctx.Query("user_code")
	userCodekey := key.WrapQrCodeLoginUserCode(userCode)
	parentCtx := contexts.WithGinContext(ctx)
	if store.RedisClient.Exists(parentCtx, userCodekey).Val() != 1 {
		api.Writer(ctx, gin.H{
			"login": QrcodeLoginStatusExpired,
		})
		return
	}
	var err error
	var userValue map[string]string
	if userValue, err = store.RedisClient.HGetAll(parentCtx, userCodekey).Result(); err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	loginCode := userValue["login_code"]
	loginCodeKey := key.WrapQrCodeLoginCode(loginCode)
	var loginValue map[string]string
	if loginValue, err = store.RedisClient.HGetAll(parentCtx, loginCodeKey).Result(); err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	if loginValue["user_code"] != userCode {
		api.Writer(ctx, errors.New("用户Code不匹配"))
		return
	}
	err = store.RedisClient.HSet(
		ctx, loginCodeKey,
		"status", "1",
	).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	session := sessions.Default(ctx)
	currentAccount := session.Get(key.SessionAccount).(*model.SessionAccount)
	// TODO:验证成功，模拟登录
	err = store.RedisClient.HSet(
		ctx, userCodekey,
		"open_id", currentAccount.UserID,
	).Err()
	if err != nil {
		logrus.Errorln(err)
		api.Writer(ctx, errors.New("服务器错误"))
		return
	}
	api.Writer(ctx, "登录成功")
}
