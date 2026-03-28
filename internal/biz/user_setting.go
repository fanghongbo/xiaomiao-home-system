package biz

import (
	"context"

	v1 "xiaomiao-home-system/api/user/setting/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// UserSettingRepo is a UserSetting repo.
type UserSettingRepo interface {
	// UpdateUserBaseSetting 更新用户基础更新
	UpdateUserBaseSetting(context.Context, *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error)
	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(context.Context, *v1.UpdateUserPasswordRequest) (*v1.UpdateUserPasswordReply, error)
	// UpdateUserSystemNotifySetting 更新用户系统通知
	UpdateUserSystemNotifySetting(context.Context, *v1.UpdateUserSystemNotifySettingRequest) (*v1.UpdateUserSystemNotifySettingReply, error)
	// UpdateUserInteractNotifySetting 更新用户互动通知
	UpdateUserInteractNotifySetting(context.Context, *v1.UpdateUserInteractNotifySettingRequest) (*v1.UpdateUserInteractNotifySettingReply, error)
	// UpdateUserAdoptNotifySetting 更新用户领养通知
	UpdateUserAdoptNotifySetting(context.Context, *v1.UpdateUserAdoptNotifySettingRequest) (*v1.UpdateUserAdoptNotifySettingReply, error)
	// UpdateUserEmailNotifySetting 更新用户邮件通知
	UpdateUserEmailNotifySetting(context.Context, *v1.UpdateUserEmailNotifySettingRequest) (*v1.UpdateUserEmailNotifySettingReply, error)
	// GetUserNotifySetting 获取用户通知设置
	GetUserNotifySetting(context.Context, *v1.GetUserNotifySettingRequest) (*v1.GetUserNotifySettingReply, error)
}

// UserSettingUsecase is a UserSetting usecase.
type UserSettingUsecase struct {
	repo UserSettingRepo
	log  *log.Helper
}

// NewUserSettingUsecase new a UserSetting usecase.
func NewUserSettingUsecase(repo UserSettingRepo, logger log.Logger) *UserSettingUsecase {
	return &UserSettingUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserSettingUsecase"))}
}

// UpdateUserBaseSetting 更新用户基础更新
func (u *UserSettingUsecase) UpdateUserBaseSetting(ctx context.Context, req *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error) {
	return u.repo.UpdateUserBaseSetting(ctx, req)
}

// UpdateUserPassword 更新用户密码
func (u *UserSettingUsecase) UpdateUserPassword(ctx context.Context, req *v1.UpdateUserPasswordRequest) (*v1.UpdateUserPasswordReply, error) {
	return u.repo.UpdateUserPassword(ctx, req)
}

// UpdateUserSystemNotifySetting 更新用户系统通知
func (u *UserSettingUsecase) UpdateUserSystemNotifySetting(ctx context.Context, req *v1.UpdateUserSystemNotifySettingRequest) (*v1.UpdateUserSystemNotifySettingReply, error) {
	return u.repo.UpdateUserSystemNotifySetting(ctx, req)
}

// UpdateUserInteractNotifySetting 更新用户互动通知
func (u *UserSettingUsecase) UpdateUserInteractNotifySetting(ctx context.Context, req *v1.UpdateUserInteractNotifySettingRequest) (*v1.UpdateUserInteractNotifySettingReply, error) {
	return u.repo.UpdateUserInteractNotifySetting(ctx, req)
}

// UpdateUserAdoptNotifySetting 更新用户领养通知
func (u *UserSettingUsecase) UpdateUserAdoptNotifySetting(ctx context.Context, req *v1.UpdateUserAdoptNotifySettingRequest) (*v1.UpdateUserAdoptNotifySettingReply, error) {
	return u.repo.UpdateUserAdoptNotifySetting(ctx, req)
}

// UpdateUserEmailNotifySetting 更新用户邮件通知
func (u *UserSettingUsecase) UpdateUserEmailNotifySetting(ctx context.Context, req *v1.UpdateUserEmailNotifySettingRequest) (*v1.UpdateUserEmailNotifySettingReply, error) {
	return u.repo.UpdateUserEmailNotifySetting(ctx, req)
}

// GetUserNotifySetting 获取用户通知设置
func (u *UserSettingUsecase) GetUserNotifySetting(ctx context.Context, req *v1.GetUserNotifySettingRequest) (*v1.GetUserNotifySettingReply, error) {
	return u.repo.GetUserNotifySetting(ctx, req)
}
