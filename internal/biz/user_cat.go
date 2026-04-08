package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/user/cat/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserCatRepo is a Greater repo.
type UserCatRepo interface {
	// GetUserCatList 查询我的小猫列表
	GetUserCatList(ctx context.Context, req *v1.GetUserCatListRequest) (*v1.GetUserCatListReply, error)
	// CreateUserCat 创建我的小猫
	CreateUserCat(ctx context.Context, req *v1.CreateUserCatRequest) (*v1.CreateUserCatReply, error)
	// UpdateUserCat 更新我的小猫
	UpdateUserCat(ctx context.Context, req *v1.UpdateUserCatRequest) (*v1.UpdateUserCatReply, error)
	// DeleteUserCat 删除我的小猫
	DeleteUserCat(ctx context.Context, req *v1.DeleteUserCatRequest) (*v1.DeleteUserCatReply, error)
	// GetUserCat 查询我的小猫信息
	GetUserCat(ctx context.Context, req *v1.GetUserCatRequest) (*v1.GetUserCatReply, error)
}

// UserCatUsecase is a UserCat usecase.
type UserCatUsecase struct {
	repo UserCatRepo
	log  *log.Helper
}

// NewUserCatUsecase new a UserCat usecase.
func NewUserCatUsecase(repo UserCatRepo, logger log.Logger) *UserCatUsecase {
	return &UserCatUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserCatUsecase"))}
}

// GetUserCatList 查询我的小猫列表
func (u *UserCatUsecase) GetUserCatList(ctx context.Context, req *v1.GetUserCatListRequest) (*v1.GetUserCatListReply, error) {
	res, err := u.repo.GetUserCatList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserCat 创建我的小猫
func (u *UserCatUsecase) CreateUserCat(ctx context.Context, req *v1.CreateUserCatRequest) (*v1.CreateUserCatReply, error) {
	res, err := u.repo.CreateUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateUserCat 更新我的小猫
func (u *UserCatUsecase) UpdateUserCat(ctx context.Context, req *v1.UpdateUserCatRequest) (*v1.UpdateUserCatReply, error) {
	res, err := u.repo.UpdateUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteUserCat 删除我的小猫
func (u *UserCatUsecase) DeleteUserCat(ctx context.Context, req *v1.DeleteUserCatRequest) (*v1.DeleteUserCatReply, error) {
	res, err := u.repo.DeleteUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetUserCat 查询我的小猫信息
func (u *UserCatUsecase) GetUserCat(ctx context.Context, req *v1.GetUserCatRequest) (*v1.GetUserCatReply, error) {
	res, err := u.repo.GetUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
