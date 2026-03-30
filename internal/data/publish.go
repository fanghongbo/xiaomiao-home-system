package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/publish/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type publishRepo struct {
	data *Data
	log  *log.Helper
}

// NewPublishRepo .
func NewPublishRepo(data *Data, logger log.Logger) biz.PublishRepo {
	return &publishRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "PublishRepo")),
	}
}

// GetPublishList 查询发布内容列表
func (u *publishRepo) GetPublishList(ctx context.Context, req *v1.GetPublishListRequest) (*v1.GetPublishListReply, error) {
	var items []*v1.PublishListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Model(&Publish{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId)

	if req.PType > 0 {
		baseQuery = baseQuery.Where("publish_type = ?", req.PType)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	result := baseQuery.Select("id", "title", "publish_status", "audit_status", "remark", "created_time", "updated_time").Order("created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            int64
			title         string
			publishStatus int
			auditStatus   int
			remark        string
			createdTime   time.Time
			updatedTime   time.Time
		)

		if err := rows.Scan(&id, &title, &publishStatus, &auditStatus, &remark, &createdTime, &updatedTime); err != nil {
			return nil, err
		}

		items = append(items, &v1.PublishListItem{
			Id:            id,
			Title:         title,
			PublishStatus: int32(publishStatus),
			AuditStatus:   int32(auditStatus),
			Remark:        remark,
			CreatedTime:   createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime:   updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &v1.GetPublishListReply{
		Code: 200, Success: true,
		Data: &v1.PublishList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CreatePublish 创建发布内容
func (u *publishRepo) CreatePublish(ctx context.Context, req *v1.CreatePublishRequest) (*v1.CreatePublishReply, error) {
	id, err := u.data.gid.NextID()
	if err != nil {
		u.log.Error("generate id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	publish := map[string]interface{}{
		"id":             id,
		"user_id":        userId,
		"title":          req.Title,
		"publish_type":   req.PublishType,
		"province_id":    req.ProvinceId,
		"city_id":        req.CityId,
		"address":        req.Address,
		"audit_status":   0,
		"publish_status": 0,
		"remark":         req.Remark,
		"created_time":   time.Now(),
		"updated_time":   time.Now(),
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Publish{}).Create(publish).Error; err != nil {
		tx.Rollback()
		u.log.Error("create publish failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("create publish failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreatePublishReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdatePublish 更新发布内容
func (u *publishRepo) UpdatePublish(ctx context.Context, req *v1.UpdatePublishRequest) (*v1.UpdatePublishReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	publish := map[string]interface{}{
		"id":             req.Id,
		"title":          req.Title,
		"publish_type":   req.PublishType,
		"province_id":    req.ProvinceId,
		"city_id":        req.CityId,
		"address":        req.Address,
		"audit_status":   0,
		"publish_status": 0,
		"remark":         req.Remark,
		"updated_time":   time.Now(),
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	result := tx.Model(&Publish{}).
		Where("id = ?", req.Id).
		Where("user_id = ?", userId).
		Where("deleted_flag = ?", 0).
		Updates(publish)
	if result.Error != nil {
		tx.Rollback()
		u.log.Error("update publish failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("update publish failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdatePublishReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeletePublish 删除发布内容
func (u *publishRepo) DeletePublish(ctx context.Context, req *v1.DeletePublishRequest) (*v1.DeletePublishReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	publish := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()

	result := tx.Model(&Publish{}).Where("id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(publish)

	if result.Error != nil {
		tx.Rollback()
		u.log.Error("delete publish failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("delete publish failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.DeletePublishReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdatePublishStatus 更新发布内容状态
func (u *publishRepo) UpdatePublishStatus(ctx context.Context, req *v1.UpdatePublishStatusRequest) (*v1.UpdatePublishStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	publish := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&Publish{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(publish)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdatePublishStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// GetPublish 查询发布内容
func (u *publishRepo) GetPublish(ctx context.Context, req *v1.GetPublishRequest) (*v1.GetPublishReply, error) {
	publish := &Publish{}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.data.db.Model(&Publish{}).Where("id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).First(publish).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetPublishReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: &v1.PublishInfo{
			Id:          publish.Id,
			Title:       publish.Title,
			PublishType: int32(publish.PublishType),
			ProvinceId:  int32(publish.ProvinceId),
			CityId:      int32(publish.CityId),
			Address:     publish.Address,
			CatType:     1,
			CatBreed:    1,
			CatGender:   1,
			Remark:      publish.Remark,
			CreatedTime: publish.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: publish.UpdatedTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
