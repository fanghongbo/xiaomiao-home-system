package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/user/post/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type userPostRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserPostRepo .
func NewUserPostRepo(data *Data, logger log.Logger) biz.UserPostRepo {
	return &userPostRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserPostRepo")),
	}
}

// GetUserPostList 查询发布内容列表
func (u *userPostRepo) GetUserPostList(ctx context.Context, req *v1.GetUserPostListRequest) (*v1.GetUserPostListReply, error) {
	var items []*v1.UserPostListItem
	var total int64

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	baseQuery := u.data.db.Table("t_post as t1").Joins("inner join t_user_post as t2 on t1.id = t2.post_id").Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Where("t2.user_id = ?", userId)

	if req.PType > 0 {
		baseQuery = baseQuery.Where("t1.post_type = ?", req.PType)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		u.log.Errorf("get user post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	result := baseQuery.Select("t1.id", "t1.title", "t1.post_status", "t1.audit_status", "t1.remark", "t1.created_time", "t1.updated_time").Order("t1.created_time DESC").
		Limit(int(req.Size)).
		Offset(int((req.Page - 1) * req.Size))

	rows, err := result.Rows()
	if err != nil {
		u.log.Errorf("get user post list failed: %v", err)
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
			u.log.Errorf("get user post list failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		items = append(items, &v1.UserPostListItem{
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
		u.log.Errorf("get user post list failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetUserPostListReply{
		Code: 200, Success: true,
		Data: &v1.UserPostList{
			Items: items,
			Total: total,
		},
	}, nil
}

// CheckPostParams 检查发布内容参数
func (u *userPostRepo) CheckCreateUserPostParams(req *v1.CreateUserPostRequest) error {
	if len(req.Title) > 64 {
		return fmt.Errorf("标题长度不能超过64个字符")
	}

	if len(req.Title) < 1 {
		return fmt.Errorf("标题不能为空")
	}

	// postType
	if req.PostType < 1 || req.PostType > 4 {
		return fmt.Errorf("发布类型错误")
	}

	// remark
	if len(req.Remark) > 500 || len(req.Remark) < 1 {
		return fmt.Errorf("描述长度不能超过500个字符，且不能为空")
	}

	switch req.PostType {
	case 1: // 领养
		if req.CatType == 1 {
			// BreedType
			if req.BreedType < 1 || req.BreedType > 10 {
				return fmt.Errorf("品种类型错误")
			}

			// Gender
			if req.Gender < 0 || req.Gender > 2 {
				return fmt.Errorf("性别错误")
			}
		}

		// CityId
		if req.CityId < 1 {
			return fmt.Errorf("城市错误")
		}

		// ProvinceId
		if req.ProvinceId < 1 {
			return fmt.Errorf("省份错误")
		}

		// Address
		if len(req.Address) > 255 || len(req.Address) < 1 {
			return fmt.Errorf("地址长度不能超过255个字符，且不能为空")
		}
	case 2: // 寻猫
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}

		// LostTime
		if len(req.LostTime) < 1 {
			return fmt.Errorf("丢失时间不能为空")
		}

		// CityId
		if req.CityId < 1 {
			return fmt.Errorf("城市错误")
		}

		// ProvinceId
		if req.ProvinceId < 1 {
			return fmt.Errorf("省份错误")
		}

		// Address
		if len(req.Address) > 255 || len(req.Address) < 1 {
			return fmt.Errorf("地址长度不能超过255个字符，且不能为空")
		}
	case 3: // 日常
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}
	case 4: // 求助
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}
	default:
		return fmt.Errorf("发布类型错误")
	}

	return nil
}

// CreatePost 创建发布内容
func (u *userPostRepo) CreateUserPost(ctx context.Context, req *v1.CreateUserPostRequest) (*v1.CreateUserPostReply, error) {
	postId, err := u.data.gid.NextID()
	if err != nil {
		u.log.Errorf("generate id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 检查发布内容参数
	err = u.CheckCreateUserPostParams(req)
	if err != nil {
		u.log.Errorf("check post params failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), err.Error())
	}

	postInfo := map[string]interface{}{
		"id":           postId,
		"title":        req.Title,
		"post_type":    req.PostType,
		"remark":       req.Remark,
		"audit_status": 0,
		"post_status":  0,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	switch req.PostType {
	case 1: // 领养
		postInfo["city_id"] = req.CityId
		postInfo["province_id"] = req.ProvinceId
		postInfo["address"] = req.Address
		postInfo["lost_time"] = nil
	case 2: // 寻猫
		postInfo["lost_time"] = req.LostTime
		postInfo["city_id"] = req.CityId
		postInfo["province_id"] = req.ProvinceId
		postInfo["address"] = req.Address
	case 3: // 日常
		postInfo["lost_time"] = nil
		postInfo["city_id"] = 0
		postInfo["province_id"] = 0
		postInfo["address"] = nil
	case 4: // 求助
		postInfo["lost_time"] = nil
		postInfo["city_id"] = 0
		postInfo["province_id"] = 0
		postInfo["address"] = nil
	default:
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布类型错误")
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Post{}).Create(postInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 创建用户发布内容
	userPostId, err := u.data.gid.NextID()
	if err != nil {
		tx.Rollback()
		u.log.Errorf("generate user post id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userPostInfo := map[string]interface{}{
		"id":           userPostId,
		"user_id":      userId,
		"post_id":      postId,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := tx.Model(&UserPost{}).Create(userPostInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create user post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 领养且为流浪猫类型的时候，需要创建小猫信息
	if req.PostType == 1 && req.CatType == 1 {
		catId, err := u.data.gid.NextID()
		if err != nil {
			tx.Rollback()
			u.log.Errorf("generate cat id failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		catInfo := map[string]interface{}{
			"id":           catId,
			"name":         "流浪猫",
			"gender":       req.Gender,
			"cat_type":     1,
			"breed_type":   req.BreedType,
			"created_time": time.Now(),
			"updated_time": time.Now(),
		}

		if err := tx.Model(&Cat{}).Create(catInfo).Error; err != nil {
			tx.Rollback()
			u.log.Errorf("create cat failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		req.CatId = int64(catId)
	}

	// 创建发布内容与小猫的关联
	postCatId, err := u.data.gid.NextID()
	if err != nil {
		tx.Rollback()
		u.log.Errorf("generate post cat id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	postCatInfo := map[string]interface{}{
		"id":           postCatId,
		"post_id":      postId,
		"cat_id":       req.CatId,
		"created_time": time.Now(),
		"updated_time": time.Now(),
	}

	if err := tx.Model(&PostCat{}).Create(postCatInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post cat failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.CreateUserPostReply{
		Code: 200, Success: true, Message: "创建成功",
	}, nil
}

// CheckUpdateUserPostParams 检查发布内容参数
func (u *userPostRepo) CheckUpdateUserPostParams(req *v1.UpdateUserPostRequest) error {
	if len(req.Title) > 64 {
		return fmt.Errorf("标题长度不能超过64个字符")
	}

	if len(req.Title) < 1 {
		return fmt.Errorf("标题不能为空")
	}

	// postType
	if req.PostType < 1 || req.PostType > 4 {
		return fmt.Errorf("发布类型错误")
	}

	// remark
	if len(req.Remark) > 500 || len(req.Remark) < 1 {
		return fmt.Errorf("描述长度不能超过500个字符，且不能为空")
	}

	switch req.PostType {
	case 1: // 领养
		if req.CatType == 1 {
			// BreedType
			if req.BreedType < 1 || req.BreedType > 10 {
				return fmt.Errorf("品种类型错误")
			}

			// Gender
			if req.Gender < 0 || req.Gender > 2 {
				return fmt.Errorf("性别错误")
			}
		}

		// CityId
		if req.CityId < 1 {
			return fmt.Errorf("城市错误")
		}

		// ProvinceId
		if req.ProvinceId < 1 {
			return fmt.Errorf("省份错误")
		}

		// Address
		if len(req.Address) > 255 || len(req.Address) < 1 {
			return fmt.Errorf("地址长度不能超过255个字符，且不能为空")
		}
	case 2: // 寻猫
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}

		// LostTime
		if len(req.LostTime) < 1 {
			return fmt.Errorf("丢失时间不能为空")
		}

		// CityId
		if req.CityId < 1 {
			return fmt.Errorf("城市错误")
		}

		// ProvinceId
		if req.ProvinceId < 1 {
			return fmt.Errorf("省份错误")
		}

		// Address
		if len(req.Address) > 255 || len(req.Address) < 1 {
			return fmt.Errorf("地址长度不能超过255个字符，且不能为空")
		}
	case 3: // 日常
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}
	case 4: // 求助
		if req.CatId < 1 {
			return fmt.Errorf("小猫id错误")
		}
	default:
		return fmt.Errorf("发布类型错误")
	}

	return nil
}

// UpdateUserPost 更新发布内容
func (u *userPostRepo) UpdateUserPost(ctx context.Context, req *v1.UpdateUserPostRequest) (*v1.UpdateUserPostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	// 检查发布内容参数
	err = u.CheckUpdateUserPostParams(req)
	if err != nil {
		u.log.Errorf("check post params failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), err.Error())
	}

	postInfo := map[string]interface{}{
		"title":        req.Title,
		"post_type":    req.PostType,
		"remark":       req.Remark,
		"audit_status": 0,
		"post_status":  0,
		"updated_time": time.Now(),
	}

	switch req.PostType {
	case 1: // 领养
		postInfo["city_id"] = req.CityId
		postInfo["province_id"] = req.ProvinceId
		postInfo["address"] = req.Address
		postInfo["lost_time"] = nil
	case 2: // 寻猫
		postInfo["lost_time"] = req.LostTime
		postInfo["city_id"] = req.CityId
		postInfo["province_id"] = req.ProvinceId
		postInfo["address"] = req.Address
	case 3: // 日常
		postInfo["lost_time"] = nil
		postInfo["city_id"] = 0
		postInfo["province_id"] = 0
		postInfo["address"] = nil
	case 4: // 求助
		postInfo["lost_time"] = nil
		postInfo["city_id"] = 0
		postInfo["province_id"] = 0
		postInfo["address"] = nil
	default:
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布类型错误")
	}

	// 查询历史已关联的小猫信息
	catInfo, err := u.GetPostCatInfo(ctx, req.Id)
	if err != nil {
		u.log.Errorf("get post cat info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 启动MySQL事务
	tx := u.data.db.Begin()

	if err := tx.Model(&Post{}).Where("id = ?", req.Id).Updates(postInfo).Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 领养且为流浪猫类型的时候，检查是否需要重新创建小猫信息
	if req.PostType == 1 && req.CatType == 1 {
		// 历史关联的小猫信息如果是流浪猫类型，则直接更新这个小猫信息
		if catInfo.CatType == 1 {
			updateCatInfo := map[string]interface{}{
				"name":         "流浪猫",
				"gender":       req.Gender,
				"cat_type":     1,
				"breed_type":   req.BreedType,
				"updated_time": time.Now(),
			}

			if err := tx.Model(&Cat{}).Where("id = ?", catInfo.Id).Updates(updateCatInfo).Error; err != nil {
				tx.Rollback()
				u.log.Errorf("update cat failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			req.CatId = catInfo.Id
		} else {
			// 如果这个小猫信息是其他类型，则需要删除这个小猫信息, 并重新创建一个流浪猫信息
			catId, err := u.data.gid.NextID()
			if err != nil {
				tx.Rollback()
				u.log.Errorf("generate cat id failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			catInfo := map[string]interface{}{
				"id":           catId,
				"name":         "流浪猫",
				"gender":       req.Gender,
				"cat_type":     1,
				"breed_type":   req.BreedType,
				"created_time": time.Now(),
				"updated_time": time.Now(),
			}

			if err := tx.Model(&Cat{}).Create(catInfo).Error; err != nil {
				tx.Rollback()
				u.log.Errorf("create cat failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}

			req.CatId = int64(catId)
		}
	}

	// 检查是否需要重新关联小猫和发布内容
	// 如果本次修改的小猫信息和历史小猫信息不一致，则先删除历史小猫和发布内容的关联, 如果一致则不需要变更
	if catInfo.Id != req.CatId {
		updateInfo := map[string]interface{}{
			"deleted_flag": 1,
			"deleted_time": time.Now(),
		}

		// 删除历史流浪猫信息
		if catInfo.CatType == 1 {
			if err := tx.Model(&Cat{}).Where("id = ?", catInfo.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
				tx.Rollback()
				u.log.Errorf("delete cat failed: %v", err)
				return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
			}
		}

		// 删除历史关联的小猫和发布内容的关联信息
		if err := tx.Model(&PostCat{}).Where("post_id = ?", req.Id).Where("cat_id = ?", catInfo.Id).Where("deleted_flag = ?", 0).Updates(updateInfo).Error; err != nil {
			tx.Rollback()
			u.log.Errorf("delete post cat failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		// 重新创建发布内容与小猫的关联
		postCatId, err := u.data.gid.NextID()
		if err != nil {
			tx.Rollback()
			u.log.Errorf("generate post cat id failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}

		postCatInfo := map[string]interface{}{
			"id":           postCatId,
			"post_id":      req.Id,
			"cat_id":       req.CatId,
			"created_time": time.Now(),
			"updated_time": time.Now(),
		}

		if err := tx.Model(&PostCat{}).Create(postCatInfo).Error; err != nil {
			tx.Rollback()
			u.log.Errorf("create post cat failed: %v", err)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("create post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserPostReply{
		Code: 200, Success: true, Message: "修改成功",
	}, nil
}

// DeleteUserPost 删除发布内容
func (u *userPostRepo) DeleteUserPost(ctx context.Context, req *v1.DeleteUserPostRequest) (*v1.DeleteUserPostReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	updateInfo := map[string]interface{}{
		"deleted_flag": 1,
		"deleted_time": time.Now(),
	}

	// 查询历史已关联的小猫信息
	catInfo, err := u.GetPostCatInfo(ctx, req.Id)
	if err != nil {
		u.log.Errorf("get post cat info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	tx := u.data.db.Begin()

	// 删除发布内容
	result := tx.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(updateInfo)

	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete post failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	// 删除用户发布内容关联信息
	result = tx.Model(&UserPost{}).Where("post_id = ?", req.Id).Where("user_id = ?", userId).Where("deleted_flag = ?", 0).Updates(updateInfo)
	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete user post failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	// 删除发布内容和小猫关联
	result = tx.Model(&PostCat{}).Where("post_id = ?", req.Id).Where("cat_id = ?", catInfo.Id).Where("deleted_flag = ?", 0).Updates(updateInfo)
	if result.Error != nil {
		tx.Rollback()
		u.log.Errorf("delete post cat failed: %v", result.Error)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容和小猫关联不存在或无权限")
	}

	// 删除小猫信息
	if catInfo.CatType == 1 {
		result = tx.Model(&Cat{}).Where("id = ?", catInfo.Id).Where("deleted_flag = ?", 0).Updates(updateInfo)
		if result.Error != nil {
			tx.Rollback()
			u.log.Errorf("delete cat failed: %v", result.Error)
			return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		u.log.Errorf("delete post failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.DeleteUserPostReply{
		Code: 200, Success: true, Message: "删除成功",
	}, nil
}

// UpdatePostStatus 更新发布内容状态
func (u *userPostRepo) UpdateUserPostStatus(ctx context.Context, req *v1.UpdateUserPostStatusRequest) (*v1.UpdateUserPostStatusReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	post := map[string]interface{}{
		"status": int(req.Status),
	}

	result := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).Updates(post)

	if result.Error != nil {
		return nil, result.Error
	}

	return &v1.UpdateUserPostStatusReply{
		Code: 200, Success: true, Message: "更新成功",
	}, nil
}

// GetUserPost 查询发布内容
func (u *userPostRepo) GetUserPost(ctx context.Context, req *v1.GetUserPostRequest) (*v1.GetUserPostReply, error) {
	post := &Post{}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Errorf("get current user id failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	// 查询当前发布内容是否属于当前用户
	belongToUser, err := u.CheckUserPostBelongToUser(ctx, userId, req.Id)
	if err != nil {
		u.log.Errorf("check user post belong to user failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if !belongToUser {
		u.log.Errorf("user post not belong to user: %v", req.Id)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在或无权限")
	}

	if err := u.data.db.Model(&Post{}).Where("id = ?", req.Id).Where("deleted_flag = ?", 0).First(post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "发布内容不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	userInfo, err := u.GetPostUserInfo(ctx, post.Id)
	if err != nil {
		u.log.Errorf("get post user info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	catInfo, err := u.GetPostCatInfo(ctx, post.Id)
	if err != nil {
		u.log.Errorf("get post cat info failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	lostTime := ""
	if post.LostTime != nil {
		lostTime = post.LostTime.Format("2006-01-02 15:04:05")
	}

	return &v1.GetUserPostReply{
		Code: 200, Success: true, Message: "查询成功",
		Data: &v1.UserPostInfo{Id: post.Id, Title: post.Title, PostType: int32(post.PostType), ProvinceId: int32(post.ProvinceId), CityId: int32(post.CityId), LostTime: lostTime, Address: post.Address, Cat: catInfo, User: userInfo, Remark: post.Remark},
	}, nil
}

// CheckUserPostBelongToUser 检查当前发布内容是否属于当前用户
func (u *userPostRepo) CheckUserPostBelongToUser(ctx context.Context, userId int64, postId int64) (bool, error) {
	var count int64

	redisKey := fmt.Sprintf("user:post:belong:%d:%d", userId, postId)

	n, err := u.data.cache.Get(ctx, redisKey).Int()
	if err != nil {
		if err != redis.Nil {
			u.log.Errorf("get user post belong to user cache error: %v", err)
			return false, nil
		}
	} else {
		return n > 0, nil
	}

	if err := u.data.db.Model(&UserPost{}).Where("user_id = ?", userId).Where("post_id = ?", postId).Where("deleted_flag = ?", 0).Count(&count).Error; err != nil {
		u.log.Errorf("get user post belong to user db error: %v", err)
		return false, err
	}

	if err := u.data.cache.Set(ctx, redisKey, count, 5*time.Minute).Err(); err != nil {
		u.log.Errorf("set user post belong to user cache error: %v", err)
		return false, err
	}

	return count > 0, nil
}

// GetPostCatInfo 查询发布内容小猫信息
func (u *userPostRepo) GetPostCatInfo(ctx context.Context, postId int64) (*v1.CatInfo, error) {
	var (
		id        int64
		name      string
		catType   int32
		breedType int32
		gender    int32
		weight    float32
	)

	row := u.data.db.Table("t_cat as t1").Joins("inner join t_post_cat as t2 on t1.id = t2.cat_id").Where("t2.post_id = ?", postId).Where("t2.deleted_flag = ?", 0).Where("t1.deleted_flag = ?", 0).Select("t1.id", "t1.name", "t1.cat_type", "t1.breed_type", "t1.gender", "t1.weight").Limit(1).Row()

	if err := row.Scan(&id, &name, &catType, &breedType, &gender, &weight); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("小猫不存在")
		}
		return nil, err
	}

	return &v1.CatInfo{
		Id:        id,
		Name:      name,
		CatType:   catType,
		BreedType: breedType,
		Gender:    gender,
		Weight:    weight,
	}, nil
}

// GetPostUserInfo 查询发布内容用户信息
func (u *userPostRepo) GetPostUserInfo(ctx context.Context, postId int64) (*v1.UserInfo, error) {
	var (
		id   int64
		name string
	)
	row := u.data.db.Table("t_user as t1").Joins("inner join t_user_post as t2 on t1.id = t2.user_id").Where("t2.post_id = ?", postId).Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0).Select("t1.id", "t1.nickname").Limit(1).Row()
	if err := row.Scan(&id, &name); err != nil {
		u.log.Errorf("get post user info failed: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户不存在")
		}
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	return &v1.UserInfo{
		Id:   id,
		Name: name,
	}, nil
}
