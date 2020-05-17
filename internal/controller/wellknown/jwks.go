package wellknown

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/square/go-jose/v3"
)

// GetJwks ...
func GetJwks(ctx *gin.Context) {
	set := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			global.Jwk,
		},
	}
	ctx.JSON(http.StatusOK, set)
}
