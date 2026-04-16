package data

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
	v1 "xiaomiao-home-system/api/user/collect/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
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

	baseQuery := u.data.db.Table("t_user_collect as t1").Joins("inner join t_post as t2 on t1.post_id = t2.id").Joins("inner join t_post_cat as t3 on t2.id = t3.post_id").Joins("inner join t_cat as t4 on t3.cat_id = t4.id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t3.deleted_flag = ?", 0).Where("t4.deleted_flag = ?", 0).Where("t1.user_id = ?", userId)

	if req.CType > 0 {
		baseQuery = baseQuery.Where("t4.cat_type = ?", req.CType)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user collect list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t2.id", "t2.title", "t2.post_status", "t2.audit_status", "t2.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
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
			Id:            id,
			Title:         title,
			PostStatus:    int32(postStatus),
			AuditStatus:   int32(auditStatus),
			CollectStatus: 1,
			CoverImage:    "",
			Remark:        remark,
			CreatedTime:   createdTime.Format("2006-01-02 15:04:05"),
			UpdatedTime:   updatedTime.Format("2006-01-02 15:04:05"),
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
	items := make([]int32, 0)
	res := &v1.GetUserCollectTypesReply{Code: 200, Success: true, Message: "查询成功"}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	types, err := u.getUserCollectTypesCache(ctx, userId)
	if err != nil  {
		if !errors.Is(err, redis.Nil) {
			u.log.Errorf("get user collect types cache failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	} else {
		res.Data = types
		return res, nil
	}

	baseQuery := u.data.db.Table("t_user_collect as t1").Joins("inner join t_post_cat as t2 on t1.post_id = t2.post_id").Joins("inner join t_cat as t3 on t2.cat_id = t3.id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t3.deleted_flag = ?", 0).Where("t1.user_id = ?", userId)

	if err := baseQuery.Distinct("t3.cat_type").Find(&items).Error; err != nil {
		u.log.Errorf("get user collect types failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})

	if err := u.setUserCollectTypesCache(ctx, userId, items); err != nil {
		u.log.Errorf("set user collect types cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	res.Data = items

	return res, nil
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
		if utils.IsDuplicateEntryError(err) {
			return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "已收藏")
		}
		u.log.Errorf("create user collect failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserPostCollectStatusCache(ctx, userId, req.Id); err != nil {
		u.log.Errorf("remove user post collect status cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserCollectTypesCache(ctx, userId); err != nil {
		u.log.Errorf("remove user collect types cache failed: %v", err)
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
		if utils.IsDuplicateEntryError(err) {
			return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "已取消收藏")
		}
		u.log.Errorf("cancel user collect failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserPostCollectStatusCache(ctx, userId, req.Id); err != nil {
		u.log.Errorf("remove user post collect status cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserCollectTypesCache(ctx, userId); err != nil {
		u.log.Errorf("remove user collect types cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CancelUserCollectReply{
		Code: 200, Success: true, Message: "取消成功",
		Data: "取消成功",
	}, nil
}

// GetUserPostCollectStatus 查询用户发布内容收藏状态
func (u *userCollectRepo) GetUserPostCollectStatus(ctx context.Context, postId int64) (bool, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return false, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	var count int64

	count, err = u.getUserPostCollectStatusCache(ctx, userId, postId)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			u.log.Errorf("get user post collect status cache failed: %v", err)
			return false, err
		}
	} else {
		return count > 0, nil
	}

	if err := u.data.db.Model(&UserCollect{}).Where("user_id = ?", userId).Where("post_id = ?", postId).Where("deleted_flag = ?", 0).Count(&count).Error; err != nil {
		u.log.Errorf("check user post collect status failed: %v", err)
		return false, err
	}

	if err := u.setUserPostCollectStatusCache(ctx, userId, postId, count); err != nil {
		u.log.Errorf("set user post collect status cache failed: %v", err)
		return false, err
	}

	return count > 0, nil
}

// getUserPostCollectStatusCacheKey 获取查询缓存key
func (u *userCollectRepo) getUserPostCollectStatusCacheKey(userId int64, postId int64) string {
	return fmt.Sprintf("user:post:collect:status:%d:%d", userId, postId)
}

// getUserPostCollectStatusCache 获取用户发布内容收藏状态缓存
func (u *userCollectRepo) getUserPostCollectStatusCache(ctx context.Context, userId int64, postId int64) (int64, error) {
	redisKey := u.getUserPostCollectStatusCacheKey(userId, postId)

	count, err := u.data.cache.Get(ctx, redisKey).Int64()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// setUserPostCollectStatusCache 设置用户发布内容收藏状态缓存
func (u *userCollectRepo) setUserPostCollectStatusCache(ctx context.Context, userId int64, postId int64, count int64) error {
	redisKey := u.getUserPostCollectStatusCacheKey(userId, postId)

	ttl := time.Minute * 1
	if err := u.data.cache.Set(ctx, redisKey, count, ttl).Err(); err != nil {
		return err
	}

	return nil
}

// removeUserPostCollectStatusCache 删除用户发布内容收藏状态缓存
func (u *userCollectRepo) removeUserPostCollectStatusCache(ctx context.Context, userId int64, postId int64) error {
	redisKey := u.getUserPostCollectStatusCacheKey(userId, postId)

	if err := u.data.cache.Del(ctx, redisKey).Err(); err != nil && err != redis.Nil {
		return err
	}

	return nil
}

// getUserCollectTypesCacheKey 获取用户收藏分类缓存key
func (u *userCollectRepo) getUserCollectTypesCacheKey(userId int64) string {
	return fmt.Sprintf("user:collect:types:%d", userId)
}

// getUserCollectTypesCache 获取用户收藏分类缓存
func (u *userCollectRepo) getUserCollectTypesCache(ctx context.Context, userId int64) ([]int32, error) {
	redisKey := u.getUserCollectTypesCacheKey(userId)

	types, err := u.data.cache.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	var items []int32
	if err := json.Unmarshal([]byte(types), &items); err != nil {
		return nil, err
	}

	return items, nil
}

// setUserCollectTypesCache 设置用户收藏分类缓存
func (u *userCollectRepo) setUserCollectTypesCache(ctx context.Context, userId int64, types []int32) error {
	redisKey := u.getUserCollectTypesCacheKey(userId)

	typesJSON, err := json.Marshal(types)
	if err != nil {
		return err
	}

	ttl := time.Second * 10
	if err := u.data.cache.Set(ctx, redisKey, typesJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

// removeUserCollectTypesCache 删除用户收藏分类缓存
func (u *userCollectRepo) removeUserCollectTypesCache(ctx context.Context, userId int64) error {
	redisKey := u.getUserCollectTypesCacheKey(userId)

	if err := u.data.cache.Del(ctx, redisKey).Err(); err != nil && err != redis.Nil {
		return err
	}

	return nil
}
