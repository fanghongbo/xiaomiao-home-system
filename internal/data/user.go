package data

import (
	"context"
	"fmt"
	"regexp"
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

type AccountType string

const (
	// 账号类型
	ACCOUNT_TYPE_USERNAME AccountType = "username"
	ACCOUNT_TYPE_EMAIL    AccountType = "email"
	ACCOUNT_TYPE_SMS      AccountType = "sms"
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

// 判断登录账号类型, 允许用户使用账号、手机号、邮箱作为账号登陆名
func (u *userRepo) CheckLoginAccountType(loginIdentity string) AccountType {

	// 检查登录身份是否是邮箱地址
	if utils.IsEmail(loginIdentity) {
		return ACCOUNT_TYPE_EMAIL
	}

	// 检查登录身份是否是手机号
	if utils.IsMobile(loginIdentity) {
		return ACCOUNT_TYPE_SMS
	}

	// 检查登录身份是否是账号
	return ACCOUNT_TYPE_USERNAME
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
	case v1.LoginType_LOGIN_TYPE_ACCOUNT: // 账户登陆
		res, err = u.WebLoginAccount(ctx, req, loginValidPeriod)
	case v1.LoginType_LOGIN_TYPE_SMS: // 短信登陆
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
func (u *userRepo) GetUserPasswordInfo(ctx context.Context, accountType AccountType, req *v1.WebLoginRequest) (*v1.UserPasswordInfo, error) {
	var userPassword UserPassword

	query := u.data.db.Table("t_user_password as t1").Joins("inner join t_user_identity as t2 on t1.user_id = t2.user_id").Select("t1.salt, t1.password").Where("t2.identity_id = ?", req.LoginIdentity).Where("t1.deleted_flag = ?", 0).Where("t2.deleted_flag = ?", 0)

	switch accountType {
	case ACCOUNT_TYPE_EMAIL:
		query = query.Where("t2.identity_type = ?", "email")
	case ACCOUNT_TYPE_SMS:
		query = query.Where("t2.identity_type = ?", "sms")
	case ACCOUNT_TYPE_USERNAME:
		query = query.Where("t2.identity_type = ?", "password")
	}

	if err := query.First(&userPassword).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "密码凭证不存在")
		}
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}
	return &v1.UserPasswordInfo{Slat: userPassword.Salt, Password: userPassword.Password}, nil
}

func (u *userRepo) CheckUserPassword(ctx context.Context, accountType AccountType, req *v1.WebLoginRequest) (bool, error) {
	userPasswordInfo, err := u.GetUserPasswordInfo(ctx, accountType, req)
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

// 获取web登陆用户信息
func (u *userRepo) GetWebLoginAccountInfo(ctx context.Context, accountType AccountType, req *v1.WebLoginRequest) (*UserInfo, error) {
	var userInfo *UserInfo

	query := u.data.db.Table("t_user as t1").Joins("inner join t_user_identity as t2 on t1.id = t2.user_id").
		Select("t1.id, t1.nickname").
		Where("t2.identity_id = ?", req.LoginIdentity).
		Where("t1.deleted_flag = ?", 0).
		Where("t2.deleted_flag = ?", 0)

	switch accountType {
	case ACCOUNT_TYPE_EMAIL:
		query = query.Where("t2.identity_type = ?", "email")
	case ACCOUNT_TYPE_SMS:
		query = query.Where("t2.identity_type = ?", "sms")
	case ACCOUNT_TYPE_USERNAME:
		query = query.Where("t2.identity_type = ?", "password")
	}

	if err := query.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户不存在")
		}
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return userInfo, nil
}

// WebLoginAccount 账号密码登录
func (u *userRepo) WebLoginAccount(ctx context.Context, req *v1.WebLoginRequest, loginValidPeriod int64) (*v1.WebLoginReply, error) {
	if loginValidPeriod <= 0 {
		loginValidPeriod = 8
	}

	// 账号类型
	accountType := u.CheckLoginAccountType(req.LoginIdentity)

	if ok, err := u.CheckUserPassword(ctx, accountType, req); err != nil {
		u.log.Error("check user password failed: %v", err)
		return nil, err
	} else if !ok {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "账号或密码错误, 请稍后再试")
	}

	// 查询用户信息
	userInfo, err := u.GetWebLoginAccountInfo(ctx, accountType, req)
	if err != nil {
		u.log.Error("get web login user info failed: %v", err)
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

func (u *userRepo) GetUserByNickname(ctx context.Context, nickname string) (*v1.UserInfo, error) {
	var user User

	if err := u.data.db.Table("t_user as t1").Joins("inner join t_user_password as t2 on t1.id = t2.user_id").Select("t1.id, t1.nickname, t1.avatar, t1.gender, t1.birthday, t1.signature, t1.status").Where("t1.nickname = ?", nickname).Where("t1.deleted_flag = ?", 0).First(&user).Error; err != nil {
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
		Gender:    int32(user.Gender),
		Birthday:  user.Birthday.Format("2006-01-02"),
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

// 查询用户信息
func (u *userRepo) GetUserInfo(ctx context.Context, userId int64) (*v1.UserInfo, error) {
	var user User
	if err := u.data.db.Table("t_user").Select("id, nickname, avatar, gender, birthday, signature, status").Where("id = ?", userId).Where("deleted_flag = ?", 0).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(v1.ErrorReason_ERR_BAD_REQUEST.String(), "用户不存在")
		}
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.UserInfo{
		Id:        user.Id,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Gender:    int32(user.Gender),
		Birthday:  user.Birthday.Format("2006-01-02"),
		Signature: user.Signature,
		Status:    int32(user.Status),
	}, nil
}

// GetWebLoginUserInfo 查询登陆用户信息
func (u *userRepo) GetWebLoginUserInfo(ctx context.Context, req *v1.GetWebLoginUserInfoRequest) (*v1.GetWebLoginUserInfoReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	userInfo, err := u.GetUserInfo(ctx, userId)
	if err != nil {
		u.log.Error("get user info failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.GetWebLoginUserInfoReply{
		Code:    200,
		Message: "查询成功",
		Success: true,
		Data: &v1.WebLoginUserInfo{
			Id:          userInfo.Id,
			Nickname:    userInfo.Nickname,
			Avatar:      userInfo.Avatar,
			Gender:      int32(userInfo.Gender),
			Birthday:    userInfo.Birthday,
			Signature:   userInfo.Signature,
			AccessCodes: []string{},
		},
	}, nil
}

// WebLogout 退出登陆
func (u *userRepo) WebLogout(ctx context.Context, req *v1.WebLogoutRequest) (*v1.WebLogoutReply, error) {
	// 获取当前用户ID（用于日志记录）
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
	} else {
		u.log.Infof("user logout: userId=%d", userId)
	}

	// TODO: 如果需要实现token黑名单，可以在这里将当前token加入Redis黑名单
	// 黑名单的key可以是token本身，过期时间设置为token的剩余有效期

	return &v1.WebLogoutReply{
		Code:    200,
		Message: "退出登录成功",
		Success: true,
		Data:    "",
	}, nil
}

// WebCheckLogin web端登陆检测
func (u *userRepo) WebCheckLogin(ctx context.Context, req *v1.WebCheckLoginRequest) (*v1.WebCheckLoginReply, error) {
	// 由于此接口需要JWT认证，如果能执行到这里说明用户已登录
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_SESSION.String(), "未登录")
	}

	u.log.Infof("user check login: userId=%d", userId)

	return &v1.WebCheckLoginReply{
		Code:    200,
		Message: "已登录",
		Success: true,
		Data:    "",
	}, nil
}

// UpdateUserBaseSetting 更新用户基础设置
func (u *userRepo) UpdateUserBaseSetting(ctx context.Context, req *v1.UpdateUserBaseSettingRequest) (*v1.UpdateUserBaseSettingReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	userInfo := map[string]interface{}{
		"nickname":  req.Nickname,
		"gender":    req.Gender,
		"birthday":  req.Birthday,
		"signature": req.Signature,
	}

	if err := u.data.db.Table("t_user").Where("id = ?", userId).Updates(userInfo).Error; err != nil {
		u.log.Error("update user base setting failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserBaseSettingReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}

// CheckPassword 校验密码：必须包含 字母+数字+特殊符号，长度6-32
func (u *userRepo) CheckPassword(password string) bool {
	if len(password) < 6 || len(password) > 32 {
		return false
	}

	// 包含字母
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	// 包含数字
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 包含特殊符号
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	return hasLetter && hasDigit && hasSpecial
}

// UpdateUserPassword 更新用户密码
func (u *userRepo) UpdateUserPassword(ctx context.Context, req *v1.UpdateUserPasswordRequest) (*v1.UpdateUserPasswordReply, error) {
	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		u.log.Error("get current user id failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	if !u.CheckPassword(req.Password) {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "密码必须包含字母+数字+特殊符号，长度6-32")
	}

	salt, err := password.NewSalt(10)
	if err != nil {
		u.log.Error("generate salt failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	passwordHash, err := password.New(req.Password, salt)
	if err != nil {
		u.log.Error("generate password hash failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	if err := u.data.db.Table("t_user_password").Where("user_id = ?", userId).Updates(map[string]interface{}{
		"password": passwordHash,
		"salt":     salt,
	}).Error; err != nil {
		u.log.Error("update user password failed: %v", err)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_BAD_REQUEST.String(), "系统错误, 请稍后再试")
	}

	return &v1.UpdateUserPasswordReply{
		Code:    200,
		Message: "更新成功",
		Success: true,
		Data:    "",
	}, nil
}
