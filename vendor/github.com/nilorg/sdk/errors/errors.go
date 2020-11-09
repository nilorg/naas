package errors

import (
	"fmt"
	"regexp"

	"github.com/nilorg/sdk/convert"
)

const (
	grpcErrPattern = "rpc error: code = (?P<rpc_code>.+) desc = (?P<code>[0-9]+)-(?P<msg>.+)"
)

// BusinessError 业务错误
type BusinessError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error
func (err *BusinessError) Error() string {
	return fmt.Sprintf("%d-%s", err.Code, err.Msg)
}

// New ...
func New(code int, msg string) error {
	return &BusinessError{
		Code: code,
		Msg:  msg,
	}
}

// NewBusinessError ...
func NewBusinessError(code int, msg string) *BusinessError {
	return &BusinessError{
		Code: code,
		Msg:  msg,
	}
}

// convertMap 转Map
func convertMap(err error, pattern string) map[string]string {
	r := regexp.MustCompile(pattern)
	match := r.FindStringSubmatch(err.Error())
	paramsMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

// FormatGRpcError 格式化
func FormatGRpcError(err error) (berr *BusinessError) {
	if err == nil {
		berr = nil
		return
	}
	errMap := convertMap(err, grpcErrPattern)
	_, rpcCodeOk := errMap["rpc_code"]
	if !rpcCodeOk {
		return
	}
	code, codeOk := errMap["code"]
	if !codeOk {
		return
	}
	msg, msgOk := errMap["msg"]
	if !msgOk {
		return
	}
	berr = &BusinessError{
		Code: convert.ToInt(code),
		Msg:  msg,
	}
	return
}
