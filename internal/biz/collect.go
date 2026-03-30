package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/collect/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// CollectRepo is a Greater repo.
type CollectRepo interface {
	// GetCollectList 查询收藏列表
	GetCollectList(ctx context.Context, req *v1.GetCollectListRequest) (*v1.GetCollectListReply, error)
	// GetCollectTypes 查询收藏分类
	GetCollectTypes(ctx context.Context, req *v1.GetCollectTypesRequest) (*v1.GetCollectTypesReply, error)
}

// CollectUsecase is a Collect usecase.
type CollectUsecase struct {
	repo CollectRepo
	log  *log.Helper
}

// NewCollectUsecase new a Collect usecase.
func NewCollectUsecase(repo CollectRepo, logger log.Logger) *CollectUsecase {
	return &CollectUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "CollectUsecase"))}
}

// GetCollectList 查询收藏列表
func (u *CollectUsecase) GetCollectList(ctx context.Context, req *v1.GetCollectListRequest) (*v1.GetCollectListReply, error) {
	res, err := u.repo.GetCollectList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetCollectTypes 查询收藏分类
func (u *CollectUsecase) GetCollectTypes(ctx context.Context, req *v1.GetCollectTypesRequest) (*v1.GetCollectTypesReply, error) {
	res, err := u.repo.GetCollectTypes(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
