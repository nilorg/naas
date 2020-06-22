package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/service"
)

type role struct {
}

func (*role) Recursive(ctx *gin.Context) {
	roles := service.Role.Recursive()
	ctx.JSON(http.StatusOK, roles)
}
