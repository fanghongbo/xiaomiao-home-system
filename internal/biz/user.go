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
	// Login 登录接口
	Login(context.Context, *v1.LoginRequest) (*v1.LoginReply, error)
	// GetUserId 查询用户ID
	GetUserId(context.Context) (int64, error)
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

// Login 用户登录
func (u *UserUsecase) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	u.log.WithContext(ctx).Infof("user login: %v", req.Username)

	res, err := u.repo.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserId 查询用户ID
func (u *UserUsecase) GetUserId(ctx context.Context) (int64, error) {
	return u.repo.GetUserId(ctx)
}
