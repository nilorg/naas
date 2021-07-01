package oauth

import (
	"encoding/json"

	wechat "github.com/nilorg/go-wechat"
	"github.com/nilorg/go-wechat/lang"
)

// OAuth 权限
type OAuth struct {
	config wechat.Configer
}

// NewOAuth ...
func NewOAuth(conf wechat.Configer) *OAuth {
	return &OAuth{
		config: conf,
	}
}

// GetAccessToken 获取 access_token
// 通过code换取网页授权access_token
func (o *OAuth) GetAccessToken(code string) (*AccessTokenReply, error) {
	result, err := wechat.Get("https://api.weixin.qq.com/sns/oauth2/access_token", map[string]string{
		"appid":      o.config.AppID(),
		"secret":     o.config.AppSecret(),
		"code":       code,
		"grant_type": "authorization_code",
	})
	if err != nil {
		return nil, err
	}
	reply := new(AccessTokenReply)
	json.Unmarshal(result, reply)
	return reply, nil
}

// RefreshToken 刷新access_token
// 由于access_token拥有较短的有效期，当access_token超时后，可以使用refresh_token进行刷新，refresh_token有效期为30天，当refresh_token失效之后，需要用户重新授权。
func (o *OAuth) RefreshToken(accessToken string) (*RefreshTokenReply, error) {
	result, err := wechat.Get("https://api.weixin.qq.com/sns/oauth2/refresh_token", map[string]string{
		"appid":         o.config.AppID(),
		"grant_type":    "refresh_token",
		"refresh_token": accessToken,
	})
	if err != nil {
		return nil, err
	}
	reply := new(RefreshTokenReply)
	json.Unmarshal(result, reply)
	return reply, nil
}

// GetUserInfo 拉取用户信息
// 如果网页授权作用域为snsapi_userinfo，则此时开发者可以通过access_token和openid拉取用户信息了。
func (o *OAuth) GetUserInfo(accessToken, openID string) (*UserInfoReply, error) {
	result, err := wechat.Get("https://api.weixin.qq.com/sns/userinfo", map[string]string{
		"access_token": accessToken,
		"openid":       openID,
		"lang":         lang.ZH_CN,
	})
	if err != nil {
		return nil, err
	}
	reply := new(UserInfoReply)
	json.Unmarshal(result, reply)
	return reply, nil
}

// CheckAccessToken 检查Token
func (o *OAuth) CheckAccessToken(accessToken, openID string) (bool, error) {
	_, err := wechat.Get("https://api.weixin.qq.com/sns/auth", map[string]string{
		"access_token": accessToken,
		"openid":       openID,
	})
	if err != nil {
		if err.Error() == "ok" {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// Code2Session 小程序登录凭证校验
// 通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程。
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (o *OAuth) Code2Session(code string) (*Code2SessionReponse, error) {
	result, err := wechat.Get("https://api.weixin.qq.com/sns/jscode2session", map[string]string{
		"appid":      o.config.AppID(),
		"secret":     o.config.AppSecret(),
		"js_code":    code,
		"grant_type": "authorization_code",
	})
	if err != nil {
		return nil, err
	}
	resp := new(Code2SessionReponse)
	json.Unmarshal(result, resp)
	return resp, nil
}
