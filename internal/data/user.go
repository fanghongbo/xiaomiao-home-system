package data

import (
	"context"
	"fmt"
	"time"
	notificationV1 "xiaomiao-home-system/api/user/notification/v1"
	v1 "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

var (
	ErrBuildForwardRequest = errors.NotFound(v1.ErrorReason_ERR_BUILD_FORWARD_REQUEST.String(), "failed to build forward request")
	ErrForwardRequest      = errors.NotFound(v1.ErrorReason_ERR_FORWARD_REQUEST.String(), "failed to forward request")
)

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserRepo")),
	}
}

// Login 登录接口
func (u *userRepo) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	res := &v1.LoginReply{}

	return res, nil
}

// GetSecretKey 查询加密密钥
func (u *userRepo) GetSecretKey(ctx context.Context, req *v1.GetSecretKeyRequest) (*v1.GetSecretKeyReply, error) {

	res := &v1.GetSecretKeyReply{}

	return res, nil
}

// GetLoginUserInfo 查询用户信息
func (u *userRepo) GetLoginUserInfo(ctx context.Context, req *v1.GetLoginUserInfoRequest) (*v1.GetLoginUserInfoReply, error) {

	res := &v1.GetLoginUserInfoReply{}

	return res, nil
}

// GetAuthMenuList 查询用户信息
func (u *userRepo) GetAuthMenuList(ctx context.Context, req *v1.GetAuthMenuListRequest) (*v1.GetAuthMenuListReply, error) {
	res := &v1.GetAuthMenuListReply{}

	return res, nil
}

// Logout 登出接口
func (u *userRepo) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutReply, error) {
	res := &v1.LogoutReply{}

	return res, nil
}

// GetUsers 查询所有用户
func (u *userRepo) GetUsers(ctx context.Context, req *v1.GetUsersRequest) (*v1.GetUsersReply, error) {
	var users []*v1.UserItem

	result := u.data.db.Model(&User{}).Select("id", "username", "nickname", "status").Where("deleted_flag = ?", 0)

	if req.Name != "" {
		result = result.Where("username LIKE ? or nickname LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}

	result.Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	res := &v1.GetUsersReply{
		Code:    200,
		Success: true,
		Data:    users,
	}

	return res, nil
}

// GetUserByName 查询用户
func (u *userRepo) GetUserByName(ctx context.Context, req *v1.GetUserByNameRequest) (*v1.GetUserByNameReply, error) {
	var user *User

	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	result := u.data.db.Model(&User{}).Select("id", "username", "nickname", "status").Where("username = ?", req.Username).Where("deleted_flag = ?", 0).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.GetUserByNameReply{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Status:   int32(user.Status),
	}, nil
}

// CreateUser 创建用户
func (u *userRepo) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
	id, err := u.data.gid.NextID()
	if err != nil {
		return nil, err
	}

	user := map[string]interface{}{
		"id":           int64(id),
		"username":     req.Username,
		"nickname":     req.Nickname,
		"status":       int(req.Status),
		"remark":       req.Remark,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&User{}).Create(user).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("用户 %s 已经存在", req.Username)
		}
		return nil, err
	}

	if err := u.CreateUserRolesRelation(tx, int64(id), req.Roles); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.CreateUserGroupsRelation(tx, int64(id), req.Groups); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateUserReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateUserByName 更新用户
func (u *userRepo) UpdateUserByName(ctx context.Context, req *v1.UpdateUserByNameRequest) error {
	var user *User

	if req.Username == "" || req.Nickname == "" {
		return fmt.Errorf("username and nickname are required")
	}

	user = &User{
		Nickname: req.Nickname,
		Status:   int(req.Status),
	}

	result := u.data.db.Model(&User{}).Where("username = ?", req.Username).Updates(user).Where("deleted_flag = ?", 0)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserList 查询用户列表
func (u *userRepo) GetUserList(ctx context.Context, req *v1.GetUserListRequest) (*v1.GetUserListReply, error) {
	var users []*v1.UserListItem
	var total int64

	baseQuery := u.data.db.Model(&User{}).Where("deleted_flag = ?", 0)

	if req.Username != "" {
		baseQuery = baseQuery.Where("username LIKE ?", "%"+req.Username+"%")
	}

	if req.Nickname != "" {
		baseQuery = baseQuery.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}

	if req.Remark != "" {
		baseQuery = baseQuery.Where("remark LIKE ?", "%"+req.Remark+"%")
	}

	if req.Status != "" {
		if req.Status == "active" {
			baseQuery = baseQuery.Where("status = ?", 1)
		} else {
			baseQuery = baseQuery.Where("status = ?", 0)
		}
	}

	if req.StartTime > 0 {
		startTimeObj := time.UnixMilli(req.StartTime)
		startTime := time.Date(startTimeObj.Year(), startTimeObj.Month(), startTimeObj.Day(), 0, 0, 0, 0, startTimeObj.Location()).Format("2006-01-02 15:04:05")
		baseQuery = baseQuery.Where("created_time >= ?", startTime)
	}

	if req.EndTime > 0 {
		endTimeObj := time.UnixMilli(req.EndTime)
		endTime := time.Date(endTimeObj.Year(), endTimeObj.Month(), endTimeObj.Day(), 23, 59, 59, 0, endTimeObj.Location()).Format("2006-01-02 15:04:05")
		baseQuery = baseQuery.Where("created_time <= ?", endTime)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	result := baseQuery.Select("id", "username", "nickname", "status", "remark", "created_time", "updated_time").
		Offset(int((req.Page - 1) * req.PageSize)).
		Limit(int(req.PageSize))

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int64
			username    string
			nickname    string
			status      int
			remark      string
			createdTime time.Time
			updatedTime time.Time
		)

		if err := rows.Scan(&id, &username, &nickname, &status, &remark, &createdTime, &updatedTime); err != nil {
			return nil, err
		}

		roles, err := u.GetUserRolesRelation(id)
		if err != nil {
			return nil, err
		}

		groups, err := u.GetUserGroupsRelation(id)
		if err != nil {
			return nil, err
		}

		users = append(users, &v1.UserListItem{
			Id:          id,
			Username:    username,
			Nickname:    nickname,
			Status:      int32(status),
			Roles:       roles,
			Groups:      groups,
			Remark:      remark,
			CreatedTime: createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &v1.GetUserListReply{
		Code: 200, Success: true,
		Data: &v1.UserList{
			Items: users,
			Total: total,
		},
	}, nil
}

// UpdateUser 更新用户
func (u *userRepo) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	updates := map[string]interface{}{
		"username":     req.Username,
		"nickname":     req.Nickname,
		"status":       int(req.Status),
		"remark":       req.Remark,
		"updated_time": time.Now(),
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&User{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updates).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("用户 %s 已经存在", req.Username)
		}
		return nil, err
	}

	if err := u.UpdateUserRolesRelation(tx, int64(req.Id), req.Roles); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.UpdateUserGroupsRelation(tx, int64(req.Id), req.Groups); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.UpdateUserReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeleteUser 删除用户
func (u *userRepo) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()

	result := tx.Model(&User{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := u.DeleteUserRolesRelation(tx, req.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.DeleteUserGroupsRelation(tx, req.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.DeleteUserReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdateUserStatus 更新用户状态
func (u *userRepo) UpdateUserStatus(ctx context.Context, req *v1.UpdateUserStatusRequest) (*v1.UpdateUserStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&User{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateUserStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// CreateUserRolesRelation 创建用户角色关系
func (u *userRepo) CreateUserRolesRelation(tx *gorm.DB, userId int64, roles []int64) error {
	if userId == 0 {
		return fmt.Errorf("user id is required")
	}

	if len(roles) == 0 {
		return nil
	}

	userRoles := make([]map[string]interface{}, 0)

	for _, role := range roles {
		relationId, err := u.data.gid.NextID()
		if err != nil {
			return err
		}

		userRoles = append(userRoles, map[string]interface{}{
			"id":           int64(relationId),
			"user_id":      userId,
			"role_id":      role,
			"created_time": time.Now(),
			"updated_time": time.Now(),
		})
	}

	if len(userRoles) == 0 {
		return nil
	}

	return tx.Model(&UserRoleRelation{}).Create(userRoles).Error
}

// CreateUserGroupsRelation 创建用户组关系
func (u *userRepo) CreateUserGroupsRelation(tx *gorm.DB, userId int64, groups []int64) error {
	if userId == 0 {
		return fmt.Errorf("user id is required")
	}

	if len(groups) == 0 {
		return nil
	}

	userGroups := make([]map[string]interface{}, 0)

	for _, group := range groups {
		relationId, err := u.data.gid.NextID()
		if err != nil {
			return err
		}

		userGroups = append(userGroups, map[string]interface{}{
			"id":           int64(relationId),
			"user_id":      userId,
			"group_id":     group,
			"created_time": time.Now(),
			"updated_time": time.Now(),
		})
	}

	if len(userGroups) == 0 {
		return nil
	}

	return tx.Model(&UserUserGroupRelation{}).Create(userGroups).Error
}

// UpdateUserRolesRelation 更新用户角色关系
func (u *userRepo) UpdateUserRolesRelation(tx *gorm.DB, userId int64, roles []int64) error {
	if userId == 0 {
		return fmt.Errorf("user id is required")
	}

	// 删除用户角色
	if err := u.DeleteUserRolesRelation(tx, userId); err != nil {
		return err
	}

	if err := u.CreateUserRolesRelation(tx, userId, roles); err != nil {
		return err
	}

	return nil
}

// UpdateUserGroupsRelation 更新用户组关系
func (u *userRepo) UpdateUserGroupsRelation(tx *gorm.DB, userId int64, groups []int64) error {
	if userId == 0 {
		return fmt.Errorf("user id is required")
	}

	// 删除用户组
	if err := u.DeleteUserGroupsRelation(tx, userId); err != nil {
		return err
	}

	if err := u.CreateUserGroupsRelation(tx, userId, groups); err != nil {
		return err
	}

	return nil
}

// DeleteUserRolesRelation 删除用户角色关系
func (u *userRepo) DeleteUserRolesRelation(tx *gorm.DB, userId int64) error {
	if userId == 0 {
		return fmt.Errorf("user id is required")
	}

	updates := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	result := tx.Model(&UserRoleRelation{}).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteUserGroupsRelation 删除用户组关系
func (u *userRepo) DeleteUserGroupsRelation(tx *gorm.DB, userId int64) error {
	if userId == 0 {
		return fmt.Errorf("userId id is required")
	}

	updates := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	result := tx.Model(&UserUserGroupRelation{}).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserRolesRelation 获取用户角色关系
func (u *userRepo) GetUserRolesRelation(userId int64) ([]int64, error) {
	var roles []int64

	if err := u.data.db.Model(&UserRoleRelation{}).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Pluck("role_id", &roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserGroupsRelation 获取用户组关系
func (u *userRepo) GetUserGroupsRelation(userId int64) ([]int64, error) {
	var groups []int64

	if err := u.data.db.Model(&UserUserGroupRelation{}).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Pluck("group_id", &groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

// GetUserId 查询当前用户ID
func (u *userRepo) GetUserId(ctx context.Context) (int64, error) {
	username, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return 0, err
	}

	var user User
	if err := u.data.db.Model(&User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}

	return user.Id, nil
}

// GetUserNotifications 查询用户消息
func (u *userRepo) GetUserNotifications(ctx context.Context, req *v1.GetUserNotificationsRequest) (*v1.GetUserNotificationsReply, error) {
	userId, err := u.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	var notifications []*v1.UserNotificationItem

	result := u.data.db.Model(&UserNotification{}).Select("id", "notification_name", "notification_type", "content", "created_time").Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Where("status = ?", 0).Order("created_time DESC").Limit(10)

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id               int64
			notificationName string
			notificationType string
			content          string
			createdTime      time.Time
		)

		if err := rows.Scan(&id, &notificationName, &notificationType, &content, &createdTime); err != nil {
			return nil, err
		}

		notifications = append(notifications, &v1.UserNotificationItem{
			Id:      id,
			Title:   notificationName,
			Type:    notificationType,
			Avatar:  "",
			Date:    createdTime.Format("2006-01-02 15:04:05"),
			Message: content,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetUserNotificationsReply{
		Code: 200, Success: true, Data: notifications,
	}, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (u *userRepo) UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	userNotificationRepo := NewUserNotificationRepo(u.data, u.log.Logger())
	_, err := userNotificationRepo.UpdateUserNotificationStatus(ctx, &notificationV1.UpdateUserNotificationStatusRequest{
		Id:     req.Id,
		Status: req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserNotificationStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// UpdateUserAllNotificationStatus 更新用户所有消息状态
func (u *userRepo) UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error) {

	userNotificationRepo := NewUserNotificationRepo(u.data, u.log.Logger())
	_, err := userNotificationRepo.UpdateUserAllNotificationStatus(ctx, &notificationV1.UpdateUserAllNotificationStatusRequest{
		Status: req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserAllNotificationStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}
