package data

import (
	"context"
	"time"
	v1 "xiaomiao-home-system/api/discover/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type discoverRepo struct {
	data *Data
	log  *log.Helper
}

// NewDiscoverRepo .
func NewDiscoverRepo(data *Data, logger log.Logger) biz.DiscoverRepo {
	return &discoverRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "DiscoverRepo")),
	}
}

// GetDiscoverList 查询发现列表
func (u *discoverRepo) GetDiscoverList(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.GetDiscoverListReply, error) {
	var items []*v1.DiscoverListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Model(&Post{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Error("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("id", "title", "post_status", "audit_status", "remark", "created_time", "updated_time").Order("created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Error("get discover list failed: %v", err)
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
			u.log.Error("get discover list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.DiscoverListItem{
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
		u.log.Error("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetDiscoverListReply{
		Code: 200, Success: true,
		Data: &v1.DiscoverList{
			Items: items,
			Total: total,
		},
	}, nil
}
