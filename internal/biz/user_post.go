package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/user/post/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserPostRepo is a Greater repo.
type UserPostRepo interface {
	// GetUserPostList 查询发布内容列表
	GetUserPostList(ctx context.Context, req *v1.GetUserPostListRequest) (*v1.GetUserPostListReply, error)
	// CreateUserPost 创建发布内容
	CreateUserPost(ctx context.Context, req *v1.CreateUserPostRequest) (*v1.CreateUserPostReply, error)
	// UpdateUserPost 更新发布内容
	UpdateUserPost(ctx context.Context, req *v1.UpdateUserPostRequest) (*v1.UpdateUserPostReply, error)
	// DeleteUserPost 删除发布内容
	DeleteUserPost(ctx context.Context, req *v1.DeleteUserPostRequest) (*v1.DeleteUserPostReply, error)
	// UpdateUserPostStatus 更新发布内容状态
	UpdateUserPostStatus(ctx context.Context, req *v1.UpdateUserPostStatusRequest) (*v1.UpdateUserPostStatusReply, error)
	// GetUserPost 查询发布内容
	GetUserPost(ctx context.Context, req *v1.GetUserPostRequest) (*v1.GetUserPostReply, error)
	// GetPostCatInfo 查询发布内容分类信息
	GetPostCatInfo(ctx context.Context, postId int64) (*v1.CatInfo, error)
	// GetPostUserInfo 查询发布内容用户信息
	GetPostUserInfo(ctx context.Context, postId int64) (*v1.UserInfo, error)
}

// UserPostUsecase is a UserPost usecase.
type UserPostUsecase struct {
	repo UserPostRepo
	log  *log.Helper
}

// NewUserPostUsecase new a UserPost usecase.
func NewUserPostUsecase(repo UserPostRepo, logger log.Logger) *UserPostUsecase {
	return &UserPostUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserPostUsecase"))}
}

// GetUserPostList 查询发布内容列表
func (u *UserPostUsecase) GetUserPostList(ctx context.Context, req *v1.GetUserPostListRequest) (*v1.GetUserPostListReply, error) {
	res, err := u.repo.GetUserPostList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserPost 创建发布内容
func (u *UserPostUsecase) CreateUserPost(ctx context.Context, req *v1.CreateUserPostRequest) (*v1.CreateUserPostReply, error) {
	res, err := u.repo.CreateUserPost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserPost 更新发布内容
func (u *UserPostUsecase) UpdateUserPost(ctx context.Context, req *v1.UpdateUserPostRequest) (*v1.UpdateUserPostReply, error) {
	res, err := u.repo.UpdateUserPost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserPost 删除发布内容
func (u *UserPostUsecase) DeleteUserPost(ctx context.Context, req *v1.DeleteUserPostRequest) (*v1.DeleteUserPostReply, error) {
	res, err := u.repo.DeleteUserPost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserPostStatus 更新发布内容状态
func (u *UserPostUsecase) UpdateUserPostStatus(ctx context.Context, req *v1.UpdateUserPostStatusRequest) (*v1.UpdateUserPostStatusReply, error) {
	res, err := u.repo.UpdateUserPostStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserPost 查询发布内容
func (u *UserPostUsecase) GetUserPost(ctx context.Context, req *v1.GetUserPostRequest) (*v1.GetUserPostReply, error) {
	res, err := u.repo.GetUserPost(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetPostCatInfo 查询发布内容分类信息
func (u *UserPostUsecase) GetPostCatInfo(ctx context.Context, postId int64) (*v1.CatInfo, error) {
	res, err := u.repo.GetPostCatInfo(ctx, postId)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetPostUserInfo 查询发布内容用户信息
func (u *UserPostUsecase) GetPostUserInfo(ctx context.Context, postId int64) (*v1.UserInfo, error) {
	res, err := u.repo.GetPostUserInfo(ctx, postId)

	if err != nil {
		return nil, err
	}

	return res, nil
}
