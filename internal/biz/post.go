package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/post/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// PostRepo is a Greater repo.
type PostRepo interface {
	// GetPostList 查询发布内容列表
	GetPostList(ctx context.Context, req *v1.GetPostListRequest) (*v1.GetPostListReply, error)
	// CreatePost 创建发布内容
	CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostReply, error)
	// UpdatePost 更新发布内容
	UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostReply, error)
	// DeletePost 删除发布内容
	DeletePost(ctx context.Context, req *v1.DeletePostRequest) (*v1.DeletePostReply, error)
	// UpdatePostStatus 更新发布内容状态
	UpdatePostStatus(ctx context.Context, req *v1.UpdatePostStatusRequest) (*v1.UpdatePostStatusReply, error)
	// GetPost 查询发布内容
	GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.GetPostReply, error)
}

// PostUsecase is a Post usecase.
type PostUsecase struct {
	repo PostRepo
	log  *log.Helper
}

// NewPostUsecase new a Post usecase.
func NewPostUsecase(repo PostRepo, logger log.Logger) *PostUsecase {
	return &PostUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "PostUsecase"))}
}

// GetPostList 查询发布内容列表
func (u *PostUsecase) GetPostList(ctx context.Context, req *v1.GetPostListRequest) (*v1.GetPostListReply, error) {
	res, err := u.repo.GetPostList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePost 创建发布内容
func (u *PostUsecase) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostReply, error) {
	res, err := u.repo.CreatePost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePost 更新发布内容
func (u *PostUsecase) UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostReply, error) {
	res, err := u.repo.UpdatePost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeletePost 删除发布内容
func (u *PostUsecase) DeletePost(ctx context.Context, req *v1.DeletePostRequest) (*v1.DeletePostReply, error) {
	res, err := u.repo.DeletePost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePostStatus 更新发布内容状态
func (u *PostUsecase) UpdatePostStatus(ctx context.Context, req *v1.UpdatePostStatusRequest) (*v1.UpdatePostStatusReply, error) {
	res, err := u.repo.UpdatePostStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetPost 查询发布内容
func (u *PostUsecase) GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.GetPostReply, error) {
	res, err := u.repo.GetPost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
