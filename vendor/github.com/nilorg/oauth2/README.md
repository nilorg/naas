# oauth2

# Usage
```bash
go get -u github.com/nilorg/oauth2
```
# Import
```bash
import "github.com/nilorg/oauth2"
```

# 例子

[oauth2-server](https://github.com/nilorg/oauth2-server)

[server/client](https://github.com/nilorg/oauth2/tree/master/examples)

# 文档参考
1. [《理解OAuth 2.0》阮一峰](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)
2. [《RFC 6749》](https://tools.ietf.org/html/rfc6749) | [《RFC 6749》](http://www.rfcreader.com/#rfc6749)
3. [《OAuth 2.0 Device Authorization Grant(RFC8628)》](https://tools.ietf.org/html/rfc8628)
4. [《OAuth 2.0 Token Introspection(RFC7662)》](https://tools.ietf.org/html/rfc7662)
5. [《OAuth 2.0 Token Revocation(RFC7009)》](https://tools.ietf.org/html/rfc7009)


### AuthorizationCode
授权码模式（authorization code）是功能最完整、流程最严密的授权模式。

它的特点就是通过客户端的后台服务器，与"服务提供商"的认证服务器进行互动。
### Implicit
简化模式（implicit grant type）不通过第三方应用程序的服务器，直接在浏览器中向认证服务器申请令牌，跳过了"授权码"这个步骤，因此得名。

所有步骤在浏览器中完成，令牌对访问者是可见的，且客户端不需要认证。
### ResourceOwnerPasswordCredentials
密码模式（Resource Owner Password Credentials Grant）中，用户向客户端提供自己的用户名和密码。

客户端使用这些信息，向"服务商提供商"索要授权。

在这种模式中，用户必须把自己的密码给客户端，但是客户端不得储存密码。

这通常用在用户对客户端高度信任的情况下，比如客户端是操作系统的一部分，或者由一个著名公司出品。

而认证服务器只有在其他授权模式无法执行的情况下，才能考虑使用这种模式。
### ClientCredentials
客户端模式（Client Credentials Grant）指客户端以自己的名义，而不是以用户的名义，向"服务提供商"进行认证。

严格地说，客户端模式并不属于OAuth框架所要解决的问题。

在这种模式中，用户直接向客户端注册，客户端以自己的名义要求"服务提供商"提供服务，其实不存在授权问题。

### DeviceCode
设备模式（Device Code）

### TokenIntrospection
内省端点（Token Introspection）
### TokenRevocation
Token销毁端点（Token Revocation）

# Server

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
)

var (
	clients = map[string]string{
		"oauth2_client": "password",
	}
)

func main() {
	srv := oauth2.NewServer()
	srv.VerifyClient = func(basic *oauth2.ClientBasic) (err error) {
		pwd, ok := clients[basic.ID]
		if !ok {
			err = oauth2.ErrInvalidClient
			return
		}
		basic = &oauth2.ClientBasic{
			ID:     basic.ID,
			Secret: pwd,
		}
		return
	}
	srv.VerifyClientID = func(clientID string) (err error) {
		_, ok := clients[clientID]
		if !ok {
			err = oauth2.ErrInvalidClient
		}
		return
	}
	srv.VerifyCode = func(code, clientID, redirectURI string) (value *oauth2.CodeValue, err error) {
		//err = oauth2.ErrUnauthorizedClient
		// 查询缓存/数据库中的code信息
		value = &oauth2.CodeValue{
			ClientID:    clientID,
			RedirectURI: redirectURI,
			Scope:       []string{"a", "b", "c"},
		}
		return
	}
	srv.GenerateCode = func(clientID, openID, redirectURI string, scope []string) (code string, err error) {
		code = oauth2.RandomCode()
		return
	}
	srv.VerifyRedirectURI = func(clientID, redirectURI string) (err error) {
		fmt.Println(clientID)
		fmt.Println(redirectURI)
		// err = oauth2.ErrInvalidRedirectURI
		return
	}

	srv.VerifyPassword = func(username, password string) (openID string, err error) {
		if username != "a" || password != "b" {
			err = oauth2.ErrUnauthorizedClient
			return
		}
		openID = "xxxx"
		return
	}

	srv.VerifyScope = func(scopes []string, clientID string) (err error) {
		// err = oauth2.ErrInvalidScope
		return
	}

	srv.GenerateAccessToken = oauth2.NewDefaultGenerateAccessToken([]byte("xxxxx"))
	srv.RefreshAccessToken = oauth2.NewDefaultRefreshAccessToken([]byte("xxxxx"))
	srv.ParseAccessToken = oauth2.NewDefaultParseAccessToken([]byte("xxxxx"))

	srv.GenerateDeviceAuthorization = func(issuer, verificationURI, clientID, scope string) (resp *oauth2.DeviceAuthorizationResponse, err error) {
		resp = &oauth2.DeviceAuthorizationResponse{
			DeviceCode:            oauth2.RandomCode(),
			UserCode:              oauth2.RandomUserCode(),
			VerificationURI:       verificationURI,
			VerificationURIQrcode: "",
			ExpiresIn:             0,
			Interval:              5,
		}
		return
	}

	srv.VerifyDeviceCode = func(deviceCode, clientID string) (value *oauth2.DeviceCodeValue, err error) {
		// err = oauth2.ErrAuthorizationPending
		return
	}

	srv.Init()

	// =============Http Default=============
	// http.HandleFunc("/authorize", srv.HandleAuthorize)
	// http.HandleFunc("/token", srv.HandleToken)
	// if err := http.ListenAndServe(":8003", srv); err != nil {
	// 	fmt.Printf("%+v\n", err)
	// }

	// =============Gin=============
	r := gin.Default()
	oauth2Group := r.Group("/oauth2")
	{
		oauth2Group.GET("/authorize", func(c *gin.Context) {
			srv.HandleAuthorize(c.Writer, c.Request)
		})
		oauth2Group.POST("/token", func(c *gin.Context) {
			srv.HandleToken(c.Writer, c.Request)
		})
		oauth2Group.POST("/device_authorization", func(c *gin.Context) {
			srv.HandleDeviceAuthorization(c.Writer, c.Request)
		})
	}

	if err := http.ListenAndServe(":8003", r); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
```

# Client

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/oauth2"
	"github.com/nilorg/pkg/logger"
)

var (
	client *oauth2.Client
)

func init()  {
	logger.Init()
	client = oauth2.NewClient("http://localhost:8003", "oauth2_client", "password")
	client.Log = logger.Default()
}
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		//err := client.AuthorizeImplicit(c.Writer, "http://localhost:8080/callback", "test", "aaaaa")
		//if err != nil {
		//	logger.Errorln(err)
		//	return
		//}
		err := client.AuthorizeAuthorizationCode(c.Writer, "http://localhost:8080/callback", "test", "bbbbb")
		if err != nil {
			logger.Errorln(err)
			return
		}
	})
	r.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		token, err := client.TokenAuthorizationCode(code, c.Request.URL.String(), state)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "callback",
				"err":     err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "callback",
				"token":   token,
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
```

# jwt playload
 
 标准中注册的声明 (建议但不强制使用) ：
 
 `iss`: 令牌**颁发者**。它表示**该令牌是由谁创建的**，在好很多OAuth部署中会将它设为授权服务器的URL。该声明是一个字符串
 
 `sub`: 令牌的**主体**。它表示**该令牌是关于谁的**，在很多OAuth部署中会将它设为资源拥有者的唯一标识。在大多数情况下，主题在同一个颁发者的范围内必须是唯一的。该声明是一个字符串
 
 `aud`: 令牌的**受众**。它表示**令牌的接收者**，在很多OAuth部署中，它包含受保护资源的URI或者能够接收该令牌的受保护资源。该声明可以是一个字符串数组，如果只有一个值，也可以是一个不用数组包装的单个字符串
 
 `exp`: 令牌的**过期**时间戳。它表示**令牌将在何时过期**，以便部署应用让令牌自行失效。该声明是一个整数，表示自UNIX新世纪（即格林威治标准时间GMT，1970年1月1日零点）以来的秒数
 
 `nbf`: 令牌的**生效**时的时间戳。它表示**令牌从什么时候开始生效**，以便部署应用可以在令牌生效之前颁发令牌。该声明是一个整数，表示自UNIX新世纪（即格林威治标准时间GMT，1970年1月1日零点）以来的秒数
 
 `iat`: 令牌**颁发时**的时间戳。它表示**令牌是何时被创建的**，它通常是颁发者在生成令牌时的系统时间戳。该声明是一个整数，表示自UNIX新世纪（即格林威治标准时间GMT，1970年1月1日零点）以来的秒数
 
 `jti`: 令牌的**唯一标识符**。该声明的值**在令牌颁发者创建的每个令牌中都是唯一的**，为了防止冲突，它通常是一个密码学随机值。这个值相当于向结构化令牌中加入了一个攻击者无法获得的随机熵组件，有利于防止令牌猜测攻击和重放攻击

 ---
 
 公共的声明 ：
 公共的声明可以添加任何的信息，一般添加用户的相关信息或其他业务需要的必要信息.但不建议添加敏感信息，因为该部分在客户端可解密.
 
 私有的声明 ：
 私有声明是提供者和消费者所共同定义的声明，一般不建议存放敏感信息，因为base64是对称解密的，意味着该部分信息可以归类为明文信息。