package random

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

// Number 随机数字
func Number(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString(strconv.Itoa(r.Intn(10)))
	}
	return buffer.String()
}

// AZaz09 随机字符串
func AZaz09(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	strs := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var buffer bytes.Buffer
	for i := 0; i < l; i++ {
		buffer.WriteByte(strs[r.Intn(len(strs))])
	}
	return buffer.String()
}
