package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/user/group/v1"
	"xiaomiao-home-system/utils"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userGroupRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserGroupRepo .
func NewUserGroupRepo(data *Data, logger log.Logger) biz.UserGroupRepo {
	return &userGroupRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserGroupRepo")),
	}
}

// GetUserGroupList 查询用户组列表
func (u *userGroupRepo) GetUserGroupList(ctx context.Context, req *v1.GetUserGroupListRequest) (*v1.GetUserGroupListReply, error) {
	var userGroups []*v1.UserGroupListItem
	var total int64

	baseQuery := u.data.db.Model(&UserGroup{}).Where("deleted_flag = ?", 0)

	if req.GroupName != "" {
		baseQuery = baseQuery.Where("group_name LIKE ?", "%"+req.GroupName+"%")
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

	result := baseQuery.Select("id", "group_name", "status", "remark", "created_time", "updated_time").
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
			groupName   string
			status      int
			remark      string
			createdTime time.Time
			updatedTime time.Time
		)

		if err := rows.Scan(&id, &groupName, &status, &remark, &createdTime, &updatedTime); err != nil {
			return nil, err
		}

		users, err := u.GetUserGroupUsersRelation(id)
		if err != nil {
			return nil, err
		}

		userGroups = append(userGroups, &v1.UserGroupListItem{
			Id:          id,
			GroupName:   groupName,
			Status:      int32(status),
			Remark:      remark,
			Users:       users,
			UserCount:   int64(len(users)),
			CreatedTime: createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetUserGroupListReply{
		Code: 200, Success: true,
		Data: &v1.UserGroupList{
			Items: userGroups,
			Total: total,
		},
	}, nil
}

// CreateUserGroup 创建用户组
func (u *userGroupRepo) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (*v1.CreateUserGroupReply, error) {
	id, err := u.data.gid.NextID()
	if err != nil {
		return nil, err
	}

	if req.GroupName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	userGroup := map[string]interface{}{
		"id":           int64(id),
		"group_name":   req.GroupName,
		"status":       int(req.Status),
		"remark":       req.Remark,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&UserGroup{}).Create(userGroup).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("用户组 %s 已经存在", req.GroupName)
		}
		return nil, err
	}

	if err := u.CreateUserGroupUsersRelation(tx, int64(id), req.Users); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateUserGroupReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateUserGroup 更新用户组
func (u *userGroupRepo) UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (*v1.UpdateUserGroupReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	if req.GroupName == "" {
		return nil, fmt.Errorf("name is required")
	}

	updates := map[string]interface{}{
		"group_name": req.GroupName,
		"status":     int(req.Status),
		"remark":     req.Remark,
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&UserGroup{}).
		Where("id = ?", req.Id).
		Where("deleted_flag = ?", 0).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("用户组名称 %s 已经存在", req.GroupName)
		}
		return nil, err
	}

	if err := u.UpdateUserGroupUsersRelation(tx, req.Id, req.Users); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.UpdateUserGroupReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeleteUserGroup 删除用户组
func (u *userGroupRepo) DeleteUserGroup(ctx context.Context, req *v1.DeleteUserGroupRequest) (*v1.DeleteUserGroupReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()
	result := tx.Model(&UserGroup{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := u.DeleteUserGroupUsersRelation(tx, req.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.DeleteUserGroupReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdateUserGroupStatus 更新用户组状态
func (u *userGroupRepo) UpdateUserGroupStatus(ctx context.Context, req *v1.UpdateUserGroupStatusRequest) (*v1.UpdateUserGroupStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	user := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&UserGroup{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateUserGroupStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// GetUserGroups 查询所有用户组
func (u *userGroupRepo) GetUserGroups(ctx context.Context, req *v1.GetUserGroupsRequest) (*v1.GetUserGroupsReply, error) {
	var groups []*v1.UserGroupItem

	result := u.data.db.Model(&UserGroup{}).Where("deleted_flag = ?", 0).Select("id", "group_name")

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int64
			groupName string
		)

		if err := rows.Scan(&id, &groupName); err != nil {
			return nil, err
		}

		groups = append(groups, &v1.UserGroupItem{
			Id:   id,
			Name: groupName,
		})
	}

	return &v1.GetUserGroupsReply{
		Code: 200, Success: true,
		Data: groups,
	}, nil
}

// CreateUserGroupUsersRelation 创建用户组用户关系
func (u *userGroupRepo) CreateUserGroupUsersRelation(tx *gorm.DB, groupId int64, users []int64) error {
	if groupId == 0 {
		return fmt.Errorf("group id is required")
	}

	if len(users) == 0 {
		return nil
	}

	userGroupUsers := make([]map[string]interface{}, 0)

	for _, user := range users {
		relationId, err := u.data.gid.NextID()
		if err != nil {
			return err
		}

		userGroupUsers = append(userGroupUsers, map[string]interface{}{
			"id":           int64(relationId),
			"group_id":     groupId,
			"user_id":      user,
			"created_time": time.Now(),
			"updated_time": time.Now(),
		})
	}

	if len(userGroupUsers) == 0 {
		return nil
	}

	return tx.Model(&UserGroupUsers{}).Create(userGroupUsers).Error
}

// DeleteUserGroupUsersRelation 删除用户组用户关系
func (u *userGroupRepo) DeleteUserGroupUsersRelation(tx *gorm.DB, groupId int64) error {
	if groupId == 0 {
		return fmt.Errorf("group id is required")
	}

	return tx.Model(&UserGroupUsers{}).Where("group_id = ?", groupId).Where("deleted_flag = ?", 0).Updates(map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}).Error
}

// UpdateUserGroupUsersRelation 更新用户组用户关系
func (u *userGroupRepo) UpdateUserGroupUsersRelation(tx *gorm.DB, groupId int64, users []int64) error {
	if groupId == 0 {
		return fmt.Errorf("group id is required")
	}

	if err := u.DeleteUserGroupUsersRelation(tx, groupId); err != nil {
		return err
	}

	if err := u.CreateUserGroupUsersRelation(tx, groupId, users); err != nil {
		return err
	}

	return nil
}

// GetUserGroupUsersRelation 获取用户组用户关系
func (u *userGroupRepo) GetUserGroupUsersRelation(groupId int64) ([]int64, error) {
	if groupId == 0 {
		return nil, fmt.Errorf("group id is required")
	}

	var users []int64

	if err := u.data.db.Model(&UserGroupUsers{}).Where("group_id = ?", groupId).Where("deleted_flag = ?", 0).Pluck("user_id", &users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
