package data

import "time"

type User struct {
	Id          int64     `gorm:"column:id"`
	Nickname    string    `gorm:"column:nickname"`
	Gender      int       `gorm:"column:gender"`
	Birthday    time.Time `gorm:"column:birthday"`
	Avatar      string    `gorm:"column:avatar"`
	Signature   string    `gorm:"column:signature"`
	ProvinceId  int32     `gorm:"column:province_id"`
	CityId      int32     `gorm:"column:city_id"`
	Address     string    `gorm:"column:address"`
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

// UserPassword 账号密码凭据；无密码登录方式的用户可无对应行。
type UserPassword struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	Password    string    `gorm:"column:password"`
	Salt        string    `gorm:"column:salt"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (UserPassword) TableName() string {
	return "t_user_password"
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

type UserInfo struct {
	Id       int64  `gorm:"column:id"`
	Nickname string `gorm:"column:nickname"`
}

type UserSetting struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	Name        string    `gorm:"column:name"`
	Value       string    `gorm:"column:value"`
	Remark      string    `gorm:"column:remark"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u UserSetting) TableName() string {
	return "t_user_setting"
}

type Post struct {
	Id          int64      `gorm:"column:id"`
	Title       string     `gorm:"column:title"`
	PostType    int        `gorm:"column:post_type"`
	ProvinceId  int64      `gorm:"column:province_id"`
	CityId      int64      `gorm:"column:city_id"`
	Address     string     `gorm:"column:address"`
	LostTime    *time.Time `gorm:"column:lost_time"`
	AuditStatus int        `gorm:"column:audit_status"`
	AuditRemark string     `gorm:"column:audit_remark"`
	AuditTime   *time.Time `gorm:"column:audit_time"`
	PostStatus  int        `gorm:"column:post_status"`
	Remark      string     `gorm:"column:remark"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u Post) TableName() string {
	return "t_post"
}

type UserPost struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	PostId      int64     `gorm:"column:post_id"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u UserPost) TableName() string {
	return "t_user_post"
}

type Cat struct {
	Id               int64     `gorm:"column:id"`
	Name             string    `gorm:"column:name"`
	Gender           int       `gorm:"column:gender"`
	CatType          int       `gorm:"column:cat_type"`
	BreedType        int       `gorm:"column:breed_type"`
	Weight           float64   `gorm:"column:weight"`
	Birthday         time.Time `gorm:"column:birthday"`
	NeuterStatus     int       `gorm:"column:neuter_status"`
	HealthStatus     int       `gorm:"column:health_status"`
	HealthDesc       string    `gorm:"column:health_desc"`
	DewormedStatus   int       `gorm:"column:dewormed_status"`
	VaccineStatus    int       `gorm:"column:vaccine_status"`
	VaccineTypes     string    `gorm:"column:vaccine_types"`
	VaccineLastDate  time.Time `gorm:"column:vaccine_last_date"`
	VaccineCertImage string    `gorm:"column:vaccine_cert_image"`
	Remark           string    `gorm:"column:remark"`
	DeletedFlag      int       `gorm:"column:deleted_flag"`
	CreatedTime      time.Time `gorm:"column:created_time"`
	UpdatedTime      time.Time `gorm:"column:updated_time"`
	DeletedTime      time.Time `gorm:"column:deleted_time"`
}

func (u Cat) TableName() string {
	return "t_cat"
}

type UserCat struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	CatId       int64     `gorm:"column:cat_id"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u UserCat) TableName() string {
	return "t_user_cat"
}

type PostCat struct {
	Id          int64     `gorm:"column:id"`
	PostId      int64     `gorm:"column:post_id"`
	CatId       int64     `gorm:"column:cat_id"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	DeletedFlag int       `gorm:"column:deleted_flag"`
	DeletedTime time.Time `gorm:"column:deleted_time"`
}

func (u PostCat) TableName() string {
	return "t_post_cat"
}
