package module

import (
	"github.com/nilorg/naas/internal/module/global"
	"github.com/nilorg/naas/internal/module/store"
)

// Init 初始化 module
func Init() {
	store.Init()
	global.Init()
}
