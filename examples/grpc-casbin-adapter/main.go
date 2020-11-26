package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	naasAdapter "github.com/nilorg/naas/pkg/casbin/adapter"
	"github.com/nilorg/ngrpc"
	"github.com/nilorg/sdk/signal"
)

var (
	// Enforcer ...
	Enforcer *casbin.SyncedEnforcer
)

func main() {
	ctx := context.Background()
	var (
		err error
	)
	grpcClient := ngrpc.NewClient("localhost:9000")
	adapter := naasAdapter.NewAdapter(ctx, grpcClient.GetConn())
	adapter.ResourceID = "1"
	adapter.ResourceSecret = "test"
	Enforcer, err = casbin.NewSyncedEnforcer("../../configs/rbac_model.conf", adapter)
	if err != nil {
		panic(err)
	}
	// Load the policy from naas.
	Enforcer.LoadPolicy()
	Enforcer.StartAutoLoadPolicy(time.Minute)

	Enforcer.AddFunction("MyDomKeyMatch2", MyDomKeyMatch2Func)
	Enforcer.AddFunction("MyRegexMatch", MyRegexMatchFunc)
	// 验证
	aaa, ww := Enforcer.Enforce("role:root", "resource:1:route", "/alice_data/111", "POST")
	log.Println(aaa, ww)
	signal.AwaitExit()
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
	return dom1 == dom2 && util.KeyMatch2(name1, name2), nil
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
	return dom1 == dom2 && util.RegexMatch(name1, name2), nil
}
