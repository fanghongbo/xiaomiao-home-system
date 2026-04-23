package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/user/like/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserLikeRepo is a Greater repo.
type UserLikeRepo interface {
	// AddUserLike 添加用户收藏
	AddUserLike(ctx context.Context, req *v1.AddUserLikeRequest) (*v1.AddUserLikeReply, error)
	// CancelUserLike 取消用户收藏
	CancelUserLike(ctx context.Context, req *v1.CancelUserLikeRequest) (*v1.CancelUserLikeReply, error)
	// GetUserPostLikeStatus 查询用户发布内容收藏状态
	GetUserPostLikeStatus(ctx context.Context, postId int64) (bool, error)
	// GetUserLikeStatus 查询用户喜欢状态
	GetUserLikeStatus(ctx context.Context, req *v1.GetUserLikeStatusRequest) (*v1.GetUserLikeStatusReply, error)
}

// UserLikeUsecase is a UserLike usecase.
type UserLikeUsecase struct {
	repo UserLikeRepo
	log  *log.Helper
}

// NewUserLikeUsecase new a UserLike usecase.
func NewUserLikeUsecase(repo UserLikeRepo, logger log.Logger) *UserLikeUsecase {
	return &UserLikeUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserLikeUsecase"))}
}

// AddUserLike 添加用户收藏
func (u *UserLikeUsecase) AddUserLike(ctx context.Context, req *v1.AddUserLikeRequest) (*v1.AddUserLikeReply, error) {
	res, err := u.repo.AddUserLike(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CancelUserLike 取消用户收藏
func (u *UserLikeUsecase) CancelUserLike(ctx context.Context, req *v1.CancelUserLikeRequest) (*v1.CancelUserLikeReply, error) {
	res, err := u.repo.CancelUserLike(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserPostLikeStatus 查询用户发布内容收藏状态
func (u *UserLikeUsecase) GetUserPostLikeStatus(ctx context.Context, postId int64) (bool, error) {
	res, err := u.repo.GetUserPostLikeStatus(ctx, postId)

	if err != nil {
		return false, err
	}

	return res, nil
}

// GetUserLikeStatus 查询用户喜欢状态
func (u *UserLikeUsecase) GetUserLikeStatus(ctx context.Context, req *v1.GetUserLikeStatusRequest) (*v1.GetUserLikeStatusReply, error) {
	res, err := u.repo.GetUserLikeStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
