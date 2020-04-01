package errors

import "errors"

var (
	// ErrUsernameOrPassword 用户名或密码错误
	ErrUsernameOrPassword = errors.New("incorrect_username_or_password")
)
