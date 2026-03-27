package biz

import (
	"context"

	v1 "xiaomiao-home-system/api/user/v1"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// UserLoginRequest 用户登录
type UserLoginRequest struct {
	Username string
	Password string
}

// UserRepo is a Greater repo.
type UserRepo interface {
	// WebLogin 登录接口
	WebLogin(context.Context, *v1.WebLoginRequest) (*v1.WebLoginReply, error)
	// AppLogin 登录接口
	AppLogin(context.Context, *v1.AppLoginRequest) (*v1.AppLoginReply, error)
	// MpLogin 登录接口
	MpLogin(context.Context, *v1.MpLoginRequest) (*v1.MpLoginReply, error)
	// GetWebLoginUserInfo 查询登陆用户信息
	GetWebLoginUserInfo(context.Context, *v1.GetWebLoginUserInfoRequest) (*v1.GetWebLoginUserInfoReply, error)
	// WebLogout 退出登录
	WebLogout(context.Context, *v1.WebLogoutRequest) (*v1.WebLogoutReply, error)
	// WebCheckLogin web端登录检测
	WebCheckLogin(context.Context, *v1.WebCheckLoginRequest) (*v1.WebCheckLoginReply, error)
	// UpdateUserBaseSetting 更新用户基础设置
	UpdateUserBaseSetting(context.Context, *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserUsecase"))}
}

// WebLogin 用户登录
func (u *UserUsecase) WebLogin(ctx context.Context, req *v1.WebLoginRequest) (*v1.WebLoginReply, error) {
	return u.repo.WebLogin(ctx, req)
}

// AppLogin 用户登录
func (u *UserUsecase) AppLogin(ctx context.Context, req *v1.AppLoginRequest) (*v1.AppLoginReply, error) {
	return u.repo.AppLogin(ctx, req)
}

// MpLogin 用户登录
func (u *UserUsecase) MpLogin(ctx context.Context, req *v1.MpLoginRequest) (*v1.MpLoginReply, error) {
	return u.repo.MpLogin(ctx, req)
}

// GetWebLoginUserInfo 查询登陆用户信息
func (u *UserUsecase) GetWebLoginUserInfo(ctx context.Context, req *v1.GetWebLoginUserInfoRequest) (*v1.GetWebLoginUserInfoReply, error) {
	return u.repo.GetWebLoginUserInfo(ctx, req)
}

// WebLogout 退出登录
func (u *UserUsecase) WebLogout(ctx context.Context, req *v1.WebLogoutRequest) (*v1.WebLogoutReply, error) {
	return u.repo.WebLogout(ctx, req)
}

// WebCheckLogin web端登录检测
func (u *UserUsecase) WebCheckLogin(ctx context.Context, req *v1.WebCheckLoginRequest) (*v1.WebCheckLoginReply, error) {
	return u.repo.WebCheckLogin(ctx, req)
}

// UpdateUserBaseSetting 更新用户基础设置
func (u *UserUsecase) UpdateUserBaseSetting(ctx context.Context, req *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error) {
	return u.repo.UpdateUserBaseSetting(ctx, req)
}
