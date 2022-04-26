package wechat

// JSONer json接口
type JSONer interface {
	JSON() string
}

// XMLer xml接口
type XMLer interface {
	XML() string
}
