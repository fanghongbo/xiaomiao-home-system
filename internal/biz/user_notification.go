package biz

import (
	"context"

	v1 "xiaomiao-home-system/api/user/notification/v1"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserNotificationNotFound is user group not found.
// ErrUserNotificationNotFound = errors.NotFound(v1.ErrorReason_USER_GROUP_NOT_FOUND.String(), "user group not found")
)

// UserNotificationRepo is a Greater repo.
type UserNotificationRepo interface {
	// GetUserNotificationList 查询用户消息列表
	GetUserNotificationList(ctx context.Context, req *v1.GetUserNotificationListRequest) (*v1.GetUserNotificationListReply, error)
	// DeleteUserNotification 删除用户消息
	DeleteUserNotification(ctx context.Context, req *v1.DeleteUserNotificationRequest) (*v1.DeleteUserNotificationReply, error)
	// UpdateUserNotificationStatus 更新用户消息状态
	UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error)
	// CreateUserNotification 创建用户消息
	CreateUserNotification(ctx context.Context, req *v1.CreateUserNotificationRequest) error
	// BatchUserNotifications 批量操作用户消息
	BatchUserNotifications(ctx context.Context, req *v1.BatchUserNotificationsRequest) (*v1.BatchUserNotificationsReply, error)
	// UpdateUserAllNotificationStatus 更新用户消息状态
	UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error)
}

// UserNotificationUsecase is a User usecase.
type UserNotificationUsecase struct {
	repo UserNotificationRepo
	log  *log.Helper
}

// NewUserNotificationUsecase new a UserNotification usecase.
func NewUserNotificationUsecase(repo UserNotificationRepo, logger log.Logger) *UserNotificationUsecase {
	return &UserNotificationUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserNotificationUsecase"))}
}

// GetUserNotificationList 查询用户消息列表
func (u *UserNotificationUsecase) GetUserNotificationList(ctx context.Context, req *v1.GetUserNotificationListRequest) (*v1.GetUserNotificationListReply, error) {
	res, err := u.repo.GetUserNotificationList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserNotification 删除用户消息
func (u *UserNotificationUsecase) DeleteUserNotification(ctx context.Context, req *v1.DeleteUserNotificationRequest) (*v1.DeleteUserNotificationReply, error) {
	res, err := u.repo.DeleteUserNotification(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (u *UserNotificationUsecase) UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error) {
	res, err := u.repo.UpdateUserNotificationStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserNotification 创建用户消息
func (u *UserNotificationUsecase) CreateUserNotification(ctx context.Context, req *v1.CreateUserNotificationRequest) error {
	err := u.repo.CreateUserNotification(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

// BatchUserNotifications 批量操作用户消息
func (u *UserNotificationUsecase) BatchUserNotifications(ctx context.Context, req *v1.BatchUserNotificationsRequest) (*v1.BatchUserNotificationsReply, error) {
	res, err := u.repo.BatchUserNotifications(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserAllNotificationStatus 更新用户消息状态
func (u *UserNotificationUsecase) UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error) {
	res, err := u.repo.UpdateUserAllNotificationStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
