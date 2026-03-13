package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/role/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo .
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "RoleRepo")),
	}
}

// GetRoleList 查询角色列表
func (u *roleRepo) GetRoleList(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRoleListReply, error) {
	var Roles []*v1.RoleListItem
	var total int64

	baseQuery := u.data.db.Model(&Role{}).Where("deleted_flag = ?", 0)

	if req.RoleName != "" {
		baseQuery = baseQuery.Where("role_name LIKE ?", "%"+req.RoleName+"%")
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

	result := baseQuery.Select("id", "role_name", "status", "remark", "created_time", "updated_time").
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
			roleName    string
			status      int
			remark      string
			createdTime time.Time
			updatedTime time.Time
		)

		if err := rows.Scan(&id, &roleName, &status, &remark, &createdTime, &updatedTime); err != nil {
			return nil, err
		}

		permissions, err := u.GetRolePermissionsCodes(id)
		if err != nil {
			return nil, err
		}

		Roles = append(Roles, &v1.RoleListItem{
			Id:          id,
			RoleName:    roleName,
			Status:      int32(status),
			Remark:      remark,
			Permissions: permissions,
			CreatedTime: createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetRoleListReply{
		Code: 200, Success: true,
		Data: &v1.RoleList{
			Items: Roles,
			Total: total,
		},
	}, nil
}

// CreateRole 创建角色
func (u *roleRepo) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.CreateRoleReply, error) {
	id, err := u.data.gid.NextID()
	if err != nil {
		return nil, err
	}

	if req.RoleName == "" {
		return nil, fmt.Errorf("group name is required")
	}

	role := map[string]interface{}{
		"id":           int64(id),
		"role_name":    req.RoleName,
		"status":       int(req.Status),
		"remark":       req.Remark,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Role{}).Create(role).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("角色 %s 已经存在", req.RoleName)
		}
		return nil, err
	}

	// 创建角色权限
	if err := u.CreateRolePermissions(tx, int64(id), req.Permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.CreateRoleReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateRole 更新角色
func (u *roleRepo) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.UpdateRoleReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	if req.RoleName == "" {
		return nil, fmt.Errorf("name is required")
	}

	updates := map[string]interface{}{
		"role_name": req.RoleName,
		"status":    int(req.Status),
		"remark":    req.Remark,
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Role{}).
		Where("id = ?", req.Id).
		Where("deleted_flag = ?", 0).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		if utils.IsDuplicateEntryError(err) {
			return nil, fmt.Errorf("角色名称 %s 已经存在", req.RoleName)
		}
		return nil, err
	}

	// 更新角色权限
	if err := u.UpdateRolePermissions(tx, req.Id, req.Permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.UpdateRoleReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeleteRole 删除角色
func (u *roleRepo) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.DeleteRoleReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	role := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()

	result := tx.Model(&Role{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(role)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// 删除角色权限
	if err := u.DeleteRolePermissions(tx, req.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &v1.DeleteRoleReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdateRoleStatus 更新角色状态
func (u *roleRepo) UpdateRoleStatus(ctx context.Context, req *v1.UpdateRoleStatusRequest) (*v1.UpdateRoleStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	role := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&Role{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(role)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateRoleStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// CreateRolePermissions 创建角色权限
func (u *roleRepo) CreateRolePermissions(tx *gorm.DB, roleId int64, permissions []string) error {
	if roleId == 0 {
		return fmt.Errorf("role id is required")
	}

	rolePermissions := make([]map[string]interface{}, 0)

	// 对permissions 去重复
	permissions = utils.RemoveDuplicate(permissions)

	for _, permission := range permissions {
		relationId, err := u.data.gid.NextID()
		if err != nil {
			return err
		}

		if permission == "" {
			continue
		}

		rolePermission := map[string]interface{}{
			"id":              int64(relationId),
			"role_id":         roleId,
			"permission_code": permission,
			"created_time":    time.Now(),
			"updated_time":    time.Now(),
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	if len(rolePermissions) == 0 {
		return nil
	}

	return tx.Model(&RolePermission{}).Create(rolePermissions).Error
}

// UpdateRolePermissions 更新角色权限
func (u *roleRepo) UpdateRolePermissions(tx *gorm.DB, roleId int64, permissions []string) error {
	if roleId == 0 {
		return fmt.Errorf("role id is required")
	}

	// 删除角色权限
	if err := u.DeleteRolePermissions(tx, roleId); err != nil {
		return err
	}

	rolePermissions := make([]map[string]interface{}, 0)

	// 对permissions 去重复
	permissions = utils.RemoveDuplicate(permissions)

	for _, permission := range permissions {
		relationId, err := u.data.gid.NextID()
		if err != nil {
			return err
		}

		if permission == "" {
			continue
		}

		rolePermission := map[string]interface{}{
			"id":              int64(relationId),
			"role_id":         roleId,
			"permission_code": permission,
			"created_time":    time.Now(),
			"updated_time":    time.Now(),
		}

		rolePermissions = append(rolePermissions, rolePermission)
	}

	if len(rolePermissions) == 0 {
		return nil
	}

	return tx.Model(&RolePermission{}).Create(rolePermissions).Error
}

// DeleteRolePermissions 删除角色权限
func (u *roleRepo) DeleteRolePermissions(tx *gorm.DB, roleId int64) error {
	if roleId == 0 {
		return fmt.Errorf("role id is required")
	}

	roleUpdates := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	return tx.Model(&RolePermission{}).Where("role_id = ?", roleId).Where("deleted_flag = ?", 0).Updates(roleUpdates).Error
}

// GetRolePermissionsCodes 获取角色权限代码列表
func (u *roleRepo) GetRolePermissionsCodes(roleId int64) ([]string, error) {
	var permissions []string

	if err := u.data.db.Model(&RolePermission{}).Where("role_id = ?", roleId).Where("deleted_flag = ?", 0).Pluck("permission_code", &permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetRoles 查询所有角色
func (u *roleRepo) GetRoles(ctx context.Context, req *v1.GetRolesRequest) (*v1.GetRolesReply, error) {
	var roles []*v1.RoleItem

	result := u.data.db.Model(&Role{}).
		Where("deleted_flag = ?", 0).
		Select("id", "role_name")

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id       int64
			roleName string
		)

		if err := rows.Scan(&id, &roleName); err != nil {
			return nil, err
		}

		roles = append(roles, &v1.RoleItem{
			Id:   id,
			Name: roleName,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetRolesReply{
		Code: 200, Success: true,
		Data: roles,
	}, nil
}
