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

// UpdateUserSystemNotifySetting 更新用户系统通知
func (s *UserSettingService) UpdateUserSystemNotifySetting(ctx context.Context, req *pb.UpdateUserSystemNotifySettingRequest) (*pb.UpdateUserSystemNotifySettingReply, error) {
	return s.userSetting.UpdateUserSystemNotifySetting(ctx, req)
}

// UpdateUserInteractNotifySetting 更新用户互动通知
func (s *UserSettingService) UpdateUserInteractNotifySetting(ctx context.Context, req *pb.UpdateUserInteractNotifySettingRequest) (*pb.UpdateUserInteractNotifySettingReply, error) {
	return s.userSetting.UpdateUserInteractNotifySetting(ctx, req)
}

// UpdateUserAdoptNotifySetting 更新用户领养通知
func (s *UserSettingService) UpdateUserAdoptNotifySetting(ctx context.Context, req *pb.UpdateUserAdoptNotifySettingRequest) (*pb.UpdateUserAdoptNotifySettingReply, error) {
	return s.userSetting.UpdateUserAdoptNotifySetting(ctx, req)
}

// UpdateUserEmailNotifySetting 更新用户邮件通知
func (s *UserSettingService) UpdateUserEmailNotifySetting(ctx context.Context, req *pb.UpdateUserEmailNotifySettingRequest) (*pb.UpdateUserEmailNotifySettingReply, error) {
	return s.userSetting.UpdateUserEmailNotifySetting(ctx, req)
}

// GetUserNotifySetting 获取用户通知设置
func (s *UserSettingService) GetUserNotifySetting(ctx context.Context, req *pb.GetUserNotifySettingRequest) (*pb.GetUserNotifySettingReply, error) {
	return s.userSetting.GetUserNotifySetting(ctx, req)
}
