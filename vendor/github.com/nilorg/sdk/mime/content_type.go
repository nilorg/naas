package mime

import (
	"io"
	"net/http"
)

// GetFileContentType 获取文件流的ContentType
func GetFileContentType(r io.Reader) (contentType string, err error) {
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)
	_, err = r.Read(buffer)
	if err != nil {
		return
	}
	contentType = http.DetectContentType(buffer)
	return
}
