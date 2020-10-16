package contexts

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrContextNotFoundGorm 上下文不存在Gorm错误
	ErrContextNotFoundGorm = errors.New("上下文中没有获取到Gorm")
)

type gormKey struct{}
type gormTranKey struct{}

// FromGormContext 从上下文中获取Gorm
func FromGormContext(ctx context.Context) (db *gorm.DB, err error) {
	var ok bool
	db, ok = ctx.Value(gormTranKey{}).(*gorm.DB)
	if !ok {
		db, ok = ctx.Value(gormKey{}).(*gorm.DB)
		if !ok {
			err = ErrContextNotFoundGorm
		}
		return
	}
	return
}

// CheckGormTranContextExist 检查gorm是否存在
func CheckGormTranContextExist(ctx context.Context) bool {
	_, ok := ctx.Value(gormTranKey{}).(*gorm.DB)
	return ok
}

// NewGormContext 创建Gorm上下文
func NewGormContext(ctx context.Context, gdb *gorm.DB) context.Context {
	return context.WithValue(ctx, gormKey{}, gdb)
}

// NewGormTranContext 创建Gorm事务上下文
func NewGormTranContext(ctx context.Context, gdb *gorm.DB) context.Context {
	return context.WithValue(ctx, gormTranKey{}, gdb)
}
