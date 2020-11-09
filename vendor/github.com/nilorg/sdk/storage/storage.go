package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/nilorg/sdk/mime"
)

// Uploader 上传
type Uploader interface {
	Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error)
}

// Downloader 下载
type Downloader interface {
	Download(ctx context.Context, write io.Writer, filename string) (info DownloadFileInfoer, err error)
}

// Remover 删除
type Remover interface {
	Remove(ctx context.Context, filename string) (err error)
}

// Storager 存储
type Storager interface {
	Uploader
	Downloader
	Remover
}

// DefaultStorage 默认存储
type DefaultStorage struct {
	BasePath string
}

// NewDefaultStorage 创建默认存储
func NewDefaultStorage() *DefaultStorage {
	return &DefaultStorage{
		BasePath: "./",
	}
}

// Upload 上传
func (ds *DefaultStorage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
	if rename, ok := FromRenameContext(ctx); ok {
		filename = rename(filename)
	}
	fullName = filepath.Join(ds.BasePath, filename)
	dir := filepath.Dir(fullName)
	_, dirErr := os.Stat(dir)
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return
			}
		} else {
			err = dirErr
			return
		}
	}
	var dist *os.File
	dist, err = os.Create(fullName)
	if err != nil {
		return
	}
	defer dist.Close()
	_, err = io.Copy(dist, read)
	return
}

// Download 下载
func (ds *DefaultStorage) Download(ctx context.Context, dist io.Writer, filename string) (info DownloadFileInfoer, err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	file, err := os.Open(fullName)
	if err != nil {
		return
	}
	defer file.Close()

	md := Metadata{}
	if mimeType, exist := mime.Lookup(filepath.Ext(filename)); exist {
		md.Set("Content-Type", mimeType)
	}
	var (
		downloadFilename      string
		downloadFilenameExist bool
	)
	if downloadFilename, downloadFilenameExist = FromDownloadFilenameContext(ctx); !downloadFilenameExist {
		downloadFilename = filepath.Base(filename)
	}
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return
	}
	info = &downloadFileInfo{
		size:     fileInfo.Size(),
		metadata: md,
		filename: downloadFilename,
	}
	if downloadBefore, downloadBeforeExist := FromDownloadBeforeContext(ctx); downloadBeforeExist {
		downloadBefore(info)
	}
	_, err = io.Copy(dist, file)
	if err != nil {
		info = nil
	}
	return
}

// Remove 删除
func (ds *DefaultStorage) Remove(_ context.Context, filename string) (err error) {
	fullName := filepath.Join(ds.BasePath, filename)
	err = os.Remove(fullName)
	return
}

// DownloadFileInfoer 下载file信息接口
type DownloadFileInfoer interface {
	// Size 文件大小
	Size() int64
	// Filename 文件名
	Filename() string
	// Metadata 获取元数据
	Metadata() Metadata
}

type downloadFileInfo struct {
	size     int64
	filename string
	metadata Metadata
}

func (dfi *downloadFileInfo) Size() int64 {
	return dfi.size
}

func (dfi *downloadFileInfo) Filename() string {
	return dfi.filename
}

func (dfi *downloadFileInfo) Metadata() Metadata {
	return dfi.metadata
}
