package data

import (
	"context"
	v1 "xiaomiao-home-system/api/file/v1"
	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/log"
)

type fileRepo struct {
	data *Data
	log  *log.Helper
}

// NewFileRepo .
func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "FileRepo")),
	}
}

// DownloadFile 下载文件
func (f *fileRepo) DownloadFile(ctx context.Context, req *v1.DownloadFileRequest) (*v1.DownloadFileReply, error) {
	return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
}
