package data

import (
	"time"
	"xiaomiao-home-system/third_party/password"

	"gorm.io/gorm/clause"
)

func InitData(d *Data) error {
	if err := initAdminUser(d); err != nil {
		return err
	}

	if err := initSystemSetting(d); err != nil {
		return err
	}

	return nil
}

func initAdminUser(d *Data) error {
	userId, err := d.gid.NextID()
	if err != nil {
		return err
	}

	salt, err := password.NewSalt(10)
	if err != nil {
		return err
	}

	password, err := password.New("style_admin@2026", salt)
	if err != nil {
		return err
	}

	user := map[string]interface{}{
		"id":           int64(userId),
		"username":     "admin",
		"nickname":     "超级管理员",
		"password":     password,
		"salt":         salt,
		"status":       1,
		"remark":       "超级管理员",
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := d.db.Model(&User{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}, {Name: "deleted_flag"}, {Name: "deleted_time"}},
		DoNothing: true,
	}).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func initSystemSetting(d *Data) error {

	settings := []map[string]interface{}{
		{
			"name":   "site_name",
			"value":  "小猫回家",
			"remark": "当前站点名称",
		},
		{
			"name":   "site_url",
			"value":  "http://localhost:8000",
			"remark": "当前站点地址",
		},
		{
			"name":   "login_valid_period",
			"value":  "8",
			"remark": "登录有效期 (单位:小时)",
		},
		{
			"name":   "session_idle_period",
			"value":  "2",
			"remark": "不活跃重登时间 (单位:分钟)",
		},
		{
			"name":   "login_valid_period",
			"value":  "5",
			"remark": "登陆校验有效期 (单位:分钟)",
		},
		{
			"name":   "invalid_login_count",
			"value":  "3",
			"remark": "登录异常次数",
		},
		{
			"name":   "login_ban_duration",
			"value":  "5",
			"remark": "限制登录时间 (单位:分钟)",
		},
	}

	for _, setting := range settings {
		settingId, err := d.gid.NextID()
		if err != nil {
			return err
		}
		setting["id"] = int64(settingId)
		setting["created_time"] = time.Now()
		setting["updated_time"] = time.Now()

		if err := d.db.Model(&SystemSetting{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "deleted_flag"}, {Name: "deleted_time"}},
			DoNothing: true,
		}).Create(setting).Error; err != nil {
			return err
		}
	}

	return nil
}
