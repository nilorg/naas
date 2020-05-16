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

客户端读取资源，要保证resource owner、client、token和数据库的记录都匹配才行。

# OpenID Connent

https://contoso.auth0.com/.well-known/openid-configuration

```json
{
    "issuer": "https://contoso.auth0.com/",
    "authorization_endpoint": "https://contoso.auth0.com/authorize",
    "token_endpoint": "https://contoso.auth0.com/oauth/token",
    "userinfo_endpoint": "https://contoso.auth0.com/userinfo",
    "mfa_challenge_endpoint": "https://contoso.auth0.com/mfa/challenge",
    "jwks_uri": "https://contoso.auth0.com/.well-known/jwks.json",
    "registration_endpoint": "https://contoso.auth0.com/oidc/register",
    "revocation_endpoint": "https://contoso.auth0.com/oauth/revoke",
    "scopes_supported": [
        "openid",
        "profile",
        "offline_access",
        "name",
        "given_name",
        "family_name",
        "nickname",
        "email",
        "email_verified",
        "picture",
        "created_at",
        "identities",
        "phone",
        "address"
    ],
    "response_types_supported": [
        "code",
        "token",
        "id_token",
        "code token",
        "code id_token",
        "token id_token",
        "code token id_token"
    ],
    "code_challenge_methods_supported": [
        "S256",
        "plain"
    ],
    "response_modes_supported": [
        "query",
        "fragment",
        "form_post"
    ],
    "subject_types_supported": [
        "public"
    ],
    "id_token_signing_alg_values_supported": [
        "HS256",
        "RS256"
    ],
    "token_endpoint_auth_methods_supported": [
        "client_secret_basic",
        "client_secret_post"
    ],
    "claims_supported": [
        "aud",
        "auth_time",
        "created_at",
        "email",
        "email_verified",
        "exp",
        "family_name",
        "given_name",
        "iat",
        "identities",
        "iss",
        "name",
        "nickname",
        "phone_number",
        "picture",
        "sub"
    ],
    "request_uri_parameter_supported": false,
    "device_authorization_endpoint": "https://contoso.auth0.com/oauth/device/code"
}
```

https://accounts.google.com/.well-known/openid-configuration

```json
{
  "issuer": "https://accounts.google.com",
  "authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
  "device_authorization_endpoint": "https://oauth2.googleapis.com/device/code",
  "token_endpoint": "https://oauth2.googleapis.com/token",
  "userinfo_endpoint": "https://openidconnect.googleapis.com/v1/userinfo",
  "revocation_endpoint": "https://oauth2.googleapis.com/revoke",
  "jwks_uri": "https://www.googleapis.com/oauth2/v3/certs",
  "response_types_supported": [
    "code",
    "token",
    "id_token",
    "code token",
    "code id_token",
    "token id_token",
    "code token id_token",
    "none"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "scopes_supported": [
    "openid",
    "email",
    "profile"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_post",
    "client_secret_basic"
  ],
  "claims_supported": [
    "aud",
    "email",
    "email_verified",
    "exp",
    "family_name",
    "given_name",
    "iat",
    "iss",
    "locale",
    "name",
    "picture",
    "sub"
  ],
  "code_challenge_methods_supported": [
    "plain",
    "S256"
  ],
  "grant_types_supported": [
    "authorization_code",
    "refresh_token",
    "urn:ietf:params:oauth:grant-type:device_code",
    "urn:ietf:params:oauth:grant-type:jwt-bearer"
  ]
}
```

内省端点是[《RFC 7662》](https://tools.ietf.org/html/rfc7662)的实现。
