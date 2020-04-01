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
