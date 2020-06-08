package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// Admin ...
	Admin = &admin{}
	// User ...
	User = &user{}
)

// ResultError example...
type ResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "error",
		"error":  err.Error(),
	})
}

func writeData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"data":   data,
	})
}
