package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServer

	user   *biz.UserUsecase
	log    *log.Helper
	config *conf.Config
}

func NewUserService(user *biz.UserUsecase, config *conf.Config, logger log.Logger) *UserService {
	return &UserService{
		user:   user,
		config: config,
		log:    log.NewHelper(log.With(logger, "service", "UserService")),
	}
}

// Login 登录接口
func (s *UserService) WebLogin(ctx context.Context, req *pb.WebLoginRequest) (*pb.WebLoginReply, error) {
	return s.user.WebLogin(ctx, req)
}

// AppLogin 登录接口
func (s *UserService) AppLogin(ctx context.Context, req *pb.AppLoginRequest) (*pb.AppLoginReply, error) {
	return s.user.AppLogin(ctx, req)
}

// MpLogin 登录接口
func (s *UserService) MpLogin(ctx context.Context, req *pb.MpLoginRequest) (*pb.MpLoginReply, error) {
	return s.user.MpLogin(ctx, req)
}

// GetWebLoginUserInfo 查询登陆用户信息
func (s *UserService) GetWebLoginUserInfo(ctx context.Context, req *pb.GetWebLoginUserInfoRequest) (*pb.GetWebLoginUserInfoReply, error) {
	return s.user.GetWebLoginUserInfo(ctx, req)
}
