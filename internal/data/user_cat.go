package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	v1 "xiaomiao-home-system/api/user/cat/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

const (
	userCatMaxCreateCount = 100
	userCatCreateCountTTL = 8 * time.Hour
	userCatMaxUpdateCount = 10
	userCatUpdateCountTTL = 8 * time.Hour
)

const redisKeyUserCatCreateCount = "user:cat:create:count"
const redisKeyUserCatUpdateCount = "user:cat:update:count"

type userCatRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserCatRepo .
func NewUserCatRepo(data *Data, logger log.Logger) biz.UserCatRepo {
	return &userCatRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserCatRepo")),
	}
}

// GetUserCatList 查询我的小猫列表
func (u *userCatRepo) GetUserCatList(ctx context.Context, req *v1.GetUserCatListRequest) (*v1.GetUserCatListReply, error) {
	var items []*v1.UserCatListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Table("t_cat as t1").Joins("inner join t_user_cat as t2 on t1.id = t2.cat_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user cat list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.name", "t1.gender", "t1.breed_type", "t1.weight", "t1.birthday", "t1.neuter_status", "t1.health_status", "t1.dewormed_status", "t1.vaccine_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get user cat list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id             int64
			name           string
			gender         int32
			breedType      int32
			weight         float32
			birthday       time.Time
			neuterStatus   int32
			healthStatus   int32
			dewormedStatus int32
			vaccineStatus  int32
			remark         string
			createdTime    time.Time
			updatedTime    time.Time
		)

		if err := rows.Scan(&id, &name, &gender, &breedType, &weight, &birthday, &neuterStatus, &healthStatus, &dewormedStatus, &vaccineStatus, &remark, &createdTime, &updatedTime); err != nil {
			u.log.Errorf("get user cat list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.UserCatListItem{
			Id:             id,
			Name:           name,
			Gender:         gender,
			BreedType:      breedType,
			Weight:         weight,
			Birthday:       birthday.Format("2006-01-02"),
			NeuterStatus:   neuterStatus,
			HealthStatus:   healthStatus,
			DewormedStatus: dewormedStatus,
			VaccineStatus:  vaccineStatus,
			Remark:         remark,
		})
	}

	if err := rows.Err(); err != nil {
		u.log.Errorf("get user cat list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserCatListReply{
		Code: 200, Success: true,
		Data: &v1.UserCatList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CreateUserCat 创建我的小猫
func (u *userCatRepo) CreateUserCat(ctx context.Context, req *v1.CreateUserCatRequest) (*v1.CreateUserCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.CheckUserCatCreateCountLimit(ctx, redisKeyUserCatCreateCount, userId); err != nil {
		return nil, err
	}

	catId, err := u.data.gid.NextID()
	if err != nil {
		u.log.Errorf("generate user cat id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		u.log.Errorf("parse birthday failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catInfo := map[string]interface{}{
		"id":              catId,
		"name":            req.Name,
		"gender":          req.Gender,
		"cat_type":        2,
		"breed_type":      req.BreedType,
		"weight":          req.Weight,
		"birthday":        birthday,
		"neuter_status":   req.NeuterStatus,
		"health_status":   req.HealthStatus,
		"dewormed_status": req.DewormedStatus,
		"vaccine_status":  req.VaccineStatus,
		"remark":          req.Remark,
		"created_time":    time.Now(),
		"updated_time":    time.Now(),
	}

	// 健康状态为: 2: 生病, 3: 残疾, 4: 其他, 需要配置健康状态信息
	if req.HealthStatus == 2 || req.HealthStatus == 3 || req.HealthStatus == 4 {
		catInfo["health_desc"] = req.HealthDesc
	}

	// 疫苗状态为: 1: 全程接种, 2: 部分接种, 需要配置疫苗欸写，最近接种日期，疫苗本凭证图片
	if req.VaccineStatus == 1 || req.VaccineStatus == 2 {
		catInfo["vaccine_cert_image"] = req.VaccineCertImage

		if len(req.VaccineTypes) > 0 {
			vaccineTypes := []string{}
			for _, vaccineType := range req.VaccineTypes {
				vaccineTypes = append(vaccineTypes, fmt.Sprintf("%d", vaccineType))
			}

			catInfo["vaccine_types"] = strings.Join(vaccineTypes, ",")
		}

		if req.VaccineLastDate != "" {
			vaccineLastDate, err := time.Parse("2006-01-02", req.VaccineLastDate)
			if err != nil {
				u.log.Errorf("parse vaccine last date failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			catInfo["vaccine_last_date"] = vaccineLastDate
		}
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&Cat{}).Create(catInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userCatId, err := u.data.gid.NextID()
	if err != nil {
		tx.Rollback()
		u.log.Errorf("generate user cat id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userCatInfo := map[string]interface{}{
		"id":           userCatId,
		"user_id":      userId,
		"cat_id":       catId,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := tx.Model(&UserCat{}).Create(userCatInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreateUserCatReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateUserCat 更新我的小猫
func (u *userCatRepo) UpdateUserCat(ctx context.Context, req *v1.UpdateUserCatRequest) (*v1.UpdateUserCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := u.CheckUserCatUpdateCountLimit(ctx, redisKeyUserCatUpdateCount, userId, req.Id); err != nil {
		return nil, err
	}

	// 查询当前小猫是否属于当前用户
	belongToUser, err := u.CheckUserCatBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user cat belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !belongToUser {
		u.log.Errorf("user cat not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "小猫不存在或无权限")
	}

	// 更新我的小猫信息
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		u.log.Errorf("parse birthday failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userCatInfo := map[string]interface{}{
		"name":            req.Name,
		"gender":          req.Gender,
		"cat_type":        2,
		"breed_type":      req.BreedType,
		"weight":          req.Weight,
		"birthday":        birthday,
		"neuter_status":   req.NeuterStatus,
		"health_status":   req.HealthStatus,
		"dewormed_status": req.DewormedStatus,
		"vaccine_status":  req.VaccineStatus,
		"remark":          req.Remark,
		"updated_time":    time.Now(),
	}

	// 健康状态为: 2: 生病, 3: 残疾, 4: 其他, 需要配置健康状态信息
	if req.HealthStatus == 2 || req.HealthStatus == 3 || req.HealthStatus == 4 {
		userCatInfo["health_desc"] = req.HealthDesc
	}

	// 疫苗状态为: 1: 全程接种, 2: 部分接种, 需要配置疫苗欸写，最近接种日期，疫苗本凭证图片
	if req.VaccineStatus == 1 || req.VaccineStatus == 2 {
		userCatInfo["vaccine_cert_image"] = req.VaccineCertImage

		if len(req.VaccineTypes) > 0 {
			vaccineTypes := []string{}
			for _, vaccineType := range req.VaccineTypes {
				vaccineTypes = append(vaccineTypes, fmt.Sprintf("%d", vaccineType))
			}

			userCatInfo["vaccine_types"] = strings.Join(vaccineTypes, ",")
		}

		if req.VaccineLastDate != "" {
			vaccineLastDate, err := time.Parse("2006-01-02", req.VaccineLastDate)
			if err != nil {
				u.log.Errorf("parse vaccine last date failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			userCatInfo["vaccine_last_date"] = vaccineLastDate
		}
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&Cat{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(userCatInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("update user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("update user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserCatReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// DeleteUserCat 删除我的小猫
func (u *userCatRepo) DeleteUserCat(ctx context.Context, req *v1.DeleteUserCatRequest) (*v1.DeleteUserCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前小猫是否属于当前用户
	belongToUser, err := u.CheckUserCatBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user cat belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !belongToUser {
		u.log.Errorf("user cat not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "小猫不存在或无权限")
	}

	tx := u.data.db.Begin()

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	if err := tx.Model(&Cat{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("delete user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Model(&UserCat{}).Where("cat_id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("delete user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("delete user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.DeleteUserCatReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// CheckUserCatBelongToUser 检查当前小猫是否属于当前用户
func (u *userCatRepo) CheckUserCatBelongToUser(ctx context.Context, userId int64, catId int64) (bool, error) {
	var count int64

	redisKey := fmt.Sprintf("user:cat:belong:%d:%d", userId, catId)

	n, err := u.data.cache.Get(ctx, redisKey).Int()
	if err != nil {
		if err != redis.Nil {
			u.log.Errorf("get user cat belong to user cache error: %v", err)
			return false, nil
		}
	} else {
		return n > 0, nil
	}

	if err := u.data.db.Model(&UserCat{}).Where("user_id = ?", userId).Where("cat_id = ?", catId).Where("deleted_flag = ?", 0).Count(&count).Error; err != nil {
		u.log.Errorf("get user cat belong to user db error: %v", err)
		return false, err
	}

	if err := u.data.cache.Set(ctx, redisKey, count, 5*time.Minute).Err(); err != nil {
		u.log.Errorf("set user cat belong to user cache error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// GetUserCat 查询我的小猫信息
func (u *userCatRepo) GetUserCat(ctx context.Context, req *v1.GetUserCatRequest) (*v1.GetUserCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前小猫是否属于当前用户
	belongToUser, err := u.CheckUserCatBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user cat belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !belongToUser {
		u.log.Errorf("user cat not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "小猫不存在或无权限")
	}

	row := u.data.db.Model(&Cat{}).Select("id", "name", "gender", "breed_type", "weight", "birthday", "neuter_status", "health_status", "health_desc", "dewormed_status", "vaccine_status", "vaccine_types", "vaccine_last_date", "vaccine_cert_image", "remark").Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Limit(1).Row()

	var (
		id               int64
		name             string
		gender           int32
		breedType        int32
		weight           float32
		birthday         time.Time
		neuterStatus     int32
		healthStatus     int32
		healthDesc       string
		dewormedStatus   int32
		vaccineStatus    int32
		vaccineTypes     string
		vaccineLastDate  *time.Time
		vaccineCertImage string
		remark           string
	)

	if err := row.Scan(&id, &name, &gender, &breedType, &weight, &birthday, &neuterStatus, &healthStatus, &healthDesc, &dewormedStatus, &vaccineStatus, &vaccineTypes, &vaccineLastDate, &vaccineCertImage, &remark); err != nil {
		u.log.Errorf("get user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userCatInfo := &v1.UserCatInfo{
		Id:               id,
		Name:             name,
		Gender:           gender,
		BreedType:        breedType,
		Weight:           weight,
		Birthday:         birthday.Format("2006-01-02"),
		NeuterStatus:     neuterStatus,
		HealthStatus:     healthStatus,
		HealthDesc:       healthDesc,
		DewormedStatus:   dewormedStatus,
		VaccineStatus:    vaccineStatus,
		VaccineTypes:     []int32{},
		VaccineCertImage: vaccineCertImage,
		Remark:           remark,
	}

	if vaccineLastDate != nil {
		userCatInfo.VaccineLastDate = vaccineLastDate.Format("2006-01-02")
	}

	if vaccineTypes != "" {
		for _, vaccineType := range strings.Split(vaccineTypes, ",") {
			vaccineTypeInt, err := strconv.ParseInt(vaccineType, 10, 32)
			if err != nil {
				u.log.Errorf("parse vaccine type failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}
			userCatInfo.VaccineTypes = append(userCatInfo.VaccineTypes, int32(vaccineTypeInt))
		}
	}

	return &v1.GetUserCatReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: userCatInfo,
	}, nil
}

// CheckUserCatCreateCountLimit 检查我的小猫创建次数限制
func (u *userCatRepo) CheckUserCatCreateCountLimit(ctx context.Context, countKeyPrefix string, userId int64) error {
	key := fmt.Sprintf("%s:%d", countKeyPrefix, userId)
	n, err := u.data.cache.Incr(ctx, key).Result()
	if err != nil {
		u.log.Errorf("increase user cat create count failed, key=%s: %v", countKeyPrefix, err)
		return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if n == 1 {
		if err := u.data.cache.Expire(ctx, key, userCatCreateCountTTL).Err(); err != nil {
			u.log.Errorf("set user cat create count ttl failed, key=%s: %v", countKeyPrefix, err)
			return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if n > userCatMaxCreateCount {
		return errors.BadRequest(v1.ErrorReason_ERR_TOO_MANY_REQUEST.String(), "今日已达最大创建次数, 请明天再试")
	}

	return nil
}

// CheckUserCatUpdateCountLimit 检查我的小猫更新次数限制
func (u *userCatRepo) CheckUserCatUpdateCountLimit(ctx context.Context, countKeyPrefix string, userId int64, catId int64) error {
	key := fmt.Sprintf("%s:%d:%d", countKeyPrefix, userId, catId)
	n, err := u.data.cache.Incr(ctx, key).Result()
	if err != nil {
		u.log.Errorf("increase user cat update count failed, key=%s: %v", countKeyPrefix, err)
		return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if n == 1 {
		if err := u.data.cache.Expire(ctx, key, userCatUpdateCountTTL).Err(); err != nil {
			u.log.Errorf("set user cat action count ttl failed, key=%s: %v", countKeyPrefix, err)
			return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if n > userCatMaxUpdateCount {
		return errors.BadRequest(v1.ErrorReason_ERR_TOO_MANY_REQUEST.String(), "今日已达最大更新次数, 请明日再试")
	}

	return nil
}

// GetUserCats 查询用户所有小猫
func (u *userCatRepo) GetUserCats(ctx context.Context, req *v1.GetUserCatsRequest) (*v1.GetUserCatsReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	items := make([]*v1.CatItem, 0)
	total := int64(0)

	baseQuery := u.data.db.Table("t_cat as t1").Joins("inner join t_user_cat as t2 on t1.id = t2.cat_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user cats failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.name", "t1.gender", "t1.breed_type").Order("t1.created_time DESC")

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get user cats failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id        int64
			name      string
			gender    int32
			breedType int32
		)

		if err := rows.Scan(&id, &name, &gender, &breedType); err != nil {
			u.log.Errorf("get user cats failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.CatItem{
			Id:        id,
			Name:      name,
			Gender:    gender,
			BreedType: breedType,
		})
	}

	if err := rows.Err(); err != nil {
		u.log.Errorf("get user cats failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserCatsReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: items,
	}, nil
}
