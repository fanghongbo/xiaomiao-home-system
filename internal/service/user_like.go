package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/like/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserLikeService struct {
	pb.UnimplementedUserLikeServer

	userLike *biz.UserLikeUsecase
	log         *log.Helper
	config      *conf.Config
}

func NewUserLikeService(userLike *biz.UserLikeUsecase, config *conf.Config, logger log.Logger) *UserLikeService {
	return &UserLikeService{
		userLike: userLike,
		config:      config,
		log:         log.NewHelper(log.With(logger, "service", "UserLikeService")),
	}
}

// AddUserLike 添加用户喜欢
func (s *UserLikeService) AddUserLike(ctx context.Context, req *pb.AddUserLikeRequest) (*pb.AddUserLikeReply, error) {
	res, err := s.userLike.AddUserLike(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CancelUserLike 取消用户喜欢
func (s *UserLikeService) CancelUserLike(ctx context.Context, req *pb.CancelUserLikeRequest) (*pb.CancelUserLikeReply, error) {
	res, err := s.userLike.CancelUserLike(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
