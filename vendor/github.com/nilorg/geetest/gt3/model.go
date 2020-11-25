package gt3

import (
	"errors"
	"strings"
)

// SDK内部与极验服务器交互接口
// https://docs.geetest.com/sensebot/apirefer/api/server#SDK%E5%86%85%E9%83%A8%E4%B8%8E%E6%9E%81%E9%AA%8C%E6%9C%8D%E5%8A%A1%E5%99%A8%E4%BA%A4%E4%BA%92%E6%8E%A5%E5%8F%A3

// RequestComm 通用请求参数
type RequestComm struct {
	UserID     string `json:"user_id"`     // user_id作为终端用户的唯一标识，确定用户的唯一性；作用于提供进阶数据分析服务，可在api1 或 api2 接口传入，不传入也不影响验证服务的使用；若担心用户信息风险，可作预处理(如哈希处理)再提供到极验
	ClientType string `json:"client_type"` // 客户端类型，web（pc浏览器），h5（手机浏览器，包括webview），native（原生app），unknown（未知）
	IPAddress  string `json:"ip_address"`  // 客户端请求SDK服务器的ip地址
	JSONFormat string `json:"json_format"` // json格式化标识
	Sdk        string `json:"sdk"`         // sdk代码版本号
}

// Validation 验证
func (req *RequestComm) Validation() (err error) {
	if req.JSONFormat == "" || strings.TrimSpace(req.JSONFormat) == "" {
		err = errors.New("json格式化标识不能为空")
	} else if req.Sdk == "" || strings.TrimSpace(req.Sdk) == "" {
		err = errors.New("sdk代码版本号不能为空")
	}
	return
}

// RegisterRequest 请求参数
type RegisterRequest struct {
	*RequestComm
	Digestmod string `json:"digestmod"` // 生成唯一标识字符串的签名算法，默认暂支持md5
	Gt        string `json:"gt"`        // 向极验申请的账号id
}

// Validation 验证
func (req *RegisterRequest) Validation() (err error) {
	err = req.RequestComm.Validation()
	if err != nil {
		return
	}
	if req.Digestmod == "" || strings.TrimSpace(req.Digestmod) == "" {
		err = errors.New("生成唯一标识字符串的签名算法不能为空")
	} else if req.Gt == "" || strings.TrimSpace(req.Gt) == "" {
		err = errors.New("向极验申请的账号id不能为空")
	}
	return
}

// RegisterResponse 响应参数
type RegisterResponse struct {
	Challenge string `json:"challenge"` // 生成唯一流水号的参考字符串，为”0”表示传参账号id有误
}

// ValidateRequest 请求参数
type ValidateRequest struct {
	*RequestComm
	Seccode   string `json:"seccode"`   // 核心校验数据
	Challenge string `json:"challenge"` // 流水号，一次完整验证流程的唯一标识
	CaptchaID string `json:"captchaid"` // 向极验申请的账号id
}

// Validation 验证
func (req *ValidateRequest) Validation() (err error) {
	err = req.RequestComm.Validation()
	if err != nil {
		return
	}
	if req.Seccode == "" || strings.TrimSpace(req.Seccode) == "" {
		err = errors.New("核心校验数据不能为空")
	} else if req.Challenge == "" || strings.TrimSpace(req.Challenge) == "" {
		err = errors.New("流水号不能为空")
	} else if req.CaptchaID == "" || strings.TrimSpace(req.CaptchaID) == "" {
		err = errors.New("向极验申请的账号id不能为空")
	}
	return
}

// ValidateResponse 响应参数
type ValidateResponse struct {
	Seccode string `json:"seccode"` // 验证结果标识，为”false”表示验证不通过
}

// RegisterResponseForWeb register response for web
type RegisterResponseForWeb struct {
	Success    int    `json:"success"`     // 流程正常、异常标识；1表示正常，0表示异常、后续走宕机模式
	NewCaptcha string `json:"new_captcha"` // 新版验证码标识，固定不变
	Challenge  string `json:"challenge"`   // 流水号，一次完整验证流程的唯一标识
	Gt         string `json:"gt"`          // 向极验申请的账号id
}
