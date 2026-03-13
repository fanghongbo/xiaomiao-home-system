package data

import "time"

type UserGroup struct {
	Id          int64      `gorm:"column:id"`
	GroupName   string     `gorm:"column:group_name"`
	Remark      string     `gorm:"column:remark"`
	Status      int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedUser string     `gorm:"column:created_user"`
	UpdatedUser string     `gorm:"column:updated_user"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u UserGroup) TableName() string {
	return "t_user_group"
}

type User struct {
	Id          int64      `gorm:"column:id"`
	Username    string     `gorm:"column:username"`
	Nickname    string     `gorm:"column:nickname"`
	Remark      string     `gorm:"column:remark"`
	Status      int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u User) TableName() string {
	return "t_user"
}

type UserUserGroupRelation struct {
	Id          int64      `gorm:"column:id"`
	UserId      int64      `gorm:"column:user_id"`
	RoleId      int64      `gorm:"column:role_id"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u UserUserGroupRelation) TableName() string {
	return "t_user_user_group"
}

type UserRoleRelation struct {
	Id          int64      `gorm:"column:id"`
	UserId      int64      `gorm:"column:user_id"`
	RoleId      int64      `gorm:"column:role_id"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u UserRoleRelation) TableName() string {
	return "t_user_role"
}

type Role struct {
	Id          int64      `gorm:"column:id"`
	RoleName    string     `gorm:"column:role_name"`
	Remark      string     `gorm:"column:remark"`
	Status      int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedUser string     `gorm:"column:created_user"`
	UpdatedUser string     `gorm:"column:updated_user"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u Role) TableName() string {
	return "t_role"
}

type RolePermissionRelation struct {
	Id             int64      `gorm:"column:id"`
	RoleId         int64      `gorm:"column:role_id"`
	PermissionCode string     `gorm:"column:permission_code"`
	CreatedTime    time.Time  `gorm:"column:created_time"`
	UpdatedTime    time.Time  `gorm:"column:updated_time"`
	DeletedFlag    int        `gorm:"column:deleted_flag"`
	DeletedTime    *time.Time `gorm:"column:deleted_time"`
}

func (u RolePermissionRelation) TableName() string {
	return "t_role_permission"
}

type DnsDomain struct {
	Id          int64      `gorm:"column:id"`
	DomainName  string     `gorm:"column:domain_name"`
	Remark      string     `gorm:"column:remark"`
	Status      int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedUser string     `gorm:"column:created_user"`
	UpdatedUser string     `gorm:"column:updated_user"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsDomain) TableName() string {
	return "t_dns_domain"
}

type DnsProvider struct {
	Id           int64      `gorm:"column:id"`
	ProviderName string     `gorm:"column:provider_name"`
	ProviderType string     `gorm:"column:provider_type"`
	ApiKey       string     `gorm:"column:api_key"`
	ApiSecret    string     `gorm:"column:api_secret"`
	Remark       string     `gorm:"column:remark"`
	Status       int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedUser  string     `gorm:"column:created_user"`
	UpdatedUser  string     `gorm:"column:updated_user"`
	CreatedTime  time.Time  `gorm:"column:created_time"`
	UpdatedTime  time.Time  `gorm:"column:updated_time"`
	DeletedFlag  int        `gorm:"column:deleted_flag"`
	DeletedTime  *time.Time `gorm:"column:deleted_time"`
}

func (u DnsProvider) TableName() string {
	return "t_dns_provider"
}

type UserGroupRoleRelation struct {
	Id          int64      `gorm:"column:id"`
	GroupId     int64      `gorm:"column:group_id"`
	RoleId      int64      `gorm:"column:role_id"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u UserGroupRoleRelation) TableName() string {
	return "t_user_group_role"
}

// DnsDomainProviderRelation  DNS域名与服务商关系
type DnsDomainProviderRelation struct {
	Id          int64      `gorm:"column:id"`
	DomainId    int64      `gorm:"column:domain_id"`
	ProviderId  int64      `gorm:"column:provider_id"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsDomainProviderRelation) TableName() string {
	return "t_dns_domain_provider"
}

// DnsRecord DNS记录
type DnsRecord struct {
	Id          int64      `gorm:"column:id"`
	DomainId    int64      `gorm:"column:domain_id"`
	VersionNum  int64      `gorm:"column:version_num"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsRecord) TableName() string {
	return "t_dns_record"
}

// DnsRecordVersion  DNS记录版本
type DnsRecordVersion struct {
	Id           int64      `gorm:"column:id"`
	DomainId     int64      `gorm:"column:domain_id"`
	RecordId     int64      `gorm:"column:record_id"`
	VersionNum   int64      `gorm:"column:version_num"`
	Host         string     `gorm:"column:host"`
	RecordType   string     `gorm:"column:type"`
	Priority     int64      `gorm:"column:priority"`
	Line         string     `gorm:"column:line"`
	Ttl          int64      `gorm:"column:ttl"`
	Value        string     `gorm:"column:value"`
	Weight       int32      `gorm:"column:weight"`
	LoadStrategy string     `gorm:"column:load_strategy"`
	Remark       string     `gorm:"column:remark"`
	Status       int        `gorm:"column:status"` // 0: 正常 1: 禁用
	CreatedTime  time.Time  `gorm:"column:created_time"`
	UpdatedTime  time.Time  `gorm:"column:updated_time"`
	DeletedFlag  int        `gorm:"column:deleted_flag"` // 0: 未删除 1: 已删除
	DeletedTime  *time.Time `gorm:"column:deleted_time"`
}

func (u DnsRecordVersion) TableName() string {
	return "t_dns_record_version"
}

// DnsTask DNS任务
type DnsTask struct {
	Id          int64      `gorm:"column:id"`
	TaskName    string     `gorm:"column:task_name"`
	DomainId    int64      `gorm:"column:domain_id"`
	TaskType    string     `gorm:"column:task_type"`
	StartTime   time.Time  `gorm:"column:start_time"`
	FinishTime  time.Time  `gorm:"column:finish_time"`
	DurationMs  int64      `gorm:"column:duration_ms"`
	Remark      string     `gorm:"column:remark"`
	Status      int        `gorm:"column:status"` // 0: pending, 1: processing, 2: success, 3: failed, 4: cancelled
	Operator    string     `gorm:"column:operator"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"`
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsTask) TableName() string {
	return "t_dns_task"
}

// DnsTaskRecord DNS任务记录
type DnsTaskRecord struct {
	Id              int64      `gorm:"column:id"`
	TaskId          int64      `gorm:"column:task_id"`
	ProviderId      int64      `gorm:"column:provider_id"`
	RecordId        int64      `gorm:"column:record_id"`
	VersionNum      int64      `gorm:"column:version_num"`
	NewVersionNum   int64      `gorm:"column:new_version_num"`
	Status          int        `gorm:"column:status"` // 0: pending, 1: processing, 2: success, 3: failed, 4: cancelled, 5: retry_waiting
	ExecuteTimes    int32      `gorm:"column:execute_times"`
	MaxRetryTimes   int32      `gorm:"column:max_retry_times"`
	NextExecuteTime time.Time  `gorm:"column:next_execute_time"`
	ErrorMsg        string     `gorm:"column:error_msg"`
	StartTime       time.Time  `gorm:"column:start_time"`
	FinishTime      time.Time  `gorm:"column:finish_time"`
	DurationMs      int64      `gorm:"column:duration_ms"`
	Operator        string     `gorm:"column:operator"`
	CreatedTime     time.Time  `gorm:"column:created_time"`
	UpdatedTime     time.Time  `gorm:"column:updated_time"`
	DeletedFlag     int        `gorm:"column:deleted_flag"`
	DeletedTime     *time.Time `gorm:"column:deleted_time"`
}

func (u DnsTaskRecord) TableName() string {
	return "t_dns_task_record"
}

// UserNotification 用户消息
type UserNotification struct {
	Id               int64      `gorm:"column:id"`
	UserId           int64      `gorm:"column:user_id"`
	NotificationName string     `gorm:"column:notification_name"`
	NotificationType string     `gorm:"column:notification_type"`
	Content          string     `gorm:"column:content"`
	Status           int        `gorm:"column:status"` // 0: 未读 1: 已读
	ReadTime         time.Time  `gorm:"column:read_time"`
	CreatedTime      time.Time  `gorm:"column:created_time"`
	UpdatedTime      time.Time  `gorm:"column:updated_time"`
	DeletedFlag      int        `gorm:"column:deleted_flag"`
	DeletedTime      *time.Time `gorm:"column:deleted_time"`
}

func (u UserNotification) TableName() string {
	return "t_user_notification"
}

// DnsAnalysis DNS分析
type DnsAnalysis struct {
	Id          int64      `gorm:"column:id"`
	UserId      int64      `gorm:"column:user_id"`
	DomainId    int64      `gorm:"column:domain_id"`
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"` // 0: 未删除 1: 已删除
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsAnalysis) TableName() string {
	return "t_dns_analysis"
}

// DnsAnalysisRecord DNS分析记录
type DnsAnalysisRecord struct {
	Id          int64      `gorm:"column:id"`
	AnalysisId  int64      `gorm:"column:analysis_id"`
	SourceId    int64      `gorm:"column:source_id"`
	TargetId    int64      `gorm:"column:target_id"`
	DiffType    int        `gorm:"column:diff_type"`   // 0: 本地和目标服务商对比, 1: 服务商之间对比
	DiffStatus  int        `gorm:"column:diff_status"` // 0: 源记录存在 1: 源记录不存在 2: 源和目标差异
	CreatedTime time.Time  `gorm:"column:created_time"`
	UpdatedTime time.Time  `gorm:"column:updated_time"`
	DeletedFlag int        `gorm:"column:deleted_flag"` // 0: 未删除 1: 已删除
	DeletedTime *time.Time `gorm:"column:deleted_time"`
}

func (u DnsAnalysisRecord) TableName() string {
	return "t_dns_analysis_record"
}
