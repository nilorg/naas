package wechat

import (
	"encoding/xml"
	"math/rand"
	"time"

	crypter "github.com/heroicyang/wechat-crypter"
	"github.com/spf13/viper"
)

const encodingAESKey = "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

var (
	msgCrypter, _ = crypter.NewMessageCrypter(viper.GetString("wechat.service.token"), encodingAESKey, viper.GetString("wechat.service.corpID"))
)

// Decrypt 解密
func Decrypt(value string) ([]byte, string) {
	msgDecrypt, corpID, _ := msgCrypter.Decrypt(value)
	return msgDecrypt, corpID
}

// encryptReply 回包加密
type encryptReply struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA    `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	TimeStamp    int64    `xml:"TimeStamp"`
	Nonce        int      `xml:"Nonce"`
}

// XML ...
func (e *encryptReply) XML() string {
	bytes, _ := xml.Marshal(e)
	return string(bytes)
}

// newEncryptReply 创建回包加密
func newEncryptReply(data string) *encryptReply {
	return &encryptReply{
		Encrypt:      CDATA(data),
		MsgSignature: encodingAESKey,
		TimeStamp:    time.Now().Unix(),
		Nonce:        rand.Intn(100),
	}
}

// Encrypt 加密
func Encrypt(value string) string {
	msgEncrypt, _ := msgCrypter.Encrypt(value)
	return newEncryptReply(msgEncrypt).XML()
}
