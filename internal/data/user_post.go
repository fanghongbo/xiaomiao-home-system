package data

import (
	"context"
	"time"
	v1 "xiaomiao-home-system/api/user/post/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userPostRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserPostRepo .
func NewUserPostRepo(data *Data, logger log.Logger) biz.UserPostRepo {
	return &userPostRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserPostRepo")),
	}
}

// GetUserPostList 查询发布内容列表
func (u *userPostRepo) GetUserPostList(ctx context.Context, req *v1.GetUserPostListRequest) (*v1.GetUserPostListReply, error) {
	var items []*v1.UserPostListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Table("t_post as t1").Joins("inner join t_user_post as t2 on t1.id = t2.post_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if req.PType > 0 {
		baseQuery = baseQuery.Where("t1.post_type = ?", req.PType)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.title", "t1.post_status", "t1.audit_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get user post list failed: %v", err)
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
			u.log.Errorf("get user post list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.UserPostListItem{
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
		u.log.Errorf("get user post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserPostListReply{
		Code: 200, Success: true,
		Data: &v1.UserPostList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CreatePost 创建发布内容
func (u *userPostRepo) CreateUserPost(ctx context.Context, req *v1.CreateUserPostRequest) (*v1.CreateUserPostReply, error) {
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

	userPostId, err := u.data.gid.NextID()
	if err != nil {
		tx.Rollback()
		u.log.Errorf("generate user post id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userPostInfo := map[string]interface{}{
		"id":           userPostId,
		"user_id":      userId,
		"post_id":      id,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := tx.Model(&UserPost{}).Create(userPostInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create user post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreateUserPostReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateUserPost 更新发布内容
func (u *userPostRepo) UpdateUserPost(ctx context.Context, req *v1.UpdateUserPostRequest) (*v1.UpdateUserPostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
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

	return &v1.UpdateUserPostReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeleteUserPost 删除发布内容
func (u *userPostRepo) DeleteUserPost(ctx context.Context, req *v1.DeleteUserPostRequest) (*v1.DeleteUserPostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	tx := u.data.db.Begin()

	result := tx.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo)

	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete post failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	result = tx.Model(&UserPost{}).Where("post_id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updateInfo)
	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete user post failed: %v", result.Error)
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

	return &v1.DeleteUserPostReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdatePostStatus 更新发布内容状态
func (u *userPostRepo) UpdateUserPostStatus(ctx context.Context, req *v1.UpdateUserPostStatusRequest) (*v1.UpdateUserPostStatusReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	post := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(post)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateUserPostStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// GetUserPost 查询发布内容
func (u *userPostRepo) GetUserPost(ctx context.Context, req *v1.GetUserPostRequest) (*v1.GetUserPostReply, error) {
	post := &Post{}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).First(post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserPostReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: &v1.UserPostInfo{
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

// CheckUserPostBelongToUser 检查当前发布内容是否属于当前用户
func (u *userPostRepo) CheckUserPostBelongToUser(ctx context.Context, userId int64, postId int64) (bool, error) {
	var count int64
	if err := u.data.db.Model(&UserPost{}).Where("user_id = ?", userId).Where("post_id = ?", postId).Where("deleted_flag = ?", 0).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
