package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/spf13/viper"
)

var (
	// Enforcer ...
	Enforcer *casbin.Enforcer
)

// Init 初始化casbin
// https://github.com/casbin/gorm-adapter
func Init() {
	var (
		adapter *gormadapter.Adapter
		err     error
	)
	adapter, err = gormadapter.NewAdapterByDB(store.DB)
	if err != nil {
		panic(err)
	}
	Enforcer, err = casbin.NewEnforcer(viper.GetString("casbin.config"), adapter)
	if err != nil {
		panic(err)
	}

	// Load the policy from DB.
	Enforcer.LoadPolicy()

	// Check the permission.
	// Enforcer.Enforce("alice", "data1", "read")

	// Modify the policy.
	// Enforcer.AddPolicy("eve", "/alice_data/*", "GET")
}
