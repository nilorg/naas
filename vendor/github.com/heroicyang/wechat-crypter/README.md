# wechat-crypter
微信开放平台加解密库 (Golang)

## Usage

```bash
$ go get github.com/heroicyang/wechat-crypter
```

```go
import "github.com/heroicyang/wechat-crypter"

token := "RMNlACHlV5ThzfRlVS4D4"
corpID := "wx5823bf96d3bd56c7"
encodingAESKey := "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

msgCrypter, _ := crypter.NewMessageCrypter(token, encodingAESKey, corpID)
msgDecrypt, corpID, err := msgCrypter.Decrypt("msgEncrypt")
```

## Doc
[http://godoc.org/github.com/heroicyang/wechat-crypter](http://godoc.org/github.com/heroicyang/wechat-crypter)
