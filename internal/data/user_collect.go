package data

import (
	"context"
	"time"
	v1 "xiaomiao-home-system/api/user/collect/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type userCollectRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserCollectRepo .
func NewUserCollectRepo(data *Data, logger log.Logger) biz.UserCollectRepo {
	return &userCollectRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserCollectRepo")),
	}
}

// GetUserCollectList 查询用户收藏列表
func (u *userCollectRepo) GetUserCollectList(ctx context.Context, req *v1.GetUserCollectListRequest) (*v1.GetUserCollectListReply, error) {
	var items []*v1.UserCollectListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Table("t_post as t1").Joins("inner join t_user_collect as t2 on t1.id = t2.post_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.title", "t1.post_status", "t1.audit_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get user collect list failed: %v", err)
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
			u.log.Errorf("get user collect list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.UserCollectListItem{
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
		u.log.Errorf("get user collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserCollectListReply{
		Code: 200, Success: true,
		Data: &v1.UserCollectList{
			Items: items,
			Total: total,
		},
	}, nil
}

// GetCollectTypes 查询用户收藏分类
func (u *userCollectRepo) GetUserCollectTypes(ctx context.Context, req *v1.GetUserCollectTypesRequest) (*v1.GetUserCollectTypesReply, error) {
	items := make([]int64, 0)

	return &v1.GetUserCollectTypesReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: items,
	}, nil
}

// AddUserCollect 添加用户收藏
func (u *userCollectRepo) AddUserCollect(ctx context.Context, req *v1.AddUserCollectRequest) (*v1.AddUserCollectReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	id, err := u.data.gid.NextID()
	if err != nil {
		u.log.Errorf("generate id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	collectInfo := map[string]interface{}{
		"id":           id,
		"user_id":      userId,
		"post_id":      req.Id,
		"deleted_flag": 0,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := u.data.db.Model(&UserCollect{}).Create(collectInfo).Error; err != nil {
		u.log.Errorf("create user collect failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.AddUserCollectReply{
		Code: 200, Success: true, Message: "添加成功",
		Data: "添加成功",
	}, nil
}

// CancelUserCollect 取消用户收藏
func (u *userCollectRepo) CancelUserCollect(ctx context.Context, req *v1.CancelUserCollectRequest) (*v1.CancelUserCollectReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	if err := u.data.db.Model(&UserCollect{}).Where("user_id = ?", userId).Where("post_id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		u.log.Errorf("cancel user collect failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CancelUserCollectReply{
		Code: 200, Success: true, Message: "取消成功",
		Data: "取消成功",
	}, nil
}
