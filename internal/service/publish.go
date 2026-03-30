package service

import (
	"context"
	pb "xiaomiao-home-system/api/publish/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type PublishService struct {
	pb.UnimplementedPublishServer

	publish *biz.PublishUsecase
	log     *log.Helper
	config  *conf.Config
}

func NewPublishService(publish *biz.PublishUsecase, config *conf.Config, logger log.Logger) *PublishService {
	return &PublishService{
		publish: publish,
		config:  config,
		log:     log.NewHelper(log.With(logger, "service", "PublishService")),
	}
}

// GetPublishList 查询发布内容列表
func (s *PublishService) GetPublishList(ctx context.Context, req *pb.GetPublishListRequest) (*pb.GetPublishListReply, error) {
	res, err := s.publish.GetPublishList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePublish 创建发布内容
func (s *PublishService) CreatePublish(ctx context.Context, req *pb.CreatePublishRequest) (*pb.CreatePublishReply, error) {
	res, err := s.publish.CreatePublish(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePublish 更新发布内容
func (s *PublishService) UpdatePublish(ctx context.Context, req *pb.UpdatePublishRequest) (*pb.UpdatePublishReply, error) {
	res, err := s.publish.UpdatePublish(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeletePublish 删除发布内容
func (s *PublishService) DeletePublish(ctx context.Context, req *pb.DeletePublishRequest) (*pb.DeletePublishReply, error) {
	res, err := s.publish.DeletePublish(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePublishStatus 更新发布内容状态
func (s *PublishService) UpdatePublishStatus(ctx context.Context, req *pb.UpdatePublishStatusRequest) (*pb.UpdatePublishStatusReply, error) {
	res, err := s.publish.UpdatePublishStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetPublish 查询发布内容
func (s *PublishService) GetPublish(ctx context.Context, req *pb.GetPublishRequest) (*pb.GetPublishReply, error) {
	res, err := s.publish.GetPublish(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
