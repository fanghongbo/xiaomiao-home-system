package data

import (
	"context"
	"time"
	v1 "xiaomiao-home-system/api/collect/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type collectRepo struct {
	data *Data
	log  *log.Helper
}

// NewCollectRepo .
func NewCollectRepo(data *Data, logger log.Logger) biz.CollectRepo {
	return &collectRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "CollectRepo")),
	}
}

// GetCollectList 查询收藏列表
func (u *collectRepo) GetCollectList(ctx context.Context, req *v1.GetCollectListRequest) (*v1.GetCollectListReply, error) {
	var items []*v1.CollectListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Model(&Publish{}).Where("deleted_flag = ?", 0).Where("user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Error("get collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("id", "title", "publish_status", "audit_status", "remark", "created_time", "updated_time").Order("created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Error("get collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
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
			u.log.Error("get collect list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.CollectListItem{
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
		u.log.Error("get collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetCollectListReply{
		Code: 200, Success: true,
		Data: &v1.CollectList{
			Items: items,
			Total: total,
		},
	}, nil
}

// GetCollectTypes 查询收藏分类
func (u *collectRepo) GetCollectTypes(ctx context.Context, req *v1.GetCollectTypesRequest) (*v1.GetCollectTypesReply, error) {
	items := make([]int64, 0)

	return &v1.GetCollectTypesReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: items,
	}, nil
}
