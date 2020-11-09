package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nilorg/sdk/mime"
)

// DownloadStorager 下载接口
type DownloadStorager interface {
	Downloader
	DispositionType() string
}

// DefaultDownloadStorage 默认下载
type DefaultDownloadStorage struct {
	*DefaultStorage
	dispositionType string
}

// DispositionType Disposition Type
func (d *DefaultDownloadStorage) DispositionType() string {
	return d.dispositionType
}

// NewDefaultDownloadStorage 创建默认存储
func NewDefaultDownloadStorage(dispositionType string) *DefaultDownloadStorage {
	return &DefaultDownloadStorage{
		DefaultStorage:  NewDefaultStorage(),
		dispositionType: dispositionType,
	}
}

// DownloadHandle 下载处理
func DownloadHandle(ctx context.Context, resp http.ResponseWriter, ds DownloadStorager, filename string) (err error) {
	if ds == nil {
		ds = NewDefaultDownloadStorage("attachment")
	}

	if _, ok := FromDownloadBeforeContext(ctx); !ok {
		ctx = NewDownloadBeforeContext(ctx, func(info DownloadFileInfoer) {
			md := info.Metadata()
			wr := false
			if md != nil {
				contextType := md.Get("Content-Type")
				if contextType != "" {
					resp.Header().Set("Content-Type", contextType)
					wr = true
				}
			}
			if mimeType, exist := mime.Lookup(info.Filename()); !wr && exist {
				resp.Header().Set("Content-Type", mimeType)
			}
			dispositionType := ds.DispositionType()
			if dispositionType == "" || (dispositionType != "inline" && dispositionType != "attachment") {
				dispositionType = "attachment"
			}
			resp.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", dispositionType, info.Filename()))
		})
	}
	_, err = ds.Download(ctx, resp, filename)
	if err != nil {
		return
	}
	resp.WriteHeader(http.StatusOK)
	return
}
