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
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	res, err := s.user.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
