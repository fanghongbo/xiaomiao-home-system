package service

import (
	"context"
	pb "xiaomiao-home-system/api/post/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type PostService struct {
	pb.UnimplementedPostServer

	post *biz.PostUsecase
	log     *log.Helper
	config  *conf.Config
}

func NewPostService(post *biz.PostUsecase, config *conf.Config, logger log.Logger) *PostService {
	return &PostService{
		post: post,
		config:  config,
		log:     log.NewHelper(log.With(logger, "service", "PostService")),
	}
}

// GetPostList 查询发布内容列表
func (s *PostService) GetPostList(ctx context.Context, req *pb.GetPostListRequest) (*pb.GetPostListReply, error) {
	res, err := s.post.GetPostList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePost 创建发布内容
func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostReply, error) {
	res, err := s.post.CreatePost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePost 更新发布内容
func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostReply, error) {
	res, err := s.post.UpdatePost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeletePost 删除发布内容
func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostReply, error) {
	res, err := s.post.DeletePost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePostStatus 更新发布内容状态
func (s *PostService) UpdatePostStatus(ctx context.Context, req *pb.UpdatePostStatusRequest) (*pb.UpdatePostStatusReply, error) {
	res, err := s.post.UpdatePostStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetPost 查询发布内容
func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostReply, error) {
	res, err := s.post.GetPost(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
