package data

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	v1 "xiaomiao-home-system/api/cat/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type catRepo struct {
	data *Data
	log  *log.Helper
}

// NewCatRepo .
func NewCatRepo(data *Data, logger log.Logger) biz.CatRepo {
	return &catRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "CatRepo")),
	}
}

// GetCatList 查询小猫列表
func (u *catRepo) GetCatList(ctx context.Context, req *v1.GetCatListRequest) (*v1.GetCatListReply, error) {
	var items []*v1.CatListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Table("t_cat as t1").Joins("inner join t_user_cat as t2 on t1.id = t2.cat_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Error("get cat list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.name", "t1.gender", "t1.breed_type", "t1.weight", "t1.birthday", "t1.neuter_status", "t1.health_status", "t1.dewormed_status", "t1.vaccine_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Error("get cat list failed: %v", err)
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
			u.log.Error("get cat list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.CatListItem{
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
		u.log.Error("get cat list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetCatListReply{
		Code: 200, Success: true,
		Data: &v1.CatList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CreateCat 创建小猫
func (u *catRepo) CreateCat(ctx context.Context, req *v1.CreateCatRequest) (*v1.CreateCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catId, err := u.data.gid.NextID()
	if err != nil {
		u.log.Error("generate cat id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		u.log.Error("parse birthday failed: %v", err)
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
				u.log.Error("parse vaccine last date failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			catInfo["vaccine_last_date"] = vaccineLastDate
		}
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&Cat{}).Create(catInfo).Error; err != nil {
		tx.Rollback()
		u.log.Error("create cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userCatId, err := u.data.gid.NextID()
	if err != nil {
		tx.Rollback()
		u.log.Error("generate user cat id failed: %v", err)
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
		u.log.Error("create user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("create cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreateCatReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// UpdateCat 更新小猫
func (u *catRepo) UpdateCat(ctx context.Context, req *v1.UpdateCatRequest) (*v1.UpdateCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前小猫是否属于当前用户
	userCatIds, err := u.GetCatIdsByUserId(ctx, userId)
	if err != nil {
		u.log.Error("get cat ids by user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !slices.Contains(userCatIds, req.Id) {
		u.log.Error("cat not found or no permission: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "小猫不存在或无权限")
	}

	// 更新小猫信息
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		u.log.Error("parse birthday failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catInfo := map[string]interface{}{
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
				u.log.Error("parse vaccine last date failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			catInfo["vaccine_last_date"] = vaccineLastDate
		}
	}

	tx := u.data.db.Begin()

	if err := tx.Model(&Cat{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(catInfo).Error; err != nil {
		tx.Rollback()
		u.log.Error("update cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("update cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateCatReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// DeleteCat 删除小猫
func (u *catRepo) DeleteCat(ctx context.Context, req *v1.DeleteCatRequest) (*v1.DeleteCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前小猫是否属于当前用户
	userCatIds, err := u.GetCatIdsByUserId(ctx, userId)
	if err != nil {
		u.log.Error("get cat ids by user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !slices.Contains(userCatIds, req.Id) {
		u.log.Error("cat not found or no permission: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "小猫不存在或无权限")
	}

	tx := u.data.db.Begin()

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	if err := tx.Model(&Cat{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		u.log.Error("delete cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Model(&UserCat{}).Where("cat_id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		u.log.Error("delete user cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Error("delete cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.DeleteCatReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// GetCatIdsByUserId 查询当前用户的小猫id
func (u *catRepo) GetCatIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var catIds []int64
	if err := u.data.db.Model(&UserCat{}).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Pluck("cat_id", &catIds).Error; err != nil {
		return nil, err
	}
	return catIds, nil
}

// GetCat 查询小猫信息
func (u *catRepo) GetCat(ctx context.Context, req *v1.GetCatRequest) (*v1.GetCatReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前小猫是否属于当前用户
	userCatIds, err := u.GetCatIdsByUserId(ctx, userId)
	if err != nil {
		u.log.Error("get cat ids by user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 判断当前小猫是否属于当前用户
	if !slices.Contains(userCatIds, req.Id) {
		u.log.Error("cat not found or no permission: %v", req.Id)
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
		u.log.Error("get cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catInfo := &v1.CatInfo{
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
		catInfo.VaccineLastDate = vaccineLastDate.Format("2006-01-02")
	}

	if vaccineTypes != "" {
		for _, vaccineType := range strings.Split(vaccineTypes, ",") {
			vaccineTypeInt, err := strconv.ParseInt(vaccineType, 10, 32)
			if err != nil {
				u.log.Error("parse vaccine type failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}
			catInfo.VaccineTypes = append(catInfo.VaccineTypes, int32(vaccineTypeInt))
		}
	}

	return &v1.GetCatReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: catInfo,
	}, nil
}
