package key

import "fmt"

// WrapOAuth2Code 包装 oauth2 code key
func WrapOAuth2Code(code string) string {
	return fmt.Sprintf("oauth2_code_%s", code)
}

const (
	// SessionAccount 当前账户
	SessionAccount = "session_account"
)
