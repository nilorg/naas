package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sdkErrors "github.com/nilorg/sdk/errors"
	grpcStatus "google.golang.org/grpc/status"
)

var (
	// User ...
	User = &user{}
	// OAuth2 ...
	OAuth2 = &oauth2{}
	// Role ...
	Role = &role{}
	// Resource ...
	Resource = &resource{}
	// File ...
	File = &file{}
	// Organization ...
	Organization = &organization{}
	// Common ...
	Common = &common{}
)

// Result example...
type Result struct {
	Status string       `json:"status"`
	Data   interface{}  `json:"data,omitempty"`
	Error  *ResultError `json:"error,omitempty"`
}

// ResultError ...
type ResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewResult 创建返回结果
func NewResult(v interface{}) (result Result) {
	switch v.(type) {
	case *Result:
		result = v.(Result)
	case *grpcStatus.Status:
		gs := v.(grpcStatus.Status)
		result.Status = "error"
		result.Error = &ResultError{
			Code:    int(gs.Code()),
			Message: gs.Message(),
		}
	case *sdkErrors.BusinessError:
		be := v.(*sdkErrors.BusinessError)
		result.Status = "error"
		result.Error = &ResultError{
			Code:    be.Code,
			Message: be.Msg,
		}
	case error:
		result.Status = "error"
		result.Error = &ResultError{
			Code:    -1,
			Message: v.(error).Error(),
		}
	default:
		result.Status = "ok"
		result.Data = v
	}
	return
}

// Writer ...
func Writer(ctx *gin.Context, v interface{}) {
	ctx.JSON(http.StatusOK, NewResult(v))
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

// QueryChildren ...
func QueryChildren(v map[string]gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		q := ctx.Query("q")
		if q == "" {
			ctx.String(400, "未知的查询类型")
			return
		}
		h, ok := v[q]
		if !ok {
			ctx.String(400, "未找到查询类型")
			return
		}
		h(ctx)
	}
}
