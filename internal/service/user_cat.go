package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/cat/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserCatService struct {
	pb.UnimplementedUserCatServer

	userCat *biz.UserCatUsecase
	log     *log.Helper
	config  *conf.Config
}

func NewUserCatService(userCat *biz.UserCatUsecase, config *conf.Config, logger log.Logger) *UserCatService {
	return &UserCatService{
		userCat: userCat,
		config:  config,
		log:     log.NewHelper(log.With(logger, "service", "UserCatService")),
	}
}

// GetUserCatList 查询我的小猫列表
func (s *UserCatService) GetUserCatList(ctx context.Context, req *pb.GetUserCatListRequest) (*pb.GetUserCatListReply, error) {
	res, err := s.userCat.GetUserCatList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserCat 创建我的小猫
func (s *UserCatService) CreateUserCat(ctx context.Context, req *pb.CreateUserCatRequest) (*pb.CreateUserCatReply, error) {
	res, err := s.userCat.CreateUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateUserCat 更新我的小猫
func (s *UserCatService) UpdateUserCat(ctx context.Context, req *pb.UpdateUserCatRequest) (*pb.UpdateUserCatReply, error) {
	res, err := s.userCat.UpdateUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteUserCat 删除我的小猫
func (s *UserCatService) DeleteUserCat(ctx context.Context, req *pb.DeleteUserCatRequest) (*pb.DeleteUserCatReply, error) {
	res, err := s.userCat.DeleteUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetUserCat 查询我的小猫信息
func (s *UserCatService) GetUserCat(ctx context.Context, req *pb.GetUserCatRequest) (*pb.GetUserCatReply, error) {
	res, err := s.userCat.GetUserCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
