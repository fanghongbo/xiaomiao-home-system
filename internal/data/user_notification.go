package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/user/notification/v1"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userNotificationRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserNotificationRepo .
func NewUserNotificationRepo(data *Data, logger log.Logger) biz.UserNotificationRepo {
	return &userNotificationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserNotificationRepo")),
	}
}

// GetUserNotificationList 查询用户消息列表
func (u *userNotificationRepo) GetUserNotificationList(ctx context.Context, req *v1.GetUserNotificationListRequest) (*v1.GetUserNotificationListReply, error) {
	userRepo := NewUserRepo(u.data, u.log.Logger())

	userId, err := userRepo.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	var userNotifications []*v1.UserNotificationListItem
	var total int64

	baseQuery := u.data.db.Model(&UserNotification{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId)

	if req.Name != "" {
		baseQuery = baseQuery.Where("notification_name LIKE ?", "%"+req.Name+"%")
	}

	if req.Content != "" {
		baseQuery = baseQuery.Where("content LIKE ?", "%"+req.Content+"%")
	}

	if req.Status != "" {
		if req.Status == "read" {
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

	if err := baseQuery.Order("created_time DESC").Count(&total).Error; err != nil {
		return nil, err
	}

	result := baseQuery.Select("id", "notification_name", "notification_type", "status", "content", "created_time", "updated_time").
		Offset(int((req.Page - 1) * req.PageSize)).
		Limit(int(req.PageSize))

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
			status           int
			content          string
			createdTime      time.Time
			updatedTime      time.Time
		)

		if err := rows.Scan(&id, &notificationName, &notificationType, &status, &content, &createdTime, &updatedTime); err != nil {
			return nil, err
		}

		userNotifications = append(userNotifications, &v1.UserNotificationListItem{
			Id:               id,
			NotificationName: notificationName,
			NotificationType: notificationType,
			Status:           int32(status),
			Content:          content,
			CreatedTime:      createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime:      updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetUserNotificationListReply{
		Code: 200, Success: true,
		Data: &v1.UserNotificationList{
			Items: userNotifications,
			Total: total,
		},
	}, nil
}

// DeleteUserNotification 删除用户消息
func (u *userNotificationRepo) DeleteUserNotification(ctx context.Context, req *v1.DeleteUserNotificationRequest) (*v1.DeleteUserNotificationReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()
	result := tx.Model(&UserNotification{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.DeleteUserNotificationReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdateUserNotificationStatus 更新用户消息状态
func (u *userNotificationRepo) UpdateUserNotificationStatus(ctx context.Context, req *v1.UpdateUserNotificationStatusRequest) (*v1.UpdateUserNotificationStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"status": int(req.Status),
	}

	if req.Status == 1 {
		user["read_time"] = time.Now()
	} else {
		user["read_time"] = nil
	}

	result := u.data.db.Model(&UserNotification{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateUserNotificationStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// CreateUserNotification 创建用户消息
func (u *userNotificationRepo) CreateUserNotification(ctx context.Context, req *v1.CreateUserNotificationRequest) error {
	if req.NotificationName == "" || req.NotificationType == "" || req.Content == "" {
		return fmt.Errorf("notification name, notification type and content are required")
	}

	notificationId, err := u.data.gid.NextID()
	if err != nil {
		return err
	}

	userNotification := map[string]interface{}{
		"id":                int64(notificationId),
		"user_id":           req.UserId,
		"notification_name": req.NotificationName,
		"notification_type": req.NotificationType,
		"content":           req.Content,
		"status":            0,
		"created_time":      time.Now(),
		"updated_time":      time.Now(),
	}

	if err := u.data.db.Model(&UserNotification{}).Create(userNotification).Error; err != nil {
		return err
	}

	return nil
}

// BatchUserNotifications 批量操作用户消息
func (u *userNotificationRepo) BatchUserNotifications(ctx context.Context, req *v1.BatchUserNotificationsRequest) (*v1.BatchUserNotificationsReply, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items is required")
	}

	switch req.Action {
	case "read":
		userNotification := map[string]interface{}{
			"status":    1,
			"read_time": time.Now(),
		}
		if err := u.data.db.Model(&UserNotification{}).Where("id IN (?)", req.Items).Where("deleted_flag = ?", 0).Updates(userNotification).Error; err != nil {
			return nil, err
		}
	case "unread":
		userNotification := map[string]interface{}{
			"status":    0,
			"read_time": nil,
		}
		if err := u.data.db.Model(&UserNotification{}).Where("id IN (?)", req.Items).Where("deleted_flag = ?", 0).Updates(userNotification).Error; err != nil {
			return nil, err
		}
	case "delete":
		userNotification := map[string]interface{}{
			"deleted_flag": 1,
			"deleted_time": time.Now(),
		}
		if err := u.data.db.Model(&UserNotification{}).Where("id IN (?)", req.Items).Where("deleted_flag = ?", 0).Updates(userNotification).Error; err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("action is required")
	}

	return &v1.BatchUserNotificationsReply{
		Code: 200, Success: true, Message: "批量操作成功",
	}, nil
}

// UpdateUserAllNotificationStatus 更新用户消息状态
func (u *userNotificationRepo) UpdateUserAllNotificationStatus(ctx context.Context, req *v1.UpdateUserAllNotificationStatusRequest) (*v1.UpdateUserAllNotificationStatusReply, error) {
	if req.Status == 0 {
		return nil, fmt.Errorf("status is required")
	}

	userRepo := NewUserRepo(u.data, u.log.Logger())
	userId, err := userRepo.GetUserId(ctx)
	if err != nil {
		return nil, err
	}

	userNotification := map[string]interface{}{
		"status": int(req.Status),
	}

	if err := u.data.db.Model(&UserNotification{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId).Updates(userNotification).Error; err != nil {
		return nil, err
	}

	return &v1.UpdateUserAllNotificationStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}
