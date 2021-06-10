package casbin

import (
	"fmt"

	"github.com/nilorg/naas/internal/module/casbin"
)

var (
	// MyDomKeyMatch2Func 定义域KeyMatch2
	MyDomKeyMatch2Func = casbin.MyDomKeyMatch2Func
	// MyRegexMatchFunc 定义域RegexMatch
	MyRegexMatchFunc = casbin.MyRegexMatchFunc
)

type PolicyV1 struct {
	Sub    string // 希望访问资源的用户
	Domain string // 域/域租户,这里以资源为单位
	Object string // 要访问的资源
	Action string // 用户对资源执行的操作
}

func FormatPolicyV1(roles []string, resourceID, obj, act string) (policys []*PolicyV1) {
	for _, role := range roles {
		policys = append(policys, &PolicyV1{
			Sub:    fmt.Sprintf("role:%s", role),
			Domain: fmt.Sprintf("resource:%s:route", resourceID),
			Object: obj,
			Action: act,
		})
	}
	return
}
