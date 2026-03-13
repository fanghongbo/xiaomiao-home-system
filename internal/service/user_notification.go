package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/notification/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserNotificationService struct {
	pb.UnimplementedUserNotificationServer

	userNotification *biz.UserNotificationUsecase
	log              *log.Helper
	config           *conf.Config
}

func NewUserNotificationService(userNotification *biz.UserNotificationUsecase, config *conf.Config, logger log.Logger) *UserNotificationService {
	return &UserNotificationService{
		userNotification: userNotification,
		config:           config,
		log:              log.NewHelper(log.With(logger, "service", "UserNotificationService")),
	}
}

// GetUserNotificationList 查询用户消息列表
func (s *UserNotificationService) GetUserNotificationList(ctx context.Context, req *pb.GetUserNotificationListRequest) (*pb.GetUserNotificationListReply, error) {
	res, err := s.userNotification.GetUserNotificationList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserNotification 删除用户消息
func (s *UserNotificationService) DeleteUserNotification(ctx context.Context, req *pb.DeleteUserNotificationRequest) (*pb.DeleteUserNotificationReply, error) {
	res, err := s.userNotification.DeleteUserNotification(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (s *UserNotificationService) UpdateUserNotificationStatus(ctx context.Context, req *pb.UpdateUserNotificationStatusRequest) (*pb.UpdateUserNotificationStatusReply, error) {
	res, err := s.userNotification.UpdateUserNotificationStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// BatchUserNotifications 批量操作用户消息
func (s *UserNotificationService) BatchUserNotifications(ctx context.Context, req *pb.BatchUserNotificationsRequest) (*pb.BatchUserNotificationsReply, error) {
	res, err := s.userNotification.BatchUserNotifications(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
