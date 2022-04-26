package http

// Header Http 请求 头
type Header map[string]string

// Set 设置值
func (h Header) Set(key, value string) {
	h[key] = value
}

// Get 获取值
func (h Header) Get(key string) string {
	return h[key]
}
