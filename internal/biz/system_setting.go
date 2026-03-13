package biz

import (
	"context"

	v1 "xiaomiao-home-system/api/system/setting/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// SystemSettingRepo is a Greater repo.
type SystemSettingRepo interface {
	// GetSystemSettingByName 查询系统设置
	GetSystemSettingByName(ctx context.Context, name string) (*v1.SystemSettingInfo, error)
	// GetSiteName 查询站点名称
	GetSiteName(ctx context.Context) (string, error)
	// GetSiteUrl 查询站点地址
	GetSiteUrl(ctx context.Context) (string, error)
	// GetBaseSetting 查询基础设置
	GetBaseSetting(ctx context.Context) (*v1.BaseSettingInfo, error)
	// GetLoginSetting 查询登录设置
	GetLoginSetting(ctx context.Context) (*v1.LoginSettingInfo, error)
	// GetMfaAuth 查询Mfa全局认证状态
	GetMfaAuth(ctx context.Context) (bool, error)
	// GetMfaValidPeriod 查询Mfa验证有效期
	GetMfaValidPeriod(ctx context.Context) (int64, error)
	// GetSecuritySetting 查询安全设置
	GetSecuritySetting(ctx context.Context) (*v1.SecuritySettingInfo, error)
	// GetCasValidPeriod 查询Cas验证有效期
	GetCasValidPeriod(ctx context.Context) (int64, error)
	// GetLoginValidPeriod 查询登录有效期
	GetLoginValidPeriod(ctx context.Context) (int64, error)
	// GetSessionIdlePeriod 查询不活跃重登时间
	GetSessionIdlePeriod(ctx context.Context) (int64, error)
	// GetInvalidLoginCount 查询限制登录错误次数
	GetInvalidLoginCount(ctx context.Context) (int64, error)
	// GetLoginBanDuration 查询限制登录时间间隔
	GetLoginBanDuration(ctx context.Context) (int64, error)
}

// SystemSettingUsecase is a User usecase.
type SystemSettingUsecase struct {
	repo SystemSettingRepo
	log  *log.Helper
}

// NewSystemSettingUsecase new a SystemSetting usecase.
func NewSystemSettingUsecase(repo SystemSettingRepo, logger log.Logger) *SystemSettingUsecase {
	return &SystemSettingUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "SystemSettingUsecase"))}
}

// GetSystemSettingByName 查询系统设置
func (u *SystemSettingUsecase) GetSystemSettingByName(ctx context.Context, name string) (*v1.SystemSettingInfo, error) {
	return u.repo.GetSystemSettingByName(ctx, name)
}

// GetSiteName 查询站点名称
func (u *SystemSettingUsecase) GetSiteName(ctx context.Context) (string, error) {
	return u.repo.GetSiteName(ctx)
}

// GetSiteUrl 查询站点地址
func (u *SystemSettingUsecase) GetSiteUrl(ctx context.Context) (string, error) {
	return u.repo.GetSiteUrl(ctx)
}

// GetBaseSetting 查询基础设置
func (u *SystemSettingUsecase) GetBaseSetting(ctx context.Context) (*v1.BaseSettingInfo, error) {
	return u.repo.GetBaseSetting(ctx)
}

// GetLoginSetting 查询登录设置
func (u *SystemSettingUsecase) GetLoginSetting(ctx context.Context) (*v1.LoginSettingInfo, error) {
	return u.repo.GetLoginSetting(ctx)
}

// GetMfaAuth 查询Mfa全局认证状态
func (u *SystemSettingUsecase) GetMfaAuth(ctx context.Context) (bool, error) {
	return u.repo.GetMfaAuth(ctx)
}

// GetMfaValidPeriod 查询Mfa验证有效期
func (u *SystemSettingUsecase) GetMfaValidPeriod(ctx context.Context) (int64, error) {
	return u.repo.GetMfaValidPeriod(ctx)
}

// GetSecuritySetting 查询安全设置
func (u *SystemSettingUsecase) GetSecuritySetting(ctx context.Context) (*v1.SecuritySettingInfo, error) {
	return u.repo.GetSecuritySetting(ctx)
}

// GetCasValidPeriod 查询Cas验证有效期
func (u *SystemSettingUsecase) GetCasValidPeriod(ctx context.Context) (int64, error) {
	return u.repo.GetCasValidPeriod(ctx)
}

// GetLoginValidPeriod 查询登录有效期
func (u *SystemSettingUsecase) GetLoginValidPeriod(ctx context.Context) (int64, error) {
	return u.repo.GetLoginValidPeriod(ctx)
}

// GetSessionIdlePeriod 查询不活跃重登时间
func (u *SystemSettingUsecase) GetSessionIdlePeriod(ctx context.Context) (int64, error) {
	return u.repo.GetSessionIdlePeriod(ctx)
}

// GetInvalidLoginCount 查询限制登录错误次数
func (u *SystemSettingUsecase) GetInvalidLoginCount(ctx context.Context) (int64, error) {
	return u.repo.GetInvalidLoginCount(ctx)
}

// GetLoginBanDuration 查询限制登录时间间隔
func (u *SystemSettingUsecase) GetLoginBanDuration(ctx context.Context) (int64, error) {
	return u.repo.GetLoginBanDuration(ctx)
}
