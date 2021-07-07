package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/pkg/tools/key"
)

// ThirdAuthRequired 第三方绑定验证
func ThirdAuthRequired(ctx *gin.Context) {
	session := sessions.Default(ctx)
	third := session.Get(key.SessionThird)
	if third == nil {
		ctx.String(http.StatusBadRequest, "未找到第三方绑定信息")
		ctx.Abort()
	} else {
		ctx.Next()
	}
}
