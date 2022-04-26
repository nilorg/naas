package storage

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/nilorg/sdk/convert"
	"github.com/nilorg/sdk/mime"
	"github.com/nilorg/sdk/storage"
)

// AliyunOssStorage 阿里云oss存储
type AliyunOssStorage struct {
	bucketNames                 []string
	ossClient                   *oss.Client
	CheckAndCreateBucketEnabled bool
}

// NewAliyunOssStorage 创建阿里云oss存储
func NewAliyunOssStorage(ossClient *oss.Client, initBucket bool, bucketNames []string) (os *AliyunOssStorage, err error) {
	os = &AliyunOssStorage{
		bucketNames:                 bucketNames,
		ossClient:                   ossClient,
		CheckAndCreateBucketEnabled: true,
	}
	if initBucket {
		err = os.initBucket()
		if err != nil {
			os = nil
		}
	}
	return
}

// initBucket 初始化桶
func (ds *AliyunOssStorage) initBucket() (err error) {
	for _, bucketName := range ds.bucketNames {
		err = ds.CheckAndCreateBucket(bucketName)
		if err != nil {
			return
		}
	}
	return
}

// CheckAndCreateBucket 检查并创建桶
func (ds *AliyunOssStorage) CheckAndCreateBucket(bucketName string) (err error) {
	var exists bool
	// 检查存储桶是否已经存在。
	exists, err = ds.ossClient.IsBucketExist(bucketName)
	if err != nil {
		return
	}
	if exists {
		return
	}
	// 创建桶
	err = ds.ossClient.CreateBucket(bucketName)
	return
}

// Upload 上传
func (ds *AliyunOssStorage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	if ds.CheckAndCreateBucketEnabled {
		err = ds.CheckAndCreateBucket(bucketName)
		if err != nil {
			return
		}
	}
	if rename, ok := storage.FromRenameContext(ctx); ok {
		filename = rename(filename)
	}
	fullName = filename

	options := []oss.Option{}
	contentType, contentTypeExist := FromContentTypeContext(ctx)
	if contentTypeExist {
		options = append(options, oss.ContentType(contentType))
	} else {
		var detectContentType string
		detectContentType, err = mime.DetectContentType(filename)
		if err != nil {
			return
		}
		options = append(options, oss.ContentType(detectContentType))
	}
	md, mdExist := storage.FromIncomingContext(ctx)
	if mdExist {
		var mdBytes []byte
		mdBytes, err = json.Marshal(md)
		if err != nil {
			return
		}
		options = append(options, oss.Meta("Data", string(mdBytes)))
	}
	var bucket *oss.Bucket
	bucket, err = ds.ossClient.Bucket(bucketName)
	if err != nil {
		return
	}
	err = bucket.PutObject(filename, read, options...)
	return
}

// Download 下载
func (ds *AliyunOssStorage) Download(ctx context.Context, dist io.Writer, filename string) (info storage.DownloadFileInfoer, err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	var bucket *oss.Bucket
	bucket, err = ds.ossClient.Bucket(bucketName)
	if err != nil {
		return
	}
	var object io.ReadCloser
	object, err = bucket.GetObject(filename)
	if err != nil {
		return
	}
	defer object.Close()
	var meta http.Header
	meta, err = bucket.GetObjectMeta(filename)
	if err != nil {
		return
	}
	md := storage.Metadata{}
	for k, v := range meta {
		md.Set(k, v[0])
	}
	var (
		downloadFilename      string
		downloadFilenameExist bool
	)
	if downloadFilename, downloadFilenameExist = storage.FromDownloadFilenameContext(ctx); !downloadFilenameExist {
		downloadFilename = filepath.Base(filename)
	}
	var size int64
	if length := md.Get(oss.HTTPHeaderContentLength); length != "" {
		size = convert.ToInt64(length)
	}
	info = &downloadFileInfo{
		filename: downloadFilename,
		size:     size,
		metadata: md,
	}
	if downloadBefore, downloadBeforeExist := storage.FromDownloadBeforeContext(ctx); downloadBeforeExist {
		downloadBefore(info)
	}
	_, err = io.Copy(dist, object)
	if err != nil {
		info = nil
		return
	}
	return
}

// Remove 删除
func (ds *AliyunOssStorage) Remove(ctx context.Context, filename string) (err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	var bucket *oss.Bucket
	bucket, err = ds.ossClient.Bucket(bucketName)
	if err != nil {
		return
	}
	err = bucket.DeleteObject(bucketName)
	return
}
