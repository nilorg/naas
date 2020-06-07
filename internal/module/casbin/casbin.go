package casbin

import (
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/spf13/viper"
)

var (
	// Enforcer ...
	Enforcer *casbin.SyncedEnforcer
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
	Enforcer, err = casbin.NewSyncedEnforcer(viper.GetString("casbin.config"), adapter)
	if err != nil {
		panic(err)
	}
	// Load the policy from DB.
	Enforcer.LoadPolicy()
	Enforcer.StartAutoLoadPolicy(time.Second * 10)

	// Enforcer.AddFunction("MyDomKeyMatch2", MyDomKeyMatch2Func)
	// Enforcer.AddFunction("MyRegexMatch", MyRegexMatchFunc)

	// aaa, ww := Enforcer.Enforce("eve", "tenant1", "/alice_data/111", "POST")
	// logger.Debugln(aaa, ww)

	// Enforcer.AddRoleForUserInDomain("role:global_admin", "role:reader", "domain1")
	// Enforcer.AddRoleForUserInDomain("role:global_admin", "role:writer", "domain1")

	// Enforcer.AddRoleForUserInDomain("alice", "role:global_admin", "domain1")

	// Enforcer.AddPolicy("role:reader", "domain1", "data1", "read")
	// Enforcer.AddPolicy("role:writer", "domain1", "data1", "write")

	// roles, _ := Enforcer.GetImplicitRolesForUser("alice", "domain1")
	// logger.Debugln("roles:", roles)
	// for _, role := range roles {
	// 	check, checkErr := Enforcer.Enforce(role, "domain1", "data1", "read")
	// 	logger.Debugln(role, "check:", check, checkErr)
	// }
}

// validate the variadic parameter size and type as string
func validateVariadicArgs(expectedLen int, args ...interface{}) error {
	if len(args) != expectedLen {
		return fmt.Errorf("Expected %d arguments, but got %d", expectedLen, len(args))
	}

	for _, p := range args {
		_, ok := p.(string)
		if !ok {
			return errors.New("Argument must be a string")
		}
	}

	return nil
}

// MyDomKeyMatch2Func 定义域KeyMatch2
func MyDomKeyMatch2Func(args ...interface{}) (interface{}, error) {
	if err := validateVariadicArgs(4, args...); err != nil {
		return false, fmt.Errorf("%s: %s", "keyMatch2", err)
	}
	name1 := args[0].(string)
	name2 := args[1].(string)
	dom1 := args[2].(string)
	dom2 := args[3].(string)
	return dom1 == dom2 && (bool)(util.KeyMatch2(name1, name2)), nil
}

// MyRegexMatchFunc 定义域RegexMatch
func MyRegexMatchFunc(args ...interface{}) (interface{}, error) {
	if err := validateVariadicArgs(4, args...); err != nil {
		return false, fmt.Errorf("%s: %s", "keyMatch2", err)
	}
	name1 := args[0].(string)
	name2 := args[1].(string)
	dom1 := args[2].(string)
	dom2 := args[3].(string)
	return dom1 == dom2 && (bool)(util.RegexMatch(name1, name2)), nil
}
