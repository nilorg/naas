package tools

// InStringSplit 在 string split存在
func InStringSplit(v string, values []string) bool {
	for _, sv := range values {
		if sv == v {
			return true
		}
	}
	return false
}
