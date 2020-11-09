package storage

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// UploadStorager 上传接口
type UploadStorager interface {
	Uploader
	MaxMemory() int64
}

// DefaultUploadStorage 默认上传
type DefaultUploadStorage struct {
	*DefaultStorage
}

// MaxMemory 最大上传大小
func (*DefaultUploadStorage) MaxMemory() int64 {
	return defaultMaxMemory
}

// NewDefaultUploadStorage 创建默认存储
func NewDefaultUploadStorage() *DefaultUploadStorage {
	return &DefaultUploadStorage{
		DefaultStorage: NewDefaultStorage(),
	}
}

// UploadFileInfoer 上传file信息接口
type UploadFileInfoer interface {
	// FullName 完整的文件名
	// 包含路径
	FullName() string
	// Filename 文件名
	Filename() string
	// Size 文件大小
	Size() int64
	// Header 获取MIMEHeader
	Header() textproto.MIMEHeader
}

// uploadFileInfo 上传file信息
type uploadFileInfo struct {
	fullName string
	filename string
	size     int64
	header   textproto.MIMEHeader
}

func (ufi *uploadFileInfo) FullName() string {
	return ufi.fullName
}

func (ufi *uploadFileInfo) Filename() string {
	return ufi.filename
}

func (ufi *uploadFileInfo) Size() int64 {
	return ufi.size
}

func (ufi *uploadFileInfo) Header() textproto.MIMEHeader {
	return ufi.header
}

// UploadHandle 上传处理
func UploadHandle(ctx context.Context, r *http.Request, us UploadStorager, name string) (infos []UploadFileInfoer, err error) {
	if us == nil {
		us = NewDefaultUploadStorage()
	}
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(us.MaxMemory())
		if err != nil {
			return nil, err
		}
	}
	for _, mfh := range r.MultipartForm.File[name] {
		var file multipart.File
		file, err = mfh.Open()
		if err != nil {
			return
		}
		var fullName string
		fullName, err = us.Upload(ctx, file, mfh.Filename)
		file.Close()
		if err != nil {
			return
		}
		// https://github.com/golang/go/issues/19501
		// var size int64
		// size, err = GetMultipartFileSize(file)
		// if err != nil {
		// 	return
		// }
		infos = append(infos, &uploadFileInfo{
			fullName: fullName,
			filename: mfh.Filename,
			size:     mfh.Size,
			header:   mfh.Header,
		})
	}
	return
}

// 获取文件大小的接口
type sizer interface {
	Size() int64
}

// 获取文件信息的接口
type stater interface {
	Stat() (os.FileInfo, error)
}

// GetMultipartFileSize 获取上传文件大小
// 相关问题：https://github.com/golang/go/issues/19501
func GetMultipartFileSize(file multipart.File) (size int64, err error) {
	if sizeImpl, ok := file.(sizer); ok {
		size = sizeImpl.Size()
	} else if statImpl, ok := file.(stater); ok {
		var fileInfo os.FileInfo
		fileInfo, err = statImpl.Stat()
		if err != nil {
			return
		}
		size = fileInfo.Size()
	}
	return
}
