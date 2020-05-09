# naas
Authentication authorization server（认证授权服务器）

```bash
http://localhost:8080/oauth2/authorize?client_id=oauth2_client&redirect_uri=http://localhost/callback&response_type=code&state=somestate&scope=read_write
```

# Develop
## Nilorg Gateway

```bash
# client_id=1000
# redirect_uri=http://localhost:8000/auth/callback
http://localhost:8080/oauth2/authorize?client_id=1000&redirect_uri=http://dev.wohuitao.vip:8000/auth/callback&response_type=code&state=somestate&scope=read_write
```

# naas-token-server
在使用AuthorizationCode（授权码模式）进行授权的时候，服务器下发的`code`变量需要开发者单独编写API进行向OAuth2Server获取Token、刷新Token。

使用`naas-token-server`能够帮您减轻不必要的接口开发，帮助您快速构建OAuth2应用API。