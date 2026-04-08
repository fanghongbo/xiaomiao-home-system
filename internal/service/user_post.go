package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/post/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserPostService struct {
	pb.UnimplementedUserPostServer

	userPost *biz.UserPostUsecase
	log      *log.Helper
	config   *conf.Config
}

func NewUserPostService(userPost *biz.UserPostUsecase, config *conf.Config, logger log.Logger) *UserPostService {
	return &UserPostService{
		userPost: userPost,
		config:   config,
		log:      log.NewHelper(log.With(logger, "service", "UserPostService")),
	}
}

// GetUserPostList 查询发布内容列表
func (s *UserPostService) GetUserPostList(ctx context.Context, req *pb.GetUserPostListRequest) (*pb.GetUserPostListReply, error) {
	res, err := s.userPost.GetUserPostList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserPost 创建发布内容
func (s *UserPostService) CreateUserPost(ctx context.Context, req *pb.CreateUserPostRequest) (*pb.CreateUserPostReply, error) {
	res, err := s.userPost.CreateUserPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserPost 更新发布内容
func (s *UserPostService) UpdateUserPost(ctx context.Context, req *pb.UpdateUserPostRequest) (*pb.UpdateUserPostReply, error) {
	res, err := s.userPost.UpdateUserPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserPost 删除发布内容
func (s *UserPostService) DeleteUserPost(ctx context.Context, req *pb.DeleteUserPostRequest) (*pb.DeleteUserPostReply, error) {
	res, err := s.userPost.DeleteUserPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserPostStatus 更新发布内容状态
func (s *UserPostService) UpdateUserPostStatus(ctx context.Context, req *pb.UpdateUserPostStatusRequest) (*pb.UpdateUserPostStatusReply, error) {
	res, err := s.userPost.UpdateUserPostStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserPost 查询发布内容
func (s *UserPostService) GetUserPost(ctx context.Context, req *pb.GetUserPostRequest) (*pb.GetUserPostReply, error) {
	res, err := s.userPost.GetUserPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
