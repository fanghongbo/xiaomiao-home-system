package data

import "time"

type User struct {
	Id          int64     `gorm:"column:id"`
	Username    string    `gorm:"column:username"`
	Nickname    string    `gorm:"column:nickname"`
	Password    string    `gorm:"column:password"`
	Salt        string    `gorm:"column:salt"`
	Avatar      string    `gorm:"column:avatar"`
	Remark      string    `gorm:"column:remark"`
	Status      int       `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u User) TableName() string {
	return "t_user"
}

type UserIdentity struct {
	Id           int64     `gorm:"column:id"`
	UserId       int64     `gorm:"column:user_id"`
	IdentityType string    `gorm:"column:identity_type"`
	IdentityId   string    `gorm:"column:identity_id"`
	VerifiedFlag int       `gorm:"column:verified_flag"`
	Remark       string    `gorm:"column:remark"`
	CreatedTime  time.Time `gorm:"column:created_time"`
	UpdatedTime  time.Time `gorm:"column:updated_time"`
	DeletedFlag  int       `gorm:"column:deleted_flag"`
	DeletedTime  time.Time `gorm:"column:deleted_time"`
}

func (u UserIdentity) TableName() string {
	return "t_user_identity"
}

type UserRoles struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	RoleId      int64     `gorm:"column:role_id"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u UserRoles) TableName() string {
	return "t_user_role"
}

type Role struct {
	Id          int64     `gorm:"column:id"`
	RoleName    string    `gorm:"column:role_name"`
	Remark      string    `gorm:"column:remark"`
	Status      int       `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedUser string    `gorm:"column:created_user"`
	UpdatedUser string    `gorm:"column:updated_user"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u Role) TableName() string {
	return "t_role"
}

type RolePermission struct {
	Id             int64      `gorm:"column:id"`
	RoleId         int64      `gorm:"column:role_id"`
	PermissionCode string     `gorm:"column:permission_code"`
	CreatedTime    time.Time  `gorm:"column:created_time"`
	UpdatedTime    time.Time  `gorm:"column:updated_time"`
	DeletedFlag    int        `gorm:"column:deleted_flag"`
	DeletedTime    *time.Time `gorm:"column:deleted_time"`
}

func (u RolePermission) TableName() string {
	return "t_role_permission"
}

type UserNotification struct {
	Id               int64      `gorm:"column:id"`
	UserId           int64      `gorm:"column:user_id"`
	NotificationName string     `gorm:"column:notification_name"`
	Content          string     `gorm:"column:content"`
	Status           int        `gorm:"column:status"` // 0: 未读 1: 已读
	CreatedTime      time.Time  `gorm:"column:created_time"`
	UpdatedTime      time.Time  `gorm:"column:updated_time"`
	DeletedFlag      int        `gorm:"column:deleted_flag"`
	DeletedTime      *time.Time `gorm:"column:deleted_time"`
}

func (u UserNotification) TableName() string {
	return "t_user_notification"
}

type SystemSetting struct {
	Id          int64     `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Value       string    `gorm:"column:value"`
	Remark      string    `gorm:"column:remark"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u SystemSetting) TableName() string {
	return "t_system_setting"
}
