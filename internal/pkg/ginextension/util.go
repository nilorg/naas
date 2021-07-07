package ginextension

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// IsMicroMessenger 判断是否是微信发起的请求
func IsMicroMessenger(ctx *gin.Context) bool {
	userAgent := ctx.GetHeader("User-Agent")
	userAgent = strings.ToLower(userAgent)
	return strings.Index(userAgent, "micromessenger") != -1
}
