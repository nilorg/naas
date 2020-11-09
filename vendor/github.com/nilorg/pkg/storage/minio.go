package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v6"
	"github.com/nilorg/sdk/mime"
	"github.com/nilorg/sdk/storage"
)

var (
	// ErrBucketNameNotIsNil 桶名称不能为空
	ErrBucketNameNotIsNil = errors.New("BucketName not is nil")
)

// MinioStorage minio存储
type MinioStorage struct {
	location                    string
	bucketNames                 []string
	minioClient                 *minio.Client
	CheckAndCreateBucketEnabled bool
}

// NewMinioStorage 创建minio存储
func NewMinioStorage(minioClient *minio.Client, location string, initBucket bool, bucketNames []string) (ms *MinioStorage, err error) {
	ms = &MinioStorage{
		location:                    location,
		bucketNames:                 bucketNames,
		minioClient:                 minioClient,
		CheckAndCreateBucketEnabled: true,
	}
	if initBucket {
		err = ms.initBucket()
		if err != nil {
			ms = nil
		}
	}
	return
}

// initBucket 初始化桶
func (ds *MinioStorage) initBucket() (err error) {
	for _, bucketName := range ds.bucketNames {
		err = ds.CheckAndCreateBucket(bucketName)
		if err != nil {
			return
		}
	}
	return
}

// CheckAndCreateBucket 检查并创建桶
func (ds *MinioStorage) CheckAndCreateBucket(bucketName string) (err error) {
	var exists bool
	// 检查存储桶是否已经存在。
	exists, err = ds.minioClient.BucketExists(bucketName)
	if err != nil {
		return
	}
	if exists {
		return
	}
	// 创建桶
	err = ds.minioClient.MakeBucket(bucketName, ds.location)
	return
}

// Upload 上传
func (ds *MinioStorage) Upload(ctx context.Context, read io.Reader, filename string) (fullName string, err error) {
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
	options := minio.PutObjectOptions{}
	contextType, contextTypeExist := FromContentTypeContext(ctx)
	if contextTypeExist {
		options.ContentType = contextType
	} else {
		suffix := filepath.Ext(filename)
		suffixContextType, suffixContextTypeExist := mime.Lookup(suffix)
		if !suffixContextTypeExist {
			err = fmt.Errorf("%s unrecognized suffix", suffix)
			return
		}
		options.ContentType = suffixContextType
	}
	md, mdExist := storage.FromIncomingContext(ctx)
	if mdExist {
		options.UserMetadata = md
	}
	_, err = ds.minioClient.PutObjectWithContext(ctx, bucketName, filename, read, -1, options)
	return
}

// Download 下载
func (ds *MinioStorage) Download(ctx context.Context, dist io.Writer, filename string) (info storage.DownloadFileInfoer, err error) {
	var (
		object *minio.Object
	)
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	object, err = ds.minioClient.GetObjectWithContext(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	defer object.Close()

	var objectInfo minio.ObjectInfo
	objectInfo, err = object.Stat()
	if err != nil {
		return
	}
	md := storage.Metadata{
		"Content-Type": objectInfo.ContentType,
	}
	for k, v := range objectInfo.UserMetadata {
		md.Set(k, v)
	}

	var (
		downloadFilename      string
		downloadFilenameExist bool
	)
	if downloadFilename, downloadFilenameExist = storage.FromDownloadFilenameContext(ctx); !downloadFilenameExist {
		downloadFilename = filepath.Base(filename)
	}

	info = &downloadFileInfo{
		filename: downloadFilename,
		size:     objectInfo.Size,
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
func (ds *MinioStorage) Remove(ctx context.Context, filename string) (err error) {
	bucketName, bucketNameOk := FromBucketNameContext(ctx)
	if !bucketNameOk {
		err = ErrBucketNameNotIsNil
		return
	}
	err = ds.minioClient.RemoveObject(bucketName, filename)
	return
}

type downloadFileInfo struct {
	size     int64
	filename string
	metadata storage.Metadata
}

func (dfi *downloadFileInfo) Size() int64 {
	return dfi.size
}

func (dfi *downloadFileInfo) Filename() string {
	return dfi.filename
}

func (dfi *downloadFileInfo) Metadata() storage.Metadata {
	return dfi.metadata
}
