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
