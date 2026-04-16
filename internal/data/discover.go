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

// getQueryCacheKey 获取查询缓存key
func (u *discoverRepo) getQueryCacheKey(req *v1.GetDiscoverListRequest) string {
	return fmt.Sprintf("discover:list:%d:%d:%d:%d:%d:%d", req.ProvinceId, req.CityId, req.PType, req.CBreed, req.Page, req.Size)
}

// getQueryCache 获取查询缓存
func (u *discoverRepo) getQueryCache(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.DiscoverList, error) {
	redisKey := u.getQueryCacheKey(req)

	jsonData, err := u.data.cache.Get(ctx, redisKey).Result()
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

// setQueryCache 设置查询缓存
func (u *discoverRepo) setQueryCache(ctx context.Context, req *v1.GetDiscoverListRequest, data *v1.DiscoverList) error {
	redisKey := u.getQueryCacheKey(req)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := time.Minute * 1
	if err := u.data.cache.Set(ctx, redisKey, jsonData, ttl).Err(); err != nil {
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

	cacheData, err := u.getQueryCache(ctx, req)
	if err != nil {
		u.log.Errorf("get query cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if cacheData != nil {
		res.Data = cacheData
		return res, nil
	}

	baseQuery := u.data.db.Table("t_post as t1").Joins("inner join t_post_cat as t2 on t1.id = t2.post_id").Joins("inner join t_cat as t3 on t2.cat_id = t3.id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t3.deleted_flag = ?", 0).Where("t1.audit_status = ?", 1)

	if req.PType > 0 {
		baseQuery = baseQuery.Where("t1.post_type = ?", req.PType)
	}

	if req.CBreed > 0 {
		baseQuery = baseQuery.Where("t3.breed_type = ?", req.CBreed)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.title", "t1.post_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
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
			id            int64
			title         string
			postStatus    int
			remark        string
			createdTime   time.Time
			updatedTime   time.Time
		)

		if err := rows.Scan(&id, &title, &postStatus, &remark, &createdTime, &updatedTime); err != nil {
			u.log.Errorf("get discover list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.DiscoverListItem{
			Id:            id,
			Title:         title,
			PostStatus:    int32(postStatus),
			Remark:        remark,
			CreatedTime:   createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime:   updatedTime.Format("2006-01-02 15:04:05"),
		})
	}

	if err := rows.Err(); err != nil {
		u.log.Errorf("get discover list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	res.Data.Items = items
	res.Data.Total = total

	if err := u.setQueryCache(ctx, req, res.Data); err != nil {
		u.log.Errorf("set query cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return res, nil
}
