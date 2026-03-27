package biz

import (
	"context"

	pb "xiaomiao-home-system/api/file/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// FileRepo is a Greater repo.
type FileRepo interface {
	// DownloadFile 下载文件
	DownloadFile(context.Context, *pb.DownloadFileRequest) (*pb.DownloadFileReply, error)
	// UploadAvatar 上传头像
	UploadAvatar(context.Context, *pb.UploadAvatarRequest) (*pb.UploadAvatarReply, error)
}

// FileUsecase is a File usecase.
type FileUsecase struct {
	repo FileRepo
	log  *log.Helper
}

// NewFileUsecase new a File usecase.
func NewFileUsecase(repo FileRepo, logger log.Logger) *FileUsecase {
	return &FileUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "FileUsecase"))}
}

// DownloadFile 下载文件
func (u *FileUsecase) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileReply, error) {
	return u.repo.DownloadFile(ctx, req)
}

// UpdateAvatar 上传头像
func (u *FileUsecase) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarReply, error) {
	return u.repo.UploadAvatar(ctx, req)
}
