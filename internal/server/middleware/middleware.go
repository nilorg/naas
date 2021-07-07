package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Header 头处理
func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", "naas(https://github.com/nilorg/naas)")
		c.Next()
	}
}

func parseAuth(auth string) (token string, ok bool) {
	const prefix = "Bearer "
	// Case insensitive prefix match. See Issue 22736.
	if auth == "" || len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	return auth[len(prefix):], true
}
