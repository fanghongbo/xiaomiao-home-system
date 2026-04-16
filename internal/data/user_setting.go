package data

import (
	"context"
	"fmt"
	"regexp"
	"time"
	v1 "xiaomiao-home-system/api/user/setting/v1"
	"xiaomiao-home-system/third_party/password"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

const userMaxSettingUpdateCount = 10
const userSettingUpdateCountTTL = 8 * time.Hour

const redisKeyUserBaseSettingUpdateCount = "user:base:setting:update:count"
const redisKeyPasswordUpdateCount = "user:password:update:count"
const redisKeySystemNotifySettingUpdateCount = "user:system:notify:setting:update:count"
const redisKeyInteractNotifySettingUpdateCount = "user:interact:notify:setting:update:count"
const redisKeyAdoptNotifySettingUpdateCount = "user:adopt:notify:setting:update:count"
const redisKeyEmailNotifySettingUpdateCount = "user:email:notify:setting:update:count"

type userSettingRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserSettingRepo .
func NewUserSettingRepo(data *Data, logger log.Logger) biz.UserSettingRepo {
	return &userSettingRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserSettingRepo")),
	}
}

// checkUserSettingUpdateCountLimit 检查用户设置更新次数限制
func (u *userSettingRepo) checkUserSettingUpdateCountLimit(ctx context.Context, userId int64, countKeyPrefix string) error {
	key := fmt.Sprintf("%s:%d", countKeyPrefix, userId)
	n, err := u.data.cache.Incr(ctx, key).Result()
	if err != nil {
		u.log.Errorf("increase notify setting update count failed, key=%s: %v", countKeyPrefix, err)
		return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if n == 1 {
		if err := u.data.cache.Expire(ctx, key, userSettingUpdateCountTTL).Err(); err != nil {
			u.log.Errorf("set notify setting update count ttl failed, key=%s: %v", countKeyPrefix, err)
			return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if n > userMaxSettingUpdateCount {
		return errors.BadRequest(v1.ErrorReason_ERR_TOO_MANY_REQUEST.String(), "今日已达最大修改次数, 请明日再试")
	}

	return nil
}

// UpdateUserBaseSetting 更新用户基础更新
func (u *userSettingRepo) UpdateUserBaseSetting(ctx context.Context, req *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if ok, word := u.data.risk.Validate(req.Nickname); !ok {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), fmt.Sprintf("昵称包含敏感词: %s, 请修改昵称", word))
	}

	if ok, word := u.data.risk.Validate(req.Address); !ok {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), fmt.Sprintf("地址包含敏感词: %s, 请修改地址", word))
	}

	if ok, word := u.data.risk.Validate(req.Signature); !ok {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), fmt.Sprintf("个人签名包含敏感词: %s, 请修改个人签名", word))
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeyUserBaseSettingUpdateCount); err != nil {
		return nil, err
	}

	userInfo := map[string]interface{}{
		"nickname":    req.Nickname,
		"gender":      req.Gender,
		"birthday":    req.Birthday,
		"province_id": req.ProvinceId,
		"city_id":     req.CityId,
		"address":     req.Address,
		"signature":   req.Signature,
	}

	if err := u.data.db.Table("t_user").Where("id = ?", userId).Updates(userInfo).Error; err != nil {
		u.log.Errorf("update user base setting failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserBaseSettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// CheckPassword 校验密码：必须包含 字母+数字+特殊符号，长度6-32
func (u *userSettingRepo) CheckPassword(password string) bool {
	if len(password) < 6 || len(password) > 32 {
		return false
	}

	// 包含字母
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	// 包含数字
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 包含特殊符号
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	return hasLetter && hasDigit && hasSpecial
}

// UpdateUserPassword 更新用户密码
func (u *userSettingRepo) UpdateUserPassword(ctx context.Context, req *v1.UpdateUserPasswordRequest) (*v1.UpdateUserPasswordReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeyPasswordUpdateCount); err != nil {
		return nil, err
	}

	if !u.CheckPassword(req.Password) {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "密码必须包含字母+数字+特殊符号，长度6-32")
	}

	salt, err := password.NewSalt(10)
	if err != nil {
		u.log.Errorf("generate salt failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	passwordHash, err := password.New(req.Password, salt)
	if err != nil {
		u.log.Errorf("generate password hash failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.data.db.Table("t_user_password").Where("user_id = ?", userId).Updates(map[string]interface{}{
		"password": passwordHash,
		"salt":     salt,
	}).Error; err != nil {
		u.log.Errorf("update user password failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserPasswordReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

func (u *userSettingRepo) upsertNotifySetting(userId int64, settingName string, enable int32) error {
	var setting UserSetting
	err := u.data.db.Model(&UserSetting{}).
		Where("user_id = ?", userId).
		Where("name = ?", settingName).
		Where("deleted_flag = ?", 0).
		First(&setting).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		settingId, e := u.data.gid.NextID()
		if e != nil {
			return e
		}

		return u.data.db.Model(&UserSetting{}).Create(map[string]interface{}{
			"id":      settingId,
			"user_id": userId,
			"name":    settingName,
			"value":   enable,
		}).Error
	}

	return u.data.db.Model(&UserSetting{}).
		Where("id = ?", setting.Id).
		Updates(map[string]interface{}{"value": enable}).Error
}

// UpdateUserSystemNotifySetting 更新用户系统通知
func (u *userSettingRepo) UpdateUserSystemNotifySetting(ctx context.Context, req *v1.UpdateUserSystemNotifySettingRequest) (*v1.UpdateUserSystemNotifySettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeySystemNotifySettingUpdateCount); err != nil {
		return nil, err
	}

	if err := u.upsertNotifySetting(userId, "system_notify_receive", req.Enable); err != nil {
		u.log.Errorf("set receive user system notify failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserSystemNotifySettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// UpdateUserInteractNotifySetting 更新用户互动通知
func (u *userSettingRepo) UpdateUserInteractNotifySetting(ctx context.Context, req *v1.UpdateUserInteractNotifySettingRequest) (*v1.UpdateUserInteractNotifySettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeyInteractNotifySettingUpdateCount); err != nil {
		return nil, err
	}

	if err := u.upsertNotifySetting(userId, "interact_notify_receive", req.Enable); err != nil {
		u.log.Errorf("set receive user interact notify failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserInteractNotifySettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// UpdateUserAdoptNotifySetting 更新用户领养通知
func (u *userSettingRepo) UpdateUserAdoptNotifySetting(ctx context.Context, req *v1.UpdateUserAdoptNotifySettingRequest) (*v1.UpdateUserAdoptNotifySettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeyAdoptNotifySettingUpdateCount); err != nil {
		return nil, err
	}

	if err := u.upsertNotifySetting(userId, "adopt_notify_receive", req.Enable); err != nil {
		u.log.Errorf("set receive user adopt notify failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserAdoptNotifySettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// UpdateUserEmailNotifySetting 更新用户邮件通知
func (u *userSettingRepo) UpdateUserEmailNotifySetting(ctx context.Context, req *v1.UpdateUserEmailNotifySettingRequest) (*v1.UpdateUserEmailNotifySettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.checkUserSettingUpdateCountLimit(ctx, userId, redisKeyEmailNotifySettingUpdateCount); err != nil {
		return nil, err
	}

	if err := u.upsertNotifySetting(userId, "email_notify_receive", req.Enable); err != nil {
		u.log.Errorf("set receive user email notify failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserEmailNotifySettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// GetUserNotifySetting 获取用户通知设置
func (u *userSettingRepo) GetUserNotifySetting(ctx context.Context, req *v1.GetUserNotifySettingRequest) (*v1.GetUserNotifySettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	var systemVal int32
	if err := u.data.db.Model(&UserSetting{}).Select("value").Where("user_id = ?", userId).Where("name = ?", "system_notify_receive").Where("deleted_flag = ?", 0).First(&systemVal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			systemVal = 1
		} else {
			u.log.Errorf("get user notify setting failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	var interactVal int32
	if err := u.data.db.Model(&UserSetting{}).Select("value").Where("user_id = ?", userId).Where("name = ?", "interact_notify_receive").Where("deleted_flag = ?", 0).First(&interactVal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			interactVal = 1
		} else {
			u.log.Errorf("get user notify setting failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	var adoptVal int32
	if err := u.data.db.Model(&UserSetting{}).Select("value").Where("user_id = ?", userId).Where("name = ?", "adopt_notify_receive").Where("deleted_flag = ?", 0).First(&adoptVal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adoptVal = 1
		} else {
			u.log.Errorf("get user notify setting failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	var emailVal int32
	if err := u.data.db.Model(&UserSetting{}).Select("value").Where("user_id = ?", userId).Where("name = ?", "email_notify_receive").Where("deleted_flag = ?", 0).First(&emailVal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			emailVal = 1
		} else {
			u.log.Errorf("get user notify setting failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	return &v1.GetUserNotifySettingReply{
		Code:    200,
		Message: "查询成功",
		Success: true,
		Data: &v1.UserNotifySettingInfo{
			System:   systemVal,
			Interact: interactVal,
			Adopt:    adoptVal,
			Email:    emailVal,
		},
	}, nil
}
