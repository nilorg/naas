package contexts

import (
	"context"
	"errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/model"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/pkg/storage"
	"github.com/spf13/viper"
)

type userIDKey struct{}

var (
	// ErrUserIDNotFound ...
	ErrUserIDNotFound = errors.New("上下文中，用户ID不存在")
)

// NewUserIDContext ...
func NewUserIDContext(ctx context.Context, userID model.ID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// FromUserIDContext ...
func FromUserIDContext(ctx context.Context) (userID model.ID, err error) {
	var ok bool
	userID, ok = ctx.Value(userIDKey{}).(model.ID)
	if !ok {
		err = ErrUserIDNotFound
	}
	return
}

// WithContext ...
func WithContext(ctx context.Context) context.Context {
	if viper.GetString("storage.type") == "oss" {
		ctx = storage.NewBucketNameContext(ctx, viper.GetString("storage.oss.bucket"))
	}
	ctx = store.NewDBContext(ctx)
	return ctx
}

// WithGinContext ...
func WithGinContext(ctx *gin.Context) context.Context {
	parent := context.Background()
	parent = WithContext(parent)
	claims := jwt.ExtractClaims(ctx)
	if id, ok := claims["user_id"].(string); ok {
		parent = NewUserIDContext(parent, model.ConvertStringToID(id))
	}
	return parent
}
