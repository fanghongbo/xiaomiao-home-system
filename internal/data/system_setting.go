package data

import (
	"context"
	v1 "xiaomiao-home-system/api/system/setting/v1"
	"xiaomiao-home-system/internal/biz"

	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type systemSettingRepo struct {
	data *Data
	log  *log.Helper
}

// NewSystemSettingRepo .
func NewSystemSettingRepo(data *Data, logger log.Logger) biz.SystemSettingRepo {
	return &systemSettingRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "SystemSettingRepo")),
	}
}

// GetSystemSettingByName 查询系统设置
func (u *systemSettingRepo) GetSystemSettingByName(ctx context.Context, name string) (*v1.SystemSettingInfo, error) {
	var systemSetting *v1.SystemSettingInfo

	if err := u.data.db.Model(&SystemSetting{}).Where("name = ?", name).Where("deleted_flag = ?", 0).First(&systemSetting).Error; err != nil {
		return nil, err
	}

	return &v1.SystemSettingInfo{
		Id:    systemSetting.Id,
		Name:  systemSetting.Name,
		Value: systemSetting.Value,
	}, nil
}

// GetSiteName 查询站点名称
func (u *systemSettingRepo) GetSiteName(ctx context.Context) (string, error) {
	var (
		setting *v1.SystemSettingInfo
		err     error
	)

	if setting, err = u.GetSystemSettingByName(ctx, "site_name"); err != nil {
		return "", err
	}

	return setting.Value, nil
}

// GetSiteUrl 查询站点地址
func (u *systemSettingRepo) GetSiteUrl(ctx context.Context) (string, error) {
	var (
		setting *v1.SystemSettingInfo
		err     error
	)

	if setting, err = u.GetSystemSettingByName(ctx, "site_url"); err != nil {
		return "", err
	}

	return setting.Value, nil
}

// GetBaseSetting 查询基础设置
func (u *systemSettingRepo) GetBaseSetting(ctx context.Context) (*v1.BaseSettingInfo, error) {
	var (
		result *v1.BaseSettingInfo
		err    error
	)

	result = &v1.BaseSettingInfo{}

	if result.SiteName, err = u.GetSiteName(ctx); err != nil {
		return nil, err
	}

	if result.SiteUrl, err = u.GetSiteUrl(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

// GetLoginValidPeriod 查询登录有效期设置
func (u *systemSettingRepo) GetLoginValidPeriod(ctx context.Context) (int64, error) {
	var (
		setting    *v1.SystemSettingInfo
		defaultVal int64
		err        error
	)

	// 默认8小时
	defaultVal = 8

	if setting, err = u.GetSystemSettingByName(ctx, "login_valid_period"); err != nil {
		return defaultVal, err
	}

	return utils.StrToInt64(setting.Value)
}

// GetSessionIdlePeriod 查询不活跃重登时间
func (u *systemSettingRepo) GetSessionIdlePeriod(ctx context.Context) (int64, error) {
	var (
		setting *v1.SystemSettingInfo
		err     error
	)

	if setting, err = u.GetSystemSettingByName(ctx, "session_idle_period"); err != nil {
		return 0, err
	}

	return utils.StrToInt64(setting.Value)
}

// GetInvalidLoginCount 查询限制登录错误次数
func (u *systemSettingRepo) GetInvalidLoginCount(ctx context.Context) (int64, error) {
	var (
		setting    *v1.SystemSettingInfo
		defaultVal int64
		err        error
	)

	// 默认3次
	defaultVal = 3

	if setting, err = u.GetSystemSettingByName(ctx, "invalid_login_count"); err != nil {
		return defaultVal, err
	}

	return utils.StrToInt64(setting.Value)
}

// GetLoginBanDuration 查询限制登录时间间隔
func (u *systemSettingRepo) GetLoginBanDuration(ctx context.Context) (int64, error) {
	var (
		setting    *v1.SystemSettingInfo
		defaultVal int64
		err        error
	)

	// 默认5分钟
	defaultVal = 5

	if setting, err = u.GetSystemSettingByName(ctx, "login_ban_duration"); err != nil {
		return defaultVal, err
	}

	return utils.StrToInt64(setting.Value)
}

// GetLoginSetting 查询登录设置
func (u *systemSettingRepo) GetLoginSetting(ctx context.Context) (*v1.LoginSettingInfo, error) {
	var (
		result *v1.LoginSettingInfo
		err    error
	)

	result = &v1.LoginSettingInfo{}

	if result.LoginValidPeriod, err = u.GetLoginValidPeriod(ctx); err != nil {
		return nil, err
	}

	if result.SessionIdlePeriod, err = u.GetSessionIdlePeriod(ctx); err != nil {
		return nil, err
	}

	if result.InvalidLoginCount, err = u.GetInvalidLoginCount(ctx); err != nil {
		return nil, err
	}

	if result.LoginBanDuration, err = u.GetLoginBanDuration(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

// GetMfaAuth 查询Mfa全局认证状态
func (u *systemSettingRepo) GetMfaAuth(ctx context.Context) (bool, error) {
	var (
		setting *v1.SystemSettingInfo
		val     int64
		err     error
	)

	if setting, err = u.GetSystemSettingByName(ctx, "mfa_auth"); err != nil {
		return false, err
	}

	if val, err = utils.StrToInt64(setting.Value); err != nil {
		return false, err
	}

	return val > 0, nil
}

// GetMfaValidPeriod 查询Mfa验证有效期
func (u *systemSettingRepo) GetMfaValidPeriod(ctx context.Context) (int64, error) {
	var (
		setting    *v1.SystemSettingInfo
		defaultVal int64
		err        error
	)

	// 默认1分钟
	defaultVal = 1

	if setting, err = u.GetSystemSettingByName(ctx, "mfa_valid_period"); err != nil {
		return defaultVal, err
	}

	return utils.StrToInt64(setting.Value)
}

// GetSecuritySetting 查询安全设置
func (u *systemSettingRepo) GetSecuritySetting(ctx context.Context) (*v1.SecuritySettingInfo, error) {
	var (
		result *v1.SecuritySettingInfo
		err    error
	)

	result = &v1.SecuritySettingInfo{}

	if result.MfaAuth, err = u.GetMfaAuth(ctx); err != nil {
		return nil, err
	}

	if result.MfaValidPeriod, err = u.GetMfaValidPeriod(ctx); err != nil {
		return nil, err
	}

	if result.CasValidPeriod, err = u.GetCasValidPeriod(ctx); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCasValidPeriod 查询Cas验证有效期
func (u *systemSettingRepo) GetCasValidPeriod(ctx context.Context) (int64, error) {
	var (
		setting    *v1.SystemSettingInfo
		defaultVal int64
		val        int64
		err        error
	)

	// 默认1分钟
	defaultVal = 1

	if setting, err = u.GetSystemSettingByName(ctx, "cas_valid_period"); err != nil {
		return defaultVal, err
	}

	if val, err = utils.StrToInt64(setting.Value); err != nil {
		return defaultVal, err
	}

	return val, nil
}
