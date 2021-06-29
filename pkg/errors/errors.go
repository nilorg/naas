package errors

import "errors"

var (
	// Is reports whether any error in err's chain matches target.
	// 使用系统标准接口
	Is = errors.Is
	// New returns an error that formats as the given text.
	// Each call to New returns a distinct error value even if the text is identical.
	// 使用系统标准接口
	New = errors.New
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user_notfound")
	// ErrUsernameOrPassword 用户名或密码错误
	ErrUsernameOrPassword = errors.New("incorrect_username_or_password")
	// ErrUsernameExist 用户名存在
	ErrUsernameExist = errors.New("username_exist")
	// ErrWxUnionIDExist 微信unionid存在
	ErrWxUnionIDExist = errors.New("wx_unionid_exist")
	// ErrWxOpenIDExist 微信openid存在
	ErrWxOpenIDExist = errors.New("wx_openid_exist")
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
	// ErrRoleCodeExist 角色code存在
	ErrRoleCodeExist = errors.New("role_code_exist")
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role_not_found")
	// ErrRoleCurrentAndParentEqual 上级角色和当前角色相同
	ErrRoleCurrentAndParentEqual = errors.New("role_parent_equal")
	// ErrRoleParentNotExist 上级角色不存在
	ErrRoleParentNotExist = errors.New("role_parent_not_exist")
	// ErrRoleUpdate 修改角色错误
	ErrRoleUpdate = errors.New("role_update")
	// ErrOAuth2CleintNotFound 客户端不存在
	ErrOAuth2CleintNotFound = errors.New("oauth2_client_notfound")
	// ErrUserExistThird 用户第三方存在
	ErrUserExistThird = errors.New("user_exist_third")
	// ErrThirdExistUser 第三方用户存在
	ErrThirdExistUser = errors.New("third_exist_user")
	// ErrThirdUserNotFound 第三方用户不存在
	ErrThirdUserNotFound = errors.New("user_third_not_found")
)
