package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/cat/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// CatRepo is a Greater repo.
type CatRepo interface {
	// GetCatList 查询小猫列表
	GetCatList(ctx context.Context, req *v1.GetCatListRequest) (*v1.GetCatListReply, error)
}

// CatUsecase is a Cat usecase.
type CatUsecase struct {
	repo CatRepo
	log  *log.Helper
}

// NewCatUsecase new a Cat usecase.
func NewCatUsecase(repo CatRepo, logger log.Logger) *CatUsecase {
	return &CatUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "CatUsecase"))}
}

// GetCatList 查询小猫列表
func (u *CatUsecase) GetCatList(ctx context.Context, req *v1.GetCatListRequest) (*v1.GetCatListReply, error) {
	res, err := u.repo.GetCatList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
