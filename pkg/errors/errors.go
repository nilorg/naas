package errors

import "errors"

var (
	// Is reports whether any error in err's chain matches target.
	// 使用系统标准接口
	Is = errors.Is
	// ErrUsernameOrPassword 用户名或密码错误
	ErrUsernameOrPassword = errors.New("incorrect_username_or_password")
	// ErrUsernameExist 用户名存在
	ErrUsernameExist = errors.New("username_exist")
	// ErrWxUnionIDExist 微信unionid存在
	ErrWxUnionIDExist = errors.New("wx_unionid_exist")
	// ErrOrganizationCodeExist 组织code存在
	ErrOrganizationCodeExist = errors.New("org_code_exist")
	// ErrOrganizationParentNotExist 上级组织不存在
	ErrOrganizationParentNotExist = errors.New("org_parent_not_exist")
	// ErrOrganizationCurrentAndParentEqual 上级组织和当前组织相同
	ErrOrganizationCurrentAndParentEqual = errors.New("org_parent_equal")
)
