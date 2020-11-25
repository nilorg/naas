package gt3

const (
	ClientTypeWeb     = "web"     // web（pc浏览器）
	ClientTypeH5      = "h5"      // h5（手机浏览器，包括webview）
	ClientTypeNative  = "native"  // native（原生app）
	ClientTypeUnknown = "unknown" // unknown（未知）
)

const (
	// GeetestChallenge 极验二次验证表单传参字段 chllenge
	GeetestChallenge string = "geetest_challenge"
	// GeetestValidate 极验二次验证表单传参字段 validate
	GeetestValidate string = "geetest_validate"
	// GeetestSeccode 极验二次验证表单传参字段 seccode
	GeetestSeccode string = "geetest_seccode"
	// GeetestServerStatusSessionKey 极验验证API服务状态Session Key
	GeetestServerStatusSessionKey string = "gt_server_status"
)
