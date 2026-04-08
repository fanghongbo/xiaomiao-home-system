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
	// CreateCat 创建小猫
	CreateCat(ctx context.Context, req *v1.CreateCatRequest) (*v1.CreateCatReply, error)
	// UpdateCat 更新小猫
	UpdateCat(ctx context.Context, req *v1.UpdateCatRequest) (*v1.UpdateCatReply, error)
	// DeleteCat 删除小猫
	DeleteCat(ctx context.Context, req *v1.DeleteCatRequest) (*v1.DeleteCatReply, error)
	// GetCat 查询小猫信息
	GetCat(ctx context.Context, req *v1.GetCatRequest) (*v1.GetCatReply, error)
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

// CreateCat 创建小猫
func (u *CatUsecase) CreateCat(ctx context.Context, req *v1.CreateCatRequest) (*v1.CreateCatReply, error) {
	res, err := u.repo.CreateCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCat 更新小猫
func (u *CatUsecase) UpdateCat(ctx context.Context, req *v1.UpdateCatRequest) (*v1.UpdateCatReply, error) {
	res, err := u.repo.UpdateCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteCat 删除小猫
func (u *CatUsecase) DeleteCat(ctx context.Context, req *v1.DeleteCatRequest) (*v1.DeleteCatReply, error) {
	res, err := u.repo.DeleteCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetCat 查询小猫信息
func (u *CatUsecase) GetCat(ctx context.Context, req *v1.GetCatRequest) (*v1.GetCatReply, error) {
	res, err := u.repo.GetCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
