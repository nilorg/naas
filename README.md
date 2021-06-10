# naas
Authentication authorization server（认证授权服务器）

# 功能
✅ 已实现&emsp;&emsp;♻️ 实现中&emsp;&emsp;❗️待实现

1. ✅ [OAuth2](https://github.com/nilorg/oauth2)
    * ✅ 授权码模式（Authorization Code）
    * ✅ 简化模式（Implicit Grant Type）
    * ✅ 密码模式（Resource Owner Password Credentials Grant）
    * ✅ 客户端模式（Client Credentials Grant）
    * ✅ 设备模式（Device Code）
    * ✅ 内省端点（Token Introspection）
    * ✅ Token销毁端点（Token Revocation）

2. ✅ OpenIDConnent
   * ✅ jwks
   * ✅ openid
   * ✅ openid-configuration
   * ✅ Userinfo（用户信息端点）
3. ♻️ [Swagger](https://github.com/swaggo/gin-swagger)
   * ✅ 集成OAuth2认证
4. ✅ [Casbin](https://github.com/casbin/casbin)
   * ✅ [RBAC](https://casbin.org/docs/zh-CN/rbac-api)
   * ✅ [多租户](https://casbin.org/docs/zh-CN/rbac-with-domains)
   * ✅ [自定义naas-casbin-adapter](./pkg/casbin/adapter/adapter.go)
5. ♻️ [NAAS后台管理系统](https://github.com/nilorg/naas-web)
   * ✅ 集成OAuth2认证
   * ✅ 基础数据
     * ✅ 组织
     * ✅ 角色
     * ✅ 用户
   * ✅ OAuth2
     * ✅ 客户端
     * ✅ 范围
   * ✅ Casbin
     * ✅ 路由-访问策略
     * ✅ 菜单-访问策略
     * ✅ 动作-访问策略
   * ✅ 资源
     * ✅ 资源服务器
     * ✅ 路由
     * ✅ 菜单
     * ✅ 动作
6. ♻️ 其他
     * ❗️手机验证码登录
     * ❗邮箱验证码登录
     * ❗️第三方登录（微信、钉钉）
     * ❗️用户日志记录
     * ❗️注册页面
     * ✅ 极验验证


# 页面展示（以下展示Nilorg任务调度平台对接使用）

1. 登录页面
  
    ![login](./examples/images/login.png)

2. 授权页面

    ![authorize](./examples/images/authorize.png)

# 部署

## 使用Docker
```bash
docker run -d \
-p 8080:8080 -p 5000:5000 -p 9000:9000 \
--name naas \
-v <local path>/naas/configs:/workspace/configs \
-v <local path>/naas/web:/workspace/web \
--link mysql:mysql \
--link redis:redis \
-e HTTP_ENABLE=true \
-e GRPC_ENABLE=true \
-e GRPC_GATEWAY_ENABLE=true \
nilorg/naas:latest
```
## 使用Kubernetes
1. 创建命名空间
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nilorg
```
```bash
kubectl apply -f ./deployments/k8s/namespace.yaml
```
2. 创建配置文件
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: naas
  namespace: nilorg
data:
  config.yaml: |
    <内容和configs/config.yaml相同>
  rbac_model.conf: |
    <内容和configs/rbac_model.conf相同>
```
```bash
kubectl apply -f ./deployments/k8s/config-cm.yaml
```
3. 创建Pod
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: naas
  namespace: nilorg
spec:
  selector:
    matchLabels:
      app: naas
      service: naas
      version: v1
  replicas: 1
  template:
    metadata:
      labels:
        app: naas
        service: naas
        version: v1
    spec:
      restartPolicy: Always
      containers:
        - name: naas
          image: nilorg/naas:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080 # 对应 HTTP_ENABLE
            - containerPort: 5000 # 对应 GRPC_ENABLE
            - containerPort: 9000 # 对应 GRPC_GATEWAY_ENABLE
          env:
            - name: GRPC_ENABLE
              value: "true"
            - name: GRPC_GATEWAY_ENABLE
              value: "true"
            - name: HTTP_ENABLE
              value: "true"
            - name: DAPR_ENABLE # 只有在使用DApr的时候才会用
              value: "false"
          volumeMounts:
            - name: config-cm # 配置文件
              mountPath: /workspace/configs/
      volumes:
        - name: config-cm # 配置文件
          configMap:
            name: naas
```
```bash
kubectl apply -f ./deployments/k8s/pod.yaml
```
4. 创建Service
```yaml
kind: Service
apiVersion: v1
metadata:
  name: naas
  namespace: nilorg
  labels:
    app: naas
spec:
  selector:
    app: naas
    service: naas
    version: v1
  ports: # 根据自己实际需求配置端口
    - name: naas-8080
      port: 8080
      protocol: TCP
      targetPort: 8080
    - name: naas-5000
      port: 5000
      protocol: TCP
      targetPort: 5000
    - name: naas-9000
      port: 9000
      protocol: TCP
      targetPort: 9000
```
```bash
kubectl apply -f ./deployments/k8s/service.yaml
```
5. 创建Ingress（我这里使用的是`traefik`，根据自己的环境进行调整）
```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: naas
  namespace: nilorg
spec:
  entryPoints:
    - web
  routes:
    - kind: Rule
      match: Host(`naas.nilorg.com`)
      services:
        - name: naas
          namespace: nilorg
          port: 8080
```
```bash
kubectl apply -f ./deployments/k8s/traefik.yaml
```
## 配置文件解答
`configs/config.yaml`
```yaml
server:
  name: naas # 服务器名
  oauth2:
    port: 8080 # http服务端口
    issuer: "https://github.com/nilorg/naas"
    device_authorization_endpoint_enabled: true # 设备授权端点
    introspection_endpoint_enabled: true # 内省端点
    revocation_endpoint_enabled: true # Token销毁端点
  grpc:
    port: 9000
    gateway:
      port: 5000
  oidc:
    enabled: true # 是否开启OpenID Connent
    userinfo_endpoint_enabled: true # 根据token获取用户信息
  open:
    enabled: true # 开放API
  admin:
    enabled: true # 管理端
    external: true # 启用外部管理，需要配置外部URL
    external_url: http://naas-admin.nilorg.com
    super_user: "root"
    oauth2: # 管理端的OAuth2Client配置信息
      client_id: 1000

log:
  level: "debug" # panic/fatal/error/warn/info/debug
  report_caller: true

jwt:
  secret: "github.com/nilorg/naas"
  timeout: 20 # Token过期时间 分钟单位
  max_refresh: 10 #Token过期容忍刷新时间
  rsa: # 用于JWT Token生成，使用脚本创建`./scripts/create.sh`私钥和证书
    private: | # 私钥
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEA20St6pqB4LQvqT1Aq2jZPbrkpSiwFeQwiu6AA2eBz3oYveYA
      SCDzl/jXfPsY36b8VahDWmhgB/ie5Ku+R6yXiZcY9SYDiu8sMONwdkhlIL4nP1oC
      97CffWf4vkt4mH7i5/rJWCd/MMLzjSmrMPdUOh9Jd2awNjUZ9QiVTBogZeMo8b5i
      nVBRfRcKAQDZYlo5/VkpaRBTqahh+RoIReX1MHy/LuPMJywPaqHpIh3dlwOvnY6Q
      uFrPo3cF4B7mi/ofTeRX7xzm6z+uxVZGkUHAxgm4VMAYmiP0dLSzyagA5IHUaPHV
      ex8luTSR6DcbINm0bw9skUzI8zYPIGzI/rchSQIDAQABAoIBAQDazaAXOfNcvbHJ
      2jvMUKZn+TXssbt1PO5L1U+dFg7tcVN7PCcP0wIBpumx6AecNtAa0fvUHc+mZKx6
      V/9bGpllTYg0KajjXWPlrTAueHOhxt73UuUfMfsVc0k+66T917Cp+RIui8taZ1AO
      j4QrKsO79Dilk61HipnKcLQ66t9liv4Uf/oxOjfvjaw0+mRDgD2eulTNE+pSIw6L
      uZXduUcpZkYenXCIS+YfRjKMJGHdCiy0bj8887vg0JiqF+mPxGo1UrOMrkWtC4am
      Fht7IMUO5KnfBveL1rMB3ed8LRie9B5EOopRoBZ7PhZ31sqlimYargHGnZwYH8BH
      HzazCGwBAoGBAO8N14JcbqEcs0VpGqyuuBffheu3+6waGt90MhYEMVJsL07qLkIw
      8P4zvPDthXMncrLBC7VJzKkZ7hmww3/qZX5xYjeSVggxG149I1Kncqn9l9BW/Qes
      IEmTUfDE8Js6mQfJVxf7qKDsN9E5N90Oj2j4XZK2ECfaLKbwWfDv3IBBAoGBAOrP
      x/jm9s6Y6KBzxBkXK0jtx2PGM1KxwJFcH9TKgz1A5yue0I1gVdU5Yf3HQowkUGJK
      lT2sUHh1JXUWd2gSrZ5ba6Fc7yITIRUYjAJaW4JKvGtk59QsdRUsHiKsMxmM1GJl
      /uDuZem+EiSA4R9ZZZSHAIfQY2VJD3MLDWVMvt8JAoGAJDebo/NvC1e2zVhMI0dh
      OrSxrHG2Xm+iDKKlB/LgqhUb4b/W/E4/5LNf97x0kGq0lOJsbK3epOv5x8ihBds0
      P0DcWYEBKcKO2+s1U8tsstZpzrWvJh9s0NjR/EFKFqp9DtHxMP/+n0rKdhdOIF6Z
      WZTvUE/nCLKkOzKE3dzpMkECgYAYkkmwyCqHkAS31aVtorkK1qcIz9LLEoK+M0+5
      ar+1BzepnuLgCHay62BPuCxEkgA/aOKZI5EAKfITgJhaMaotag+nQRxdCndpx7nO
      /TmaNsvkyRhhYY2W+5jjs/Vc9Rm8ekPjsc7EWPl5DGuCZk507nOlwq7ECJMvTLbI
      JPHMUQKBgF9O0xzJu7NwR1njqeU1MWdo8nzmb9F2itsYRXmOtC+rjTs3uqWBqlu3
      TE+L0j3o3S6navSHhzzcZLwozW6otHfDcmfFBQG48zbH7YgBVuTnSQyegEpSUHRa
      Pk78NMGbTCMJ65lA96vscXaSk0hF9Y83YY9Jjiju+uwWdnx74khb
      -----END RSA PRIVATE KEY-----
    cert: | # 签名后的证书
      -----BEGIN CERTIFICATE-----
      MIIDSjCCAjICCQDWXqh/wC9VZjANBgkqhkiG9w0BAQUFADBnMQswCQYDVQQGEwJD
      TjERMA8GA1UECAwIU2hhbmRvbmcxDjAMBgNVBAcMBUppbmFuMQ8wDQYDVQQKDAZk
      ZXZvcHMxDzANBgNVBAsMBmRldm9wczETMBEGA1UEAwwKbmlsb3JnLmNvbTAeFw0y
      MDA1MTYxMjA5MjNaFw0yMTA1MTYxMjA5MjNaMGcxCzAJBgNVBAYTAkNOMREwDwYD
      VQQIDAhTaGFuZG9uZzEOMAwGA1UEBwwFSmluYW4xDzANBgNVBAoMBmRldm9wczEP
      MA0GA1UECwwGZGV2b3BzMRMwEQYDVQQDDApuaWxvcmcuY29tMIIBIjANBgkqhkiG
      9w0BAQEFAAOCAQ8AMIIBCgKCAQEA20St6pqB4LQvqT1Aq2jZPbrkpSiwFeQwiu6A
      A2eBz3oYveYASCDzl/jXfPsY36b8VahDWmhgB/ie5Ku+R6yXiZcY9SYDiu8sMONw
      dkhlIL4nP1oC97CffWf4vkt4mH7i5/rJWCd/MMLzjSmrMPdUOh9Jd2awNjUZ9QiV
      TBogZeMo8b5inVBRfRcKAQDZYlo5/VkpaRBTqahh+RoIReX1MHy/LuPMJywPaqHp
      Ih3dlwOvnY6QuFrPo3cF4B7mi/ofTeRX7xzm6z+uxVZGkUHAxgm4VMAYmiP0dLSz
      yagA5IHUaPHVex8luTSR6DcbINm0bw9skUzI8zYPIGzI/rchSQIDAQABMA0GCSqG
      SIb3DQEBBQUAA4IBAQAxCCdWsJjI0BNja2VhW4UjN+E2NiE5YQU0wZWtoPtc//lt
      RziOGrZP82W6uh6BreonBu9JdNOJ0z+FYO957OrCrk6YBoFHe3l38KkQa13Vc4yG
      2I4s1QPwor9rPRLcRQv4rB/ZS42IXXQBaCEHg+RfQ6oOX8E8YVpmRI8i3fBL4Zcf
      KPiaI5i2Ey9p7ncV+7LhZ9+rZvMeA10v1jdXhl0rRphJjN+EyC+pHCu01NAaQKAo
      Cj3vnvAfK8f8dEsZ9hUHLw1olVz0PbdsoUwdvULvVU5weVNyIGFfFMQeoZESrhxr
      B36K98eWEdm2Wc3IY6OL2xj+DaYm8Tuyh9KzL9hU
      -----END CERTIFICATE-----

session:
  name: "naas-session" # session使用的cookie名称
  secret: "github.com/nilorg/naas" # 用于session的加密
  options: # session 配置的可选项
    path: "/"
    domain: "naas.nilorg.com"
    max_age: 86400
    secure: false # 要在HTTPS下开启才可以，HTTP下开启导致Session不可用问题
    http_only: true
  redis: # 用于存储Session的Redis配置信息
    address: "localhost:6379"
    password: ""

mysql: # MySQL数据库
  address: "root:test123@tcp(localhost:3306)/naas?charset=utf8&parseTime=True&loc=Local"
  log: true # 是否打印log

redis: # Redis
  address: "localhost:6379"
  password: ""
  db: 0

swagger: # https://swagger.io
  enabled: true # 是否启用Swagger
  oauth2: # 用于Swagger中的OAuth2配置信息
    client_id: 1000
    client_secret: 22222
    realm:
    app_name: naas-server
    redirect_url: http://naas.nilorg.com/swagger/oauth2-redirect.html # 授权回调地址

casbin: # https://casbin.org
  init:
    enabled: false # 是否初始化Casbin信息，用于项目第一次初始化使用。
  config: configs/rbac_model.conf # casbin配置文件

storage: # 对象存储，目前支持两种方式default和oss，default使用指定文件夹目录进行存储、oss使用阿里云对象存储进行存储
  type: default # default/oss
  default: 
    base_path: ./web/storage
  oss: # 阿里云对象存储配置信息
    endpoint: oss-cn-shanghai.aliyuncs.com
    bucket: xxx
    access:
      key_id: aaaaa
      key_secret: bbbbb
  public_path: http://localhost:8080/storage # 文件前缀地址，用于访问文件使用.oss的方式可以使用外网地址
  max_memory: 20 # 20MB

naas:
  resource: # 用于后端API授权资源使用
    id: 1

geetest: # https://www.geetest.com 极验验证
  enabled: true
  id: "c9c4facd1a6feeb80802222cbb74ca8e" # 可更换为自己的
  key: "f7475f921a41f7ba79ae15e41658627c" # 可更换为自己的
```
`configs/rbac_model.conf` 查看 [Casbin Model语法](https://casbin.org/docs/zh-CN/syntax-for-models)
```conf
# Model语法 https://casbin.org/docs/zh-CN/syntax-for-models
# sub, obj, act 表示经典三元组: 访问实体 (Subject)，访问资源 (Object) 和访问方法 (Action)。
# sub:希望访问资源的用户
# dom:域/域租户 https://casbin.org/docs/zh-CN/rbac-with-domains
# obj:要访问的资源
# act:用户对资源执行的操作

# request_definition:请求定义
[request_definition]
r = sub, dom, obj, act

# policy_definition:策略定义
[policy_definition]
p = sub, dom, obj, act

# role_definition:角色定义
[role_definition]
g = _, _, _

# policy_effect:政策的影响
[policy_effect]
e = some(where (p.eft == allow))

# matchers:匹配器
[matchers]
m = g(r.sub, p.sub, r.dom) == true \
&& MyDomKeyMatch2(r.obj, p.obj, r.dom, p.dom) == true \
&& MyRegexMatch(r.act, p.act, r.dom, p.dom) == true \
|| r.sub == "role:naas_root"
```