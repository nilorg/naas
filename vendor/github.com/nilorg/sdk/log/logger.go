package log

// Logger logger
type Logger interface {
	// Debugf 测试
	Debugf(format string, args ...interface{})
	// Debugln 测试
	Debugln(args ...interface{})
	// Infof 信息
	Infof(format string, args ...interface{})
	// Infoln 消息
	Infoln(args ...interface{})
	// Warnf 警告
	Warnf(format string, args ...interface{})
	// Warnln 警告
	Warnln(args ...interface{})
	// Warningf 警告
	Warningf(format string, args ...interface{})
	// Warningln 警告
	Warningln(args ...interface{})
	// Errorf 错误
	Errorf(format string, args ...interface{})
	// Errorln 错误
	Errorln(args ...interface{})
}
