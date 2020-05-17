#!/bin/bash
# 第一步：生成私钥
# genra	生成RSA私钥
# -des3	des3算法
# -out server.key 生成的私钥文件名
# 2048 私钥长度
# pwd：123456
openssl genrsa -des3 -out ca/server.pass.key 2048

# 第二步：去除私钥中的密码
# 注意：有密码的私钥是server.pass.key，没有密码的私钥是server.key
openssl rsa -in ca/server.pass.key -out ca/server.key

# 第三步：生成CSR(证书签名请求)
# req 生成证书签名请求
# -new 新生成
# -key 私钥文件
# -out 生成的CSR文件
# -subj 生成CSR证书的参数
openssl req -new -key ca/server.key -out ca/server.csr -subj "/C=CN/ST=Shandong/L=Jinan/O=devops/OU=devops/CN=nilorg.com"

# 第四步：生成自签名SSL证书
# -days 证书有效期
openssl x509 -req -days 365 -in ca/server.csr -signkey ca/server.key -out ca/server.crt

# X.509证书包含三个文件：key，csr，crt。
# key是服务器上的私钥文件，用于对发送给客户端数据的加密，以及对从客户端接收到数据的解密
# csr是证书签名请求文件，用于提交给证书颁发机构（CA）对证书签名
# crt是由证书颁发机构（CA）签名后的证书，或者是开发者自签名的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息
# 备注：在密码学中，X.509是一个标准，规范了公开秘钥认证、证书吊销列表、授权凭证、凭证路径验证算法等。
