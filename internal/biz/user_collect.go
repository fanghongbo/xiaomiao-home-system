package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/user/collect/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserCollectRepo is a Greater repo.
type UserCollectRepo interface {
	// GetUserCollectList 查询用户收藏列表
	GetUserCollectList(ctx context.Context, req *v1.GetUserCollectListRequest) (*v1.GetUserCollectListReply, error)
	// GetUserCollectTypes 查询用户收藏分类
	GetUserCollectTypes(ctx context.Context, req *v1.GetUserCollectTypesRequest) (*v1.GetUserCollectTypesReply, error)
	// AddUserCollect 添加用户收藏
	AddUserCollect(ctx context.Context, req *v1.AddUserCollectRequest) (*v1.AddUserCollectReply, error)
	// CancelUserCollect 取消用户收藏
	CancelUserCollect(ctx context.Context, req *v1.CancelUserCollectRequest) (*v1.CancelUserCollectReply, error)
}

// UserCollectUsecase is a UserCollect usecase.
type UserCollectUsecase struct {
	repo UserCollectRepo
	log  *log.Helper
}

// NewUserCollectUsecase new a UserCollect usecase.
func NewUserCollectUsecase(repo UserCollectRepo, logger log.Logger) *UserCollectUsecase {
	return &UserCollectUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserCollectUsecase"))}
}

// GetUserCollectList 查询用户收藏列表
func (u *UserCollectUsecase) GetUserCollectList(ctx context.Context, req *v1.GetUserCollectListRequest) (*v1.GetUserCollectListReply, error) {
	res, err := u.repo.GetUserCollectList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserCollectTypes 查询用户收藏分类
func (u *UserCollectUsecase) GetUserCollectTypes(ctx context.Context, req *v1.GetUserCollectTypesRequest) (*v1.GetUserCollectTypesReply, error) {
	res, err := u.repo.GetUserCollectTypes(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// AddUserCollect 添加用户收藏
func (u *UserCollectUsecase) AddUserCollect(ctx context.Context, req *v1.AddUserCollectRequest) (*v1.AddUserCollectReply, error) {
	res, err := u.repo.AddUserCollect(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CancelUserCollect 取消用户收藏
func (u *UserCollectUsecase) CancelUserCollect(ctx context.Context, req *v1.CancelUserCollectRequest) (*v1.CancelUserCollectReply, error) {
	res, err := u.repo.CancelUserCollect(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
