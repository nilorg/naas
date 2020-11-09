package oauth2

import "fmt"

// Logger logger
type Logger interface {
	// Debugf 测试
	Debugf(format string, args ...interface{})
	// Debugln 测试
	Debugln(args ...interface{})
	// Errorf 错误
	Errorf(format string, args ...interface{})
	// Errorln 错误
	Errorln(args ...interface{})
}

// DefaultLogger ...
type DefaultLogger struct{}

// Debugf ...
func (*DefaultLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf("OAuth2 [DEBUG] "+format+"\n", args...)
}

// Debugln ...
func (*DefaultLogger) Debugln(args ...interface{}) {
	fmt.Println("OAuth2 [DEBUG] ", args)
}

// Errorf ...
func (*DefaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("OAuth2 [ERROR] "+format+"\n", args...)
}

// Errorln ...
func (*DefaultLogger) Errorln(args ...interface{}) {
	fmt.Println("OAuth2 [ERROR] ", args)
}
