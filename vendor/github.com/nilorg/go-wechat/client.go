package wechat

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/nilorg/sdk/http"
)

var (
	// AppID 应用Key
	AppID = ""
	// AppSecret 秘密
	AppSecret = ""
)

// Parameter 参数
type Parameter map[string]string

// Execute 执行
func Execute(url string, param Parameter) (json *simplejson.Json, err error) {
	err = checkConfig()
	if err != nil {
		return
	}
	param["appid"] = AppID
	param["secret"] = AppSecret

	result, err := http.Get(url, param)
	if err != nil {
		return
	}
	json, err = simplejson.NewJson([]byte(result))
	if err != nil {
		return
	}
	if errmsg, ok := json.CheckGet("errcode"); ok {
		bytes, _ := errmsg.Encode()
		err = errors.New(string(bytes))
		return
	}
	return
}

func checkConfig() error {
	if AppID == "" {
		return errors.New("AppID 不能为空")
	}
	if AppSecret == "" {
		return errors.New("AppSecret 不能为空")
	}
	return nil
}

// Configer 微信客户端接口
type Configer interface {
	AppID() string
	AppSecret() string
}

// Clienter 微信客户端接口
type Clienter interface {
	GetAccessToken() string
	GetJsAPITicket() string
}

// Client 客户端
type Client struct {
	appID       string
	appSecret   string
	accessToken string
	jsAPITicket string
	ticker      *time.Ticker
	rwMutex     *sync.RWMutex // 读写锁
}

// startTicker 启动打点任务
func (c *Client) startTicker() {
	go func() { // 异步
		for range c.ticker.C {
			log.Println("刷新AccessToken和JsAPITicket")
			c.refreshAccessToken() // 刷新AccessToken
			c.refreshJsAPITicket()
		}
	}()
}

// NewClient 创建客户端
func NewClient(appID, appSecret string) *Client {
	client := &Client{
		appID:     appID,
		appSecret: appSecret,
		ticker:    time.NewTicker(time.Hour), // 1小时执行一次
		rwMutex:   &sync.RWMutex{},
	}
	client.refreshAccessToken() // 刷新AccessToken
	client.refreshJsAPITicket()
	client.startTicker()
	return client
}

// AccessTokenReply ...
type AccessTokenReply struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken 获取 accessToken
func (c *Client) GetAccessToken() string {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.accessToken
}

// refreshAccessToken ...
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
func (c *Client) refreshAccessToken() {
	result, err := Get("https://api.weixin.qq.com/cgi-bin/token", map[string]string{
		"appid":      c.AppID(),
		"secret":     c.AppSecret(),
		"grant_type": "client_credential",
	})
	if err != nil {
		log.Printf("刷新AccessToken错误：%v\n", err)
		// c.accessToken = ""
		return
	}
	reply := new(AccessTokenReply)
	json.Unmarshal(result, reply)
	c.rwMutex.Lock()
	c.accessToken = reply.AccessToken
	c.rwMutex.Unlock()
}

// JsAPITicketReply ...
type JsAPITicketReply struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// GetJsAPITicket 获取 js api ticket
func (c *Client) GetJsAPITicket() string {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.jsAPITicket
}

// refreshJsAPITicket ...
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
func (c *Client) refreshJsAPITicket() {
	result, err := Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket", map[string]string{
		"access_token": c.GetAccessToken(),
		"type":         "jsapi",
	})
	if err != nil {
		log.Printf("刷新Ticket错误：%v", err)
		// c.jsAPITicket = ""
		return
	}
	reply := new(JsAPITicketReply)
	json.Unmarshal(result, reply)

	c.rwMutex.Lock()
	c.jsAPITicket = reply.Ticket
	c.rwMutex.Unlock()
}

// AppID ...
func (c *Client) AppID() string {
	return c.appID
}

// AppSecret ...
func (c *Client) AppSecret() string {
	return c.appSecret
}
