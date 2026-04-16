package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/user/like/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type userLikeRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserLikeRepo .
func NewUserLikeRepo(data *Data, logger log.Logger) biz.UserLikeRepo {
	return &userLikeRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserLikeRepo")),
	}
}

// AddUserLike 添加用户喜欢
func (u *userLikeRepo) AddUserLike(ctx context.Context, req *v1.AddUserLikeRequest) (*v1.AddUserLikeReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 检查用户添加喜欢次数限制
	if err := u.checkUserAddUserLikeCountLimit(ctx, userId); err != nil {
		return nil, err
	}

	id, err := u.data.gid.NextID()
	if err != nil {
		u.log.Errorf("generate id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	likeInfo := map[string]interface{}{
		"id":           id,
		"user_id":      userId,
		"post_id":      req.Id,
		"deleted_flag": 0,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := u.data.db.Model(&UserLike{}).Create(likeInfo).Error; err != nil {
		if utils.IsDuplicateEntryError(err) {
			return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "已喜欢")
		}
		u.log.Errorf("create user like failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserPostLikeStatusCache(ctx, userId, req.Id); err != nil {
		u.log.Errorf("remove user post like status cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.AddUserLikeReply{
		Code: 200, Success: true, Message: "添加成功",
		Data: "添加成功",
	}, nil
}

// CancelUserLike 取消用户喜欢
func (u *userLikeRepo) CancelUserLike(ctx context.Context, req *v1.CancelUserLikeRequest) (*v1.CancelUserLikeReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	if err := u.data.db.Model(&UserLike{}).Where("user_id = ?", userId).Where("post_id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		if utils.IsDuplicateEntryError(err) {
			return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "已取消喜欢")
		}
		u.log.Errorf("cancel user like failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.removeUserPostLikeStatusCache(ctx, userId, req.Id); err != nil {
		u.log.Errorf("remove user post like status cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CancelUserLikeReply{
		Code: 200, Success: true, Message: "取消成功",
		Data: "取消成功",
	}, nil
}

// GetUserPostLikeStatus 查询用户发布内容喜欢状态
func (u *userLikeRepo) GetUserPostLikeStatus(ctx context.Context, postId int64) (bool, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return false, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	var count int64

	count, err = u.getUserPostLikeStatusCache(ctx, userId, postId)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			u.log.Errorf("get user post like status cache failed: %v", err)
			return false, err
		}
	} else {
		return count > 0, nil
	}

	if err := u.data.db.Model(&UserLike{}).Where("user_id = ?", userId).Where("post_id = ?", postId).Where("deleted_flag = ?", 0).Count(&count).Error; err != nil {
		u.log.Errorf("check user post like status failed: %v", err)
		return false, err
	}

	if err := u.setUserPostLikeStatusCache(ctx, userId, postId, count); err != nil {
		u.log.Errorf("set user post like status cache failed: %v", err)
		return false, err
	}

	return count > 0, nil
}

// getUserPostLikeStatusCacheKey 获取查询缓存key
func (u *userLikeRepo) getUserPostLikeStatusCacheKey(userId int64, postId int64) string {
	return fmt.Sprintf("user:post:like:status:%d:%d", userId, postId)
}

// getUserPostLikeStatusCache 获取用户发布内容喜欢状态缓存
func (u *userLikeRepo) getUserPostLikeStatusCache(ctx context.Context, userId int64, postId int64) (int64, error) {
	redisKey := u.getUserPostLikeStatusCacheKey(userId, postId)

	count, err := u.data.cache.Get(ctx, redisKey).Int64()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// setUserPostLikeStatusCache 设置用户发布内容喜欢状态缓存
func (u *userLikeRepo) setUserPostLikeStatusCache(ctx context.Context, userId int64, postId int64, count int64) error {
	redisKey := u.getUserPostLikeStatusCacheKey(userId, postId)

	ttl := time.Minute * 1
	if err := u.data.cache.Set(ctx, redisKey, count, ttl).Err(); err != nil {
		return err
	}

	return nil
}

// removeUserPostLikeStatusCache 删除用户发布内容喜欢状态缓存
func (u *userLikeRepo) removeUserPostLikeStatusCache(ctx context.Context, userId int64, postId int64) error {
	redisKey := u.getUserPostLikeStatusCacheKey(userId, postId)

	if err := u.data.cache.Del(ctx, redisKey).Err(); err != nil && err != redis.Nil {
		return err
	}

	return nil
}

// checkUserAddUserLikeCountLimit 检查用户添加喜欢次数限制
func (u *userLikeRepo) checkUserAddUserLikeCountLimit(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("user:like:add:count:%d", userId)
	n, err := u.data.cache.Incr(ctx, key).Result()
	if err != nil {
		u.log.Errorf("increase user like add count failed, key=%s: %v", key, err)
		return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if n == 1 {
		if err := u.data.cache.Expire(ctx, key, 8*time.Hour).Err(); err != nil {
			u.log.Errorf("set user like add count ttl failed, key=%s: %v", key, err)
			return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if n > 100 {
		return errors.BadRequest(v1.ErrorReason_ERR_TOO_MANY_REQUEST.String(), "今日已达最大添加喜欢次数, 请明日再试")
	}

	return nil
}
