package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/discover/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// DiscoverRepo is a Greater repo.
type DiscoverRepo interface {
	// GetDiscoverList 查询发现列表
	GetDiscoverList(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.GetDiscoverListReply, error)
}

// DiscoverUsecase is a Discover usecase.
type DiscoverUsecase struct {
	repo DiscoverRepo
	log  *log.Helper
}

// NewDiscoverUsecase new a Discover usecase.
func NewDiscoverUsecase(repo DiscoverRepo, logger log.Logger) *DiscoverUsecase {
	return &DiscoverUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "DiscoverUsecase"))}
}

// GetDiscoverList 查询发现列表
func (u *DiscoverUsecase) GetDiscoverList(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.GetDiscoverListReply, error) {
	res, err := u.repo.GetDiscoverList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
