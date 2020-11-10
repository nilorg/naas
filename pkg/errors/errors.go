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
	// ErrOrganizationNotFound 组织不存在
	ErrOrganizationNotFound = errors.New("org_not_found")
	// ErrOrganizationCreate 创建组织错误
	ErrOrganizationCreate = errors.New("org_create")
	// ErrOrganizationUpdate 修改组织错误
	ErrOrganizationUpdate = errors.New("org_update")
	// ErrOrganizationCurrentAndParentEqual 上级组织和当前组织相同
	ErrOrganizationCurrentAndParentEqual = errors.New("org_parent_equal")
	// ErrResourceCreate 创建资源错误
	ErrResourceCreate = errors.New("resource_create")
	// ErrResourceUpdate 修改资源错误
	ErrResourceUpdate = errors.New("resource_update")
	// ErrResourceNotFound 资源服务器不存在
	ErrResourceNotFound = errors.New("resource_not_found")
)
