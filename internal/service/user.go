package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServer

	user   *biz.UserUsecase
	log    *log.Helper
	config *conf.Config
}

func NewUserService(user *biz.UserUsecase, config *conf.Config, logger log.Logger) *UserService {
	return &UserService{
		user:   user,
		config: config,
		log:    log.NewHelper(log.With(logger, "service", "UserService")),
	}
}

// Login 登录接口
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	res, err := s.user.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetSecretKey 查询加密串
func (s *UserService) GetSecretKey(ctx context.Context, req *pb.GetSecretKeyRequest) (*pb.GetSecretKeyReply, error) {
	res, err := s.user.GetSecretKey(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetLoginUserInfo 查询用户信息
func (s *UserService) GetLoginUserInfo(ctx context.Context, req *pb.GetLoginUserInfoRequest) (*pb.GetLoginUserInfoReply, error) {
	res, err := s.user.GetLoginUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAuthMenuList 查询用户菜单信息
func (s *UserService) GetAuthMenuList(ctx context.Context, req *pb.GetAuthMenuListRequest) (*pb.GetAuthMenuListReply, error) {
	res, err := s.user.GetAuthMenuList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Logout 登出接口
func (s *UserService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	res, err := s.user.Logout(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUsers 查询所有用户
func (s *UserService) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	res, err := s.user.GetUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserList 查询用户列表
func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	res, err := s.user.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	res, err := s.user.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	res, err := s.user.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	res, err := s.user.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*pb.UpdateUserStatusReply, error) {
	res, err := s.user.UpdateUserStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserNotifications 查询用户消息
func (s *UserService) GetUserNotifications(ctx context.Context, req *pb.GetUserNotificationsRequest) (*pb.GetUserNotificationsReply, error) {
	res, err := s.user.GetUserNotifications(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (s *UserService) UpdateUserNotificationStatus(ctx context.Context, req *pb.UpdateUserNotificationStatusRequest) (*pb.UpdateUserNotificationStatusReply, error) {
	res, err := s.user.UpdateUserNotificationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserAllNotificationStatus 更新用户所有消息状态
func (s *UserService) UpdateUserAllNotificationStatus(ctx context.Context, req *pb.UpdateUserAllNotificationStatusRequest) (*pb.UpdateUserAllNotificationStatusReply, error) {
	res, err := s.user.UpdateUserAllNotificationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
