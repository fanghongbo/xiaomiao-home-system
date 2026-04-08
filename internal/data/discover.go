package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/discover/v1"
	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
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

// GetQueryCache 获取查询缓存
func (u *discoverRepo) GetQueryCache(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.DiscoverList, error) {
	redisKey := fmt.Sprintf("discover:list:%d:%d:%d:%d:%d:%d:%d", req.ProvinceId, req.CityId, req.PType, req.CBreed, req.CType, req.Page, req.Size)

	jsonData, err := u.data.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data v1.DiscoverList
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// SetQueryCache 设置查询缓存
func (u *discoverRepo) SetQueryCache(ctx context.Context, req *v1.GetDiscoverListRequest, data *v1.DiscoverList) error {
	redisKey := fmt.Sprintf("discover:list:%d:%d:%d:%d:%d:%d:%d", req.ProvinceId, req.CityId, req.PType, req.CBreed, req.CType, req.Page, req.Size)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := time.Minute * 1
	if err := u.data.rdb.Set(ctx, redisKey, jsonData, ttl).Err(); err != nil {
		return err
	}

	return nil
}

// GetDiscoverList 查询发现列表
func (u *discoverRepo) GetDiscoverList(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.GetDiscoverListReply, error) {
	var items []*v1.DiscoverListItem
	var total int64

	res := &v1.GetDiscoverListReply{
		Code: 200, Success: true,
		Data: &v1.DiscoverList{
			Items: []*v1.DiscoverListItem{},
			Total: 0,
		},
	}

	// 登录后查看更多内容
	if req.Page > 3 || req.Size > 30 {
		return res, nil
	}

	cacheData, err := u.GetQueryCache(ctx, req)
	if err != nil {
		u.log.Errorf("get query cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if cacheData != nil {
		res.Data = cacheData
		return res, nil
	}

	baseQuery := u.data.db.Model(&Post{}).Where("deleted_flag = ?", 0)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("id", "title", "post_status", "audit_status", "remark", "created_time", "updated_time").Order("created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get discover list failed: %v", err)
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
			u.log.Errorf("get discover list failed: %v", err)
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
		u.log.Errorf("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	res.Data.Items = items
	res.Data.Total = total

	if err := u.SetQueryCache(ctx, req, res.Data); err != nil {
		u.log.Errorf("set query cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return res, nil
}
