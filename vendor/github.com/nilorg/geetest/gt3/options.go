package gt3

// ClientOptions 客户端可选参数
type ClientOptions struct {
	APIURL      string
	RegisterURL string
	ValidteURL  string
	HTTPTimeout int
	Version     string
	ClientType  string
	IPAddress   string
	JSONFormat  string
}

// DefaultClientOptions 默认客户端配置
var DefaultClientOptions = ClientOptions{
	APIURL:      "http://api.geetest.com",
	RegisterURL: "register.php",
	ValidteURL:  "validate.php",
	HTTPTimeout: 5, // 单位：秒
	Version:     "golang-gin:3.1.0",
	ClientType:  ClientTypeUnknown,
	JSONFormat:  "1",
}

// newOptions 创建可选参数
func newOptions(opts ...Option) ClientOptions {
	opt := DefaultClientOptions
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Option 为可选参数赋值的函数
type Option func(*ClientOptions)

// OptionAPIURL ...
func OptionAPIURL(v string) Option {
	return func(o *ClientOptions) {
		o.APIURL = v
	}
}

// OptionRegisterURL ...
func OptionRegisterURL(v string) Option {
	return func(o *ClientOptions) {
		o.RegisterURL = v
	}
}
