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
	// GetSecretKey 查询加密串
	GetSecretKey(context.Context, *v1.GetSecretKeyRequest) (*v1.GetSecretKeyReply, error)
	// GetLoginUserInfo 查询用户信息
	GetLoginUserInfo(context.Context, *v1.GetLoginUserInfoRequest) (*v1.GetLoginUserInfoReply, error)
	// GetAuthMenuList 查询用户菜单信息
	GetAuthMenuList(context.Context, *v1.GetAuthMenuListRequest) (*v1.GetAuthMenuListReply, error)
	// Logout 登出接口
	Logout(context.Context, *v1.LogoutRequest) (*v1.LogoutReply, error)
	// GetUsers 查询所有用户
	GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersReply, error)
	// GetUserList 查询用户列表
	GetUserList(ctx context.Context, req *v1.GetUserListRequest) (*v1.GetUserListReply, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error)
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error)
	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, req *v1.UpdateUserStatusRequest) (*v1.UpdateUserStatusReply, error)
	// GetUserId 查询当前用户ID
	GetUserId(ctx context.Context) (int64, error)
	// GetUserNotifications 查询用户消息
	GetUserNotifications(ctx context.Context, req *v1.GetUserNotificationsRequest) (*v1.GetUserNotificationsReply, error)
	// UpdateUserNotificationStatus 更新用户消息状态
	UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error)
	// UpdateUserAllNotificationStatus 更新用户所有消息状态
	UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error)
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

// GetSecretKey 查询加密串
func (u *UserUsecase) GetSecretKey(ctx context.Context, req *v1.GetSecretKeyRequest) (*v1.GetSecretKeyReply, error) {

	res, err := u.repo.GetSecretKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetLoginUserInfo 查询用户信息
func (u *UserUsecase) GetLoginUserInfo(ctx context.Context, req *v1.GetLoginUserInfoRequest) (*v1.GetLoginUserInfoReply, error) {
	res, err := u.repo.GetLoginUserInfo(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAuthMenuList 查询用户菜单信息
func (u *UserUsecase) GetAuthMenuList(ctx context.Context, req *v1.GetAuthMenuListRequest) (*v1.GetAuthMenuListReply, error) {
	res, err := u.repo.GetAuthMenuList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// Logout 登出接口
func (u *UserUsecase) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutReply, error) {
	res, err := u.repo.Logout(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUsers 查询所有用户
func (u *UserUsecase) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersReply, error) {
	res, err := u.repo.GetUsers(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserList 查询用户列表
func (u *UserUsecase) GetUserList(ctx context.Context, req *v1.GetUserListRequest) (*v1.GetUserListReply, error) {
	res, err := u.repo.GetUserList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUser 创建用户
func (u *UserUsecase) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	res, err := u.repo.CreateUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUser 更新用户
func (u *UserUsecase) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	res, err := u.repo.UpdateUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUser 删除用户
func (u *UserUsecase) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	res, err := u.repo.DeleteUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserStatus 更新用户状态
func (u *UserUsecase) UpdateUserStatus(ctx context.Context, req *v1.UpdateUserStatusRequest) (*v1.UpdateUserStatusReply, error) {
	res, err := u.repo.UpdateUserStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserId 查询当前用户ID
func (u *UserUsecase) GetUserId(ctx context.Context) (int64, error) {
	res, err := u.repo.GetUserId(ctx)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GetUserNotifications 查询用户消息
func (u *UserUsecase) GetUserNotifications(ctx context.Context, req *v1.GetUserNotificationsRequest) (*v1.GetUserNotificationsReply, error) {
	res, err := u.repo.GetUserNotifications(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (u *UserUsecase) UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error) {
	res, err := u.repo.UpdateUserNotificationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserAllNotificationStatus 更新用户所有消息状态
func (u *UserUsecase) UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error) {
	res, err := u.repo.UpdateUserAllNotificationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
