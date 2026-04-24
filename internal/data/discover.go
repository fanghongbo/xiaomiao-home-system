package data

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	v1 "xiaomiao-home-system/api/discover/v1"
	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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

// getDiscoverListCacheKey 获取发现列表缓存key
func (u *discoverRepo) getDiscoverListCacheKey(req *v1.GetDiscoverListRequest) string {
	return fmt.Sprintf("discover:list:%d:%d:%d:%d:%d:%d", req.ProvinceId, req.CityId, req.PType, req.CBreed, req.Page, req.Size)
}

// getDiscoverListCache 获取发现列表缓存
func (u *discoverRepo) getDiscoverListCache(ctx context.Context, req *v1.GetDiscoverListRequest) (*v1.DiscoverList, error) {
	redisKey := u.getDiscoverListCacheKey(req)

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

// setDiscoverListCache 设置发现列表缓存
func (u *discoverRepo) setDiscoverListCache(ctx context.Context, req *v1.GetDiscoverListRequest, data *v1.DiscoverList) error {
	redisKey := u.getDiscoverListCacheKey(req)

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

	cacheData, err := u.getDiscoverListCache(ctx, req)
	if err != nil {
		u.log.Errorf("get discover list cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if cacheData != nil {
		res.Data = cacheData
		return res, nil
	}

	baseQuery := u.data.db.Table("t_post as t1").Joins("inner join t_post_cat as t2 on t1.id = t2.post_id").Joins("inner join t_cat as t3 on t2.cat_id = t3.id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t3.deleted_flag = ?", 0)

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
			id          int64
			title       string
			postStatus  int
			remark      string
			createdTime time.Time
			updatedTime time.Time
		)

		if err := rows.Scan(&id, &title, &postStatus, &remark, &createdTime, &updatedTime); err != nil {
			u.log.Errorf("get discover list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.DiscoverListItem{
			Id:          id,
			Title:       title,
			PostStatus:  int32(postStatus),
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

	if err := u.setDiscoverListCache(ctx, req, res.Data); err != nil {
		u.log.Errorf("set discover list cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return res, nil
}

// GetDiscover 查询发现内容
func (u *discoverRepo) GetDiscover(ctx context.Context, req *v1.GetDiscoverRequest) (*v1.GetDiscoverReply, error) {
	post := &Post{}

	res := &v1.GetDiscoverReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: &v1.DiscoverInfo{},
	}

	cacheData, err := u.getDiscoverInfoCache(ctx, req)
	if err != nil {
		u.log.Errorf("get discover info cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if cacheData != nil {
		res.Data = cacheData
		return res, nil
	}

	if err := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).First(post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发现内容不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userPostRepo := NewUserPostRepo(u.data, u.log.Logger())

	userInfo, err := userPostRepo.GetPostUserInfo(ctx, post.Id)
	if err != nil {
		u.log.Errorf("get post user info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catInfo, err := userPostRepo.GetPostCatInfo(ctx, post.Id)
	if err != nil {
		u.log.Errorf("get post cat info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	lostTime := ""
	if post.LostTime != nil {
		lostTime = post.LostTime.Format("2006-01-02 15:04:05")
	}

	res.Data = &v1.DiscoverInfo{
		Id:         post.Id,
		Title:      post.Title,
		PostType:   int32(post.PostType),
		ProvinceId: int32(post.ProvinceId),
		CityId:     int32(post.CityId),
		LostTime:   lostTime,
		Address:    post.Address,
		Cat: &v1.CatInfo{
			Id:        catInfo.Id,
			Name:      catInfo.Name,
			CatType:   catInfo.CatType,
			BreedType: catInfo.BreedType,
			Gender:    catInfo.Gender,
			Weight:    catInfo.Weight,
		},
		User: &v1.UserInfo{
			Id:         userInfo.Id,
			Name:       userInfo.Name,
			ProvinceId: int32(userInfo.ProvinceId),
			CityId:     int32(userInfo.CityId),
		},
		Remark:      post.Remark,
		CreatedTime: post.CreatedTime.Format("2006-01-02 15:04:05"),
		UpdatedTime: post.UpdatedTime.Format("2006-01-02 15:04:05"),
	}

	if err := u.setDiscoverInfoCache(ctx, req, res.Data); err != nil {
		u.log.Errorf("set discover info cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return res, nil
}

// getDiscoverInfoCacheKey 获取发现内容缓存key
func (u *discoverRepo) getDiscoverInfoCacheKey(req *v1.GetDiscoverRequest) string {
	return fmt.Sprintf("discover:info:%d", req.Id)
}

// getDiscoverInfoCache 获取发现内容缓存
func (u *discoverRepo) getDiscoverInfoCache(ctx context.Context, req *v1.GetDiscoverRequest) (*v1.DiscoverInfo, error) {
	redisKey := u.getDiscoverInfoCacheKey(req)

	jsonData, err := u.data.cache.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data v1.DiscoverInfo
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// setDiscoverInfoCache 设置发现内容缓存
func (u *discoverRepo) setDiscoverInfoCache(ctx context.Context, req *v1.GetDiscoverRequest, data *v1.DiscoverInfo) error {
	redisKey := u.getDiscoverInfoCacheKey(req)

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

// GetDiscoverRecommend 查询推荐内容
func (u *discoverRepo) GetDiscoverRecommend(ctx context.Context, req *v1.GetDiscoverRecommendRequest) (*v1.GetDiscoverRecommendReply, error) {
	res := &v1.GetDiscoverRecommendReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: []*v1.DiscoverRecommendItem{},
	}

	cacheData, err := u.getDiscoverRecommendCache(ctx, req)
	if err != nil {
		u.log.Errorf("get discover recommend cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if cacheData != nil {
		res.Data = cacheData
		return res, nil
	}

	query := u.data.db.Table("t_post as t1").Joins("inner join t_post_cat as t2 on t1.id = t2.post_id").Joins("inner join t_cat as t3 on t2.cat_id = t3.id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t3.deleted_flag = ?", 0).Where("t1.audit_status = ?", 1).Order("t1.created_time DESC").Limit(20)

	if req.Id > 0 {
		query = query.Where("t1.id != ?", req.Id)
	}

	rows, err := query.Select("t1.id", "t1.title", "t1.province_id", "t1.city_id").Rows()
	if err != nil {
		u.log.Errorf("get discover recommend failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	defer rows.Close()

	var items []*v1.DiscoverRecommendItem

	for rows.Next() {
		var (
			id         int64
			title      string
			provinceId int32
			cityId     int32
		)

		if err := rows.Scan(&id, &title, &provinceId, &cityId); err != nil {
			u.log.Errorf("get discover recommend failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		if provinceId == 0 || cityId == 0 {
			// 使用主人的省份和城市
			userPostRepo := NewUserPostRepo(u.data, u.log.Logger())
			userInfo, err := userPostRepo.GetPostUserInfo(ctx, id)
			if err != nil {
				u.log.Errorf("get post user info failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}
			provinceId = userInfo.ProvinceId
			cityId = userInfo.CityId
		}

		items = append(items, &v1.DiscoverRecommendItem{
			Id:         id,
			Title:      title,
			ProvinceId: provinceId,
			CityId:     cityId,
		})
	}

	if len(items) > 5 {
		// 随机获取5个
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		res.Data = items[:5]
	} else {
		res.Data = items
	}

	if err := rows.Err(); err != nil {
		u.log.Errorf("get discover recommend failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.setDiscoverRecommendCache(ctx, req, res.Data); err != nil {
		u.log.Errorf("set discover recommend cache failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return res, nil
}

// getDiscoverRecommendCacheKey 获取推荐内容缓存key
func (u *discoverRepo) getDiscoverRecommendCacheKey(req *v1.GetDiscoverRecommendRequest) string {
	if req.Id > 0 {
		return fmt.Sprintf("discover:recommend:%d", req.Id)
	}
	return "discover:recommend"
}

// getDiscoverRecommendCache 获取推荐内容缓存
func (u *discoverRepo) getDiscoverRecommendCache(ctx context.Context, req *v1.GetDiscoverRecommendRequest) ([]*v1.DiscoverRecommendItem, error) {
	redisKey := u.getDiscoverRecommendCacheKey(req)

	jsonData, err := u.data.cache.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var data []*v1.DiscoverRecommendItem
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	return data, nil
}

// setDiscoverRecommendCache 设置推荐内容缓存
func (u *discoverRepo) setDiscoverRecommendCache(ctx context.Context, req *v1.GetDiscoverRecommendRequest, data []*v1.DiscoverRecommendItem) error {
	redisKey := u.getDiscoverRecommendCacheKey(req)

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
