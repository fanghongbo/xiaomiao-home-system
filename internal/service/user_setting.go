package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/setting/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserSettingService struct {
	pb.UnimplementedUserSettingServer

	userSetting *biz.UserSettingUsecase
	log         *log.Helper
	config      *conf.Config
}

func NewUserSettingService(userSetting *biz.UserSettingUsecase, config *conf.Config, logger log.Logger) *UserSettingService {
	return &UserSettingService{
		userSetting: userSetting,
		config:      config,
		log:         log.NewHelper(log.With(logger, "service", "UserSettingService")),
	}
}

// UpdateUserBaseSetting 更新用户基础更新
func (s *UserSettingService) UpdateUserBaseSetting(ctx context.Context, req *pb.UpdateUserBaseSettingRequest) (*pb.UpdateUserBaseSettingReply, error) {
	return s.userSetting.UpdateUserBaseSetting(ctx, req)
}

// UpdateUserPassword 更新用户密码
func (s *UserSettingService) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordReply, error) {
	return s.userSetting.UpdateUserPassword(ctx, req)
}

// UpdateUserSystemNotifyRecevieSetting 更新用户系统通知
func (s *UserSettingService) UpdateUserSystemNotifyRecevieSetting(ctx context.Context, req *pb.UpdateUserSystemNotifyRecevieSettingRequest) (*pb.UpdateUserSystemNotifyRecevieSettingReply, error) {
	return s.userSetting.UpdateUserSystemNotifyRecevieSetting(ctx, req)
}

// UpdateUserInteractNotifyRecevieSetting 更新用户互动通知
func (s *UserSettingService) UpdateUserInteractNotifyRecevieSetting(ctx context.Context, req *pb.UpdateUserInteractNotifyRecevieSettingRequest) (*pb.UpdateUserInteractNotifyRecevieSettingReply, error) {
	return s.userSetting.UpdateUserInteractNotifyRecevieSetting(ctx, req)
}

// UpdateUserAdoptNotifyRecevieSetting 更新用户领养通知
func (s *UserSettingService) UpdateUserAdoptNotifyRecevieSetting(ctx context.Context, req *pb.UpdateUserAdoptNotifyRecevieSettingRequest) (*pb.UpdateUserAdoptNotifyRecevieSettingReply, error) {
	return s.userSetting.UpdateUserAdoptNotifyRecevieSetting(ctx, req)
}

// UpdateUserEmailNotifyRecevieSetting 更新用户邮件通知
func (s *UserSettingService) UpdateUserEmailNotifyRecevieSetting(ctx context.Context, req *pb.UpdateUserEmailNotifyRecevieSettingRequest) (*pb.UpdateUserEmailNotifyRecevieSettingReply, error) {
	return s.userSetting.UpdateUserEmailNotifyRecevieSetting(ctx, req)
}


// GetUserNotifySetting 获取用户通知设置
func (s *UserSettingService) GetUserNotifySetting(ctx context.Context, req *pb.GetUserNotifySettingRequest) (*pb.GetUserNotifySettingReply, error) {
	return s.userSetting.GetUserNotifySetting(ctx, req)
}