package key

import "fmt"

// WrapOAuth2Code 包装 oauth2 code key
func WrapOAuth2Code(code string) string {
	return fmt.Sprintf("oauth2:code:%s", code)
}

// WrapOAuth2DeviceCode 包装 oauth2 device code key
func WrapOAuth2DeviceCode(code string) string {
	return fmt.Sprintf("oauth2:device:code:%s", code)
}

// WrapOAuth2UserCode 包装 oauth2 user code key
func WrapOAuth2UserCode(code string) string {
	return fmt.Sprintf("oauth2:user:code:%s", code)
}

const (
	// SessionAccount 当前账户
	SessionAccount = "session_account"
	// SessionDeviceUserCode 设备用户code
	SessionDeviceUserCode = "device_user_code"
)
