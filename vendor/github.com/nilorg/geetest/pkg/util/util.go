package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

// MD5Encode md5 加密
func MD5Encode(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Sha256Encode sha256加密
func Sha256Encode(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HmacSha256Encode hmac-sha256 加密
func HmacSha256Encode(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}
