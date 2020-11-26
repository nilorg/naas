package api

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/module/store"
	"github.com/nilorg/sdk/storage"
	"github.com/spf13/viper"
)

type file struct {
}

// Upload 文件上传
// @Tags 		文件
// @Summary		文件上传
// @Accept multipart/form-data
// @Produce	json
// @Param file formData file true "upload file"
// @Success 200	{object}	Result
// @Router /files [POST]
// @Security OAuth2AccessCode
func (*file) Upload(ctx *gin.Context) {
	var (
		finfos []storage.UploadFileInfoer
		err    error
	)
	ctxUpload := storage.NewRenameContext(context.Background(), store.FileRename)
	finfos, err = storage.UploadHandle(ctxUpload, ctx.Request, store.Picture, "file")
	if err != nil {
		writeError(ctx, err)
		return
	}
	var u *url.URL
	u, err = url.Parse(viper.GetString("storage.public_path"))
	if err != nil {
		writeError(ctx, err)
		return
	}
	var q string
	if q = ctx.Query("q"); q == "" {
		q = "picture"
	}
	var values = make([]gin.H, 0)
	for _, file := range finfos {
		u.Path = path.Join(u.Path, fmt.Sprintf("/%s/", q), filepath.Base(file.FullName()))
		values = append(values, gin.H{
			"fullName": u.String(),
			"filename": file.Filename(),
			"size":     file.Size(),
		})
	}
	writeData(ctx, values)
}
