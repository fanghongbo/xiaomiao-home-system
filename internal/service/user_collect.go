package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/collect/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserCollectService struct {
	pb.UnimplementedUserCollectServer

	userCollect *biz.UserCollectUsecase
	log         *log.Helper
	config      *conf.Config
}

func NewUserCollectService(userCollect *biz.UserCollectUsecase, config *conf.Config, logger log.Logger) *UserCollectService {
	return &UserCollectService{
		userCollect: userCollect,
		config:      config,
		log:         log.NewHelper(log.With(logger, "service", "UserCollectService")),
	}
}

// GetUserCollectList 查询用户收藏列表
func (s *UserCollectService) GetUserCollectList(ctx context.Context, req *pb.GetUserCollectListRequest) (*pb.GetUserCollectListReply, error) {
	res, err := s.userCollect.GetUserCollectList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserCollectTypes 查询用户收藏分类
func (s *UserCollectService) GetUserCollectTypes(ctx context.Context, req *pb.GetUserCollectTypesRequest) (*pb.GetUserCollectTypesReply, error) {
	res, err := s.userCollect.GetUserCollectTypes(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
