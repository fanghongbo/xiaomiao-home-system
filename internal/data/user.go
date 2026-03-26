package data

import (
	"context"
	"fmt"
	"time"
	v1 "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/third_party/jwt"
	"xiaomiao-home-system/third_party/password"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"xiaomiao-home-system/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "UserRepo")),
	}
}

// 判断登录身份类型
func (u *userRepo) CheckLoginIdentityType(loginIdentity string) v1.LoginType {

	// 检查登录身份是否是邮箱地址
	if utils.IsEmail(loginIdentity) {
		return v1.LoginType_LOGIN_TYPE_EMAIL
	}

	// 检查登录身份是否是手机号
	if utils.IsMobile(loginIdentity) {
		return v1.LoginType_LOGIN_TYPE_SMS
	}

	// 检查登录身份是否是账号
	return v1.LoginType_LOGIN_TYPE_ACCOUNT
}

// Login 登录接口
func (u *userRepo) WebLogin(ctx context.Context, req *v1.WebLoginRequest) (*v1.WebLoginReply, error) {
	settingRepo := NewSystemSettingRepo(u.data, u.log.Logger())

	// 获取限制登录错误次数
	invalidLoginCount, err := settingRepo.GetInvalidLoginCount(ctx)
	if err != nil || invalidLoginCount <= 3 {
		invalidLoginCount = 3
	}

	// 获取限制登录时间间隔
	loginBanDuration, err := settingRepo.GetLoginBanDuration(ctx)
	if err != nil || loginBanDuration <= 5 {
		loginBanDuration = 5
	}

	// 登录有效期
	loginValidPeriod, err := settingRepo.GetLoginValidPeriod(ctx)
	if err != nil || loginValidPeriod <= 8 {
		loginValidPeriod = 8
	}

	clientIp, err := utils.GetUserIP(ctx)
	if err != nil {
		u.log.Error("get user ip failed: %v", err)
		clientIp = "127.0.0.1"
	}

	// 获取登录错误次数
	loginErrorCount, err := u.GetLoginErrorCount(ctx, req.LoginType, req.LoginIdentity, clientIp)
	if err != nil {
		u.log.Error("get login error count failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	if loginErrorCount >= invalidLoginCount {
		u.log.Error("login error count too many: %d", loginErrorCount)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "登录错误次数过多, 请稍后再试")
	}

	var res *v1.WebLoginReply

	switch req.LoginType {
	case v1.LoginType_LOGIN_TYPE_ACCOUNT:
		res, err = u.WebLoginAccount(ctx, req, loginValidPeriod)
	case v1.LoginType_LOGIN_TYPE_SMS:
		res, err = u.WebLoginSms(ctx, req)
	default:
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "暂不支持当前登录方式, 请联系管理员")
	}

	if err != nil {
		u.log.Error("web login failed: %v", err)

		if err = u.IncLoginErrorCount(ctx, req.LoginType, req.LoginIdentity, clientIp, time.Duration(loginBanDuration)*time.Minute); err != nil {
			u.log.Error("inc login error count failed: %v", err)
		}

		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "账号或密码错误, 请稍后再试")
	}

	return res, nil
}

// GetLoginErrorCount 查询登录错误次数
func (u *userRepo) GetLoginErrorCount(ctx context.Context, loginType v1.LoginType, loginIdentity string, clientIp string) (int64, error) {
	redisKey := fmt.Sprintf("login:error:count:%d:%s:%s", loginType.Number(), loginIdentity, clientIp)

	count, err := u.data.rdb.Get(ctx, redisKey).Int64()
	if err != nil && err != redis.Nil {
		u.log.Error("get login error count failed: %v", err)
		return 0, err
	}

	return count, nil
}

// IncLoginErrorCount 增加登录错误次数
func (u *userRepo) IncLoginErrorCount(ctx context.Context, loginType v1.LoginType, loginIdentity string, clientIp string, ttl time.Duration) error {
	redisKey := fmt.Sprintf("login:error:count:%d:%s:%s", loginType.Number(), loginIdentity, clientIp)

	_, err := u.data.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		u.log.Error("inc login error count failed: %v", err)
		return err
	}

	if err = u.data.rdb.Expire(ctx, redisKey, ttl).Err(); err != nil {
		u.log.Error("set login error count expire failed: %v", err)
		return err
	}

	return nil
}

// ClearLoginErrorCount 清空登录错误次数
func (u *userRepo) ClearLoginErrorCount(ctx context.Context, loginType v1.LoginType, loginIdentity string, clientIp string) error {
	redisKey := fmt.Sprintf("login:error:count:%d:%s:%s", loginType.Number(), loginIdentity, clientIp)

	err := u.data.rdb.Del(ctx, redisKey).Err()
	if err != nil {
		u.log.Error("clear login error count failed: %v", err)
		return err
	}

	return nil
}

// GetUserPassword 查询用户密码信息
func (u *userRepo) GetUserPasswordInfo(ctx context.Context, req *v1.WebLoginRequest) (*v1.UserPasswordInfo, error) {
	var userPassword UserPassword

	loginType := u.CheckLoginIdentityType(req.LoginIdentity)

	query := u.data.db.Table("t_user_password as t1").Joins("inner join t_user_identity as t2 on t1.user_id = t2.user_id").Select("t1.salt, t1.password").Where("t2.identity_type = ?", "password").Where("t2.identity_id = ?", req.LoginIdentity).Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0)

	switch loginType {
	case v1.LoginType_LOGIN_TYPE_EMAIL:
		query.Where("t2.identity_type = ?", "email")
	case v1.LoginType_LOGIN_TYPE_SMS:
		query.Where("t2.identity_type = ?", "sms")
	case v1.LoginType_LOGIN_TYPE_ACCOUNT:
		query.Where("t2.identity_type = ?", "account")
	}

	if err := query.First(&userPassword).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "密码凭证不存在")
		}
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}
	return &v1.UserPasswordInfo{Slat: userPassword.Salt, Password: userPassword.Password}, nil
}

func (u *userRepo) CheckUserPassword(ctx context.Context, req *v1.WebLoginRequest) (bool, error) {
	userPasswordInfo, err := u.GetUserPasswordInfo(ctx, req)
	if err != nil {
		u.log.Error("get user password info failed: %v", err)
		return false, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户名或密码错误, 请检查后重试")
	}

	pwd, err := password.New(req.Password, userPasswordInfo.Slat)
	if err != nil {
		u.log.Error("check password failed: %v", err)
		return false, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return pwd == userPasswordInfo.Password, nil
}

// WebLoginAccount 账号密码登录
func (u *userRepo) WebLoginAccount(ctx context.Context, req *v1.WebLoginRequest, loginValidPeriod int64) (*v1.WebLoginReply, error) {
	if loginValidPeriod <= 0 {
		loginValidPeriod = 8
	}

	if ok, err := u.CheckUserPassword(ctx, req); err != nil {
		u.log.Error("check user password failed: %v", err)
		return nil, err
	} else if !ok {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "账号或密码错误, 请稍后再试")
	}

	// 查询用户信息
	var userInfo UserInfo

	if err := u.data.db.Table("t_user as t1").Joins("inner join t_user_identity as t2 on t1.id = t2.user_id").
		Select("t1.id, t1.nickname").
		Where("t2.identity_type = ?", "password").
		Where("t2.identity_id = ?", req.LoginIdentity).
		Where("t1.deleted_flag = ?", 0).
		Where("t2.deleted_flag = ?", 0).
		First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户不存在")
		}
		u.log.Error("get user identity failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	jwtConfig := &jwt.Config{
		SecretKey: []byte(u.data.jwt.SecretKey),
		Issuer:    "xiaomiao-home-system",
		TTL:       time.Duration(loginValidPeriod) * time.Hour,
	}

	// 生成jwt token
	token, err := jwt.GenerateToken(jwtConfig, userInfo.Id, userInfo.Nickname)
	if err != nil {
		u.log.Error("generate jwt token failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.WebLoginReply{
		Code:    200,
		Message: "登录成功",
		Success: true,
		Data: &v1.WebLoginInfo{
			Token:            token,
			LoginValidPeriod: loginValidPeriod,
		},
	}, nil
}

func (u *userRepo) WebLoginSms(ctx context.Context, req *v1.WebLoginRequest) (*v1.WebLoginReply, error) {
	return nil, nil
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

func (u *userRepo) GetUserByNickname(ctx context.Context, nickname string) (*v1.UserInfo, error) {
	var user User

	if err := u.data.db.Table("t_user as t1").Joins("inner join t_user_password as t2 on t1.id = t2.user_id").Select("t1.id, t1.nickname, t1.avatar, t1.signature, t1.status").Where("t1.nickname = ?", nickname).Where("t1.deleted_flag = ?", 0).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			u.log.Error("user not found: %v", err)
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户不存在")
		}
		u.log.Error("get user by nickname failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.UserInfo{
		Id:        user.Id,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Signature: user.Signature,
		Status:    int32(user.Status),
	}, nil
}

// AppLogin 登录接口
func (u *userRepo) AppLogin(ctx context.Context, req *v1.AppLoginRequest) (*v1.AppLoginReply, error) {
	return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "暂不支持当前登录方式, 请联系管理员")
}

// MpLogin 登录接口
func (u *userRepo) MpLogin(ctx context.Context, req *v1.MpLoginRequest) (*v1.MpLoginReply, error) {
	return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "暂不支持当前登录方式, 请联系管理员")
}
