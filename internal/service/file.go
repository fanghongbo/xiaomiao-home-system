package service

import (
	"context"
	pb "xiaomiao-home-system/api/file/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type FileService struct {
	pb.UnimplementedFileServer

	file   *biz.FileUsecase
	log    *log.Helper
	config *conf.Config
}

func NewFileService(file *biz.FileUsecase, config *conf.Config, logger log.Logger) *FileService {
	return &FileService{
		file:   file,
		config: config,
		log:    log.NewHelper(log.With(logger, "service", "FileService")),
	}
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileReply, error) {
	return s.file.DownloadFile(ctx, req)
}

// UpdateAvatar 上传头像
func (s *FileService) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarReply, error) {
	return s.file.UploadAvatar(ctx, req)
}
