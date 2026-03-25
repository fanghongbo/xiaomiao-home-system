package data

import (
	"context"
	"fmt"
	v1 "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

var (
	ErrBuildForwardRequest = errors.NotFound(v1.ErrorReason_ERR_BUILD_FORWARD_REQUEST.String(), "failed to build forward request")
	ErrForwardRequest      = errors.NotFound(v1.ErrorReason_ERR_FORWARD_REQUEST.String(), "failed to forward request")
)

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserRepo")),
	}
}

// Login 登录接口
func (u *userRepo) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {

	fmt.Println(req.Username)
	fmt.Println(req.Password)

	settingRepo := NewSystemSettingRepo(u.data, u.log.Logger())

	siteName, err := settingRepo.GetSiteName(ctx)
	if err != nil || siteName == "" {
		siteName = "单点登录系统"
	}

	invalidLoginCount, err := settingRepo.GetInvalidLoginCount(ctx)
	if err != nil || invalidLoginCount <= 3 {
		invalidLoginCount = 3
	}

	loginBanDuration, err := settingRepo.GetLoginBanDuration(ctx)
	if err != nil || loginBanDuration <= 5 {
		loginBanDuration = 5
	}

	fmt.Println(siteName)
	fmt.Println(invalidLoginCount)
	fmt.Println(loginBanDuration)

	userInfo, err := u.CheckLogin(ctx, req, invalidLoginCount, loginBanDuration)
	if err != nil {
		return nil, err
	}

	fmt.Println("userInfo =>", userInfo)

	return &v1.LoginReply{}, nil
}

// GetUserId 查询当前用户ID
func (u *userRepo) GetUserId(ctx context.Context) (int64, error) {
	username, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return 0, err
	}

	var user User
	if err := u.data.db.Model(&User{}).Where("username = ?", username).Where("deleted_flag = ?", 0).First(&user).Error; err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (u *userRepo) GetUserByUsername(ctx context.Context, username string) (*v1.UserInfo, error) {
	var user User

	if err := u.data.db.Model(&User{}).Where("username = ?", username).Where("deleted_flag = ?", 0).First(&user).Error; err != nil {
		return nil, err
	}

	return &v1.UserInfo{
		Id:       user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Roles:    []*v1.RoleItem{},
		Status:   int64(user.Status),
	}, nil
}

func (u *userRepo) CheckLogin(ctx context.Context, login *v1.LoginRequest, maxInvalidLoginCount int64, loginBanDuration int64) (*v1.UserInfo, error) {
	maxInvalidLoginCount += 1

	userInfo, err := u.GetUserByUsername(ctx, login.Username)
	if err != nil {
		u.log.Errorf("获取用户信息失败: %v", err)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	userIp, err := utils.GetUserIP(ctx)
	if err != nil {
		u.log.Errorf("获取用户IP失败: %v", err)
		return nil, fmt.Errorf("获取用户IP失败")
	}

	fmt.Println("ip =>", userIp)

	return userInfo, nil
}
