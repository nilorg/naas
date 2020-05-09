package main

import (
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/pkg/token"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
)

var (
	// client oauth2客户端
	client *oauth2.Client
)

func init() {
	// 初始化线程数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Init()
}

func main() {
	var (
		oauth2Server       string
		oauth2ClientID     string
		oauth2ClientSecret string
		oauth2RedirectURI  string
	)
	if oauth2Server = os.Getenv("OAUTH2_SERVER"); oauth2Server == "" {
		logger.Errorln("env OAUTH2_SERVER is empty")
		return
	}
	if oauth2ClientID = os.Getenv("OAUTH2_CLIENT_ID"); oauth2ClientID == "" {
		logger.Errorln("env OAUTH2_CLIENT_ID is empty")
		return
	}
	if oauth2ClientSecret = os.Getenv("OAUTH2_CLIENT_SECRET"); oauth2ClientSecret == "" {
		logger.Errorln("env OAUTH2_CLIENT_SECRET is empty")
		return
	}
	if oauth2RedirectURI = os.Getenv("OAUTH2_REDIRECT_URI"); oauth2RedirectURI == "" {
		logger.Errorln("env OAUTH2_REDIRECT_URI is empty")
		return
	}

	client = oauth2.NewClient(oauth2Server, oauth2ClientID, oauth2ClientSecret)

	r := gin.Default()
	r.GET("/auth/token", token.AuthToken(client, oauth2RedirectURI))
	r.GET("/auth/refresh_token", token.AuthRefreshToken(client))
	if err := r.Run(":8081"); err != nil {
		logger.Errorf("gin Run Error: %s", err)
	}
}
