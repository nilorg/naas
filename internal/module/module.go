package module

import (
	"github.com/nilorg/naas/internal/module/casbin"
	"github.com/nilorg/naas/internal/module/geetest"
	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/naas/internal/module/logger"
	"github.com/nilorg/naas/internal/module/store"
)

// Init 初始化 module
func Init() {
	logger.Init()
	store.Init()
	global.Init()
	casbin.Init()
	geetest.Init()
}
