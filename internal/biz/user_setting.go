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
	// UpdateUserSystemNotifyRecevieSetting 更新用户系统通知
	UpdateUserSystemNotifyRecevieSetting(context.Context, *v1.UpdateUserSystemNotifyRecevieSettingRequest) (*v1.UpdateUserSystemNotifyRecevieSettingReply, error)
	// UpdateUserInteractNotifyRecevieSetting 更新用户互动通知
	UpdateUserInteractNotifyRecevieSetting(context.Context, *v1.UpdateUserInteractNotifyRecevieSettingRequest) (*v1.UpdateUserInteractNotifyRecevieSettingReply, error)
	// UpdateUserAdoptNotifyRecevieSetting 更新用户领养通知
	UpdateUserAdoptNotifyRecevieSetting(context.Context, *v1.UpdateUserAdoptNotifyRecevieSettingRequest) (*v1.UpdateUserAdoptNotifyRecevieSettingReply, error)
	// UpdateUserEmailNotifyRecevieSetting 更新用户邮件通知
	UpdateUserEmailNotifyRecevieSetting(context.Context, *v1.UpdateUserEmailNotifyRecevieSettingRequest) (*v1.UpdateUserEmailNotifyRecevieSettingReply, error)
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

// UpdateUserSystemNotifyRecevieSetting 更新用户系统通知
func (u *UserSettingUsecase) UpdateUserSystemNotifyRecevieSetting(ctx context.Context, req *v1.UpdateUserSystemNotifyRecevieSettingRequest) (*v1.UpdateUserSystemNotifyRecevieSettingReply, error) {
	return u.repo.UpdateUserSystemNotifyRecevieSetting(ctx, req)
}

// UpdateUserInteractNotifyRecevieSetting 更新用户互动通知
func (u *UserSettingUsecase) UpdateUserInteractNotifyRecevieSetting(ctx context.Context, req *v1.UpdateUserInteractNotifyRecevieSettingRequest) (*v1.UpdateUserInteractNotifyRecevieSettingReply, error) {
	return u.repo.UpdateUserInteractNotifyRecevieSetting(ctx, req)
}

// UpdateUserAdoptNotifyRecevieSetting 更新用户领养通知
func (u *UserSettingUsecase) UpdateUserAdoptNotifyRecevieSetting(ctx context.Context, req *v1.UpdateUserAdoptNotifyRecevieSettingRequest) (*v1.UpdateUserAdoptNotifyRecevieSettingReply, error) {
	return u.repo.UpdateUserAdoptNotifyRecevieSetting(ctx, req)
}

// UpdateUserEmailNotifyRecevieSetting 更新用户邮件通知
func (u *UserSettingUsecase) UpdateUserEmailNotifyRecevieSetting(ctx context.Context, req *v1.UpdateUserEmailNotifyRecevieSettingRequest) (*v1.UpdateUserEmailNotifyRecevieSettingReply, error) {
	return u.repo.UpdateUserEmailNotifyRecevieSetting(ctx, req)
}

// GetUserNotifySetting 获取用户通知设置
func (u *UserSettingUsecase) GetUserNotifySetting(ctx context.Context, req *v1.GetUserNotifySettingRequest) (*v1.GetUserNotifySettingReply, error) {
	return u.repo.GetUserNotifySetting(ctx, req)
}
