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
