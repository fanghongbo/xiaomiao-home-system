package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/post/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type postRepo struct {
	data *Data
	log  *log.Helper
}

// NewPostRepo .
func NewPostRepo(data *Data, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "PostRepo")),
	}
}

// GetPostList 查询发布内容列表
func (u *postRepo) GetPostList(ctx context.Context, req *v1.GetPostListRequest) (*v1.GetPostListReply, error) {
	var items []*v1.PostListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Model(&Post{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId)

	if req.PType > 0 {
		baseQuery = baseQuery.Where("post_type = ?", req.PType)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("id", "title", "post_status", "audit_status", "remark", "created_time", "updated_time").Order("created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id          int64
			title       string
			postStatus  int
			auditStatus int
			remark      string
			createdTime time.Time
			updatedTime time.Time
		)

		if err := rows.Scan(&id, &title, &postStatus, &auditStatus, &remark, &createdTime, &updatedTime); err != nil {
			u.log.Errorf("get post list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.PostListItem{
			Id:          id,
			Title:       title,
			PostStatus:  int32(postStatus),
			AuditStatus: int32(auditStatus),
			Remark:      remark,
			CreatedTime: createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		u.log.Errorf("get post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetPostListReply{
		Code: 200, Success: true,
		Data: &v1.PostList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CreatePost 创建发布内容
func (u *postRepo) CreatePost(ctx context.Context, req *v1.CreatePostRequest) (*v1.CreatePostReply, error) {
	id, err := u.data.gid.NextID()
	if err != nil {
		u.log.Errorf("generate id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	post := map[string]interface{}{
		"id":           id,
		"user_id":      userId,
		"title":        req.Title,
		"post_type":    req.PostType,
		"province_id":  req.ProvinceId,
		"city_id":      req.CityId,
		"address":      req.Address,
		"audit_status": 0,
		"post_status":  0,
		"remark":       req.Remark,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Post{}).Create(post).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreatePostReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdatePost 更新发布内容
func (u *postRepo) UpdatePost(ctx context.Context, req *v1.UpdatePostRequest) (*v1.UpdatePostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	post := map[string]interface{}{
		"id":           req.Id,
		"title":        req.Title,
		"post_type":    req.PostType,
		"province_id":  req.ProvinceId,
		"city_id":      req.CityId,
		"address":      req.Address,
		"audit_status": 0,
		"post_status":  0,
		"remark":       req.Remark,
		"updated_time": time.Now(),
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	result := tx.Model(&Post{}).
		Where("id = ?", req.Id).
		Where("user_id = ?", userId).
		Where("deleted_flag = ?", 0).
		Updates(post)
	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("update post failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("update post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdatePostReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeletePost 删除发布内容
func (u *postRepo) DeletePost(ctx context.Context, req *v1.DeletePostRequest) (*v1.DeletePostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	post := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()

	result := tx.Model(&Post{}).Where("id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(post)

	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete post failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("delete post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.DeletePostReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdatePostStatus 更新发布内容状态
func (u *postRepo) UpdatePostStatus(ctx context.Context, req *v1.UpdatePostStatusRequest) (*v1.UpdatePostStatusReply, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	post := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(post)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdatePostStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// GetPost 查询发布内容
func (u *postRepo) GetPost(ctx context.Context, req *v1.GetPostRequest) (*v1.GetPostReply, error) {
	post := &Post{}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).First(post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetPostReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: &v1.PostInfo{
			Id:          post.Id,
			Title:       post.Title,
			PostType:    int32(post.PostType),
			ProvinceId:  int32(post.ProvinceId),
			CityId:      int32(post.CityId),
			Address:     post.Address,
			CatType:     1,
			CatBreed:    1,
			CatGender:   1,
			Remark:      post.Remark,
			CreatedTime: post.CreatedTime.Format("2006-01-02 15:04:05"),
			UpdatedTime: post.UpdatedTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
