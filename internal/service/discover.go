package service

import (
	"context"
	pb "xiaomiao-home-system/api/discover/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type DiscoverService struct {
	pb.UnimplementedDiscoverServer

	discover *biz.DiscoverUsecase
	log      *log.Helper
	config   *conf.Config
}

func NewDiscoverService(discover *biz.DiscoverUsecase, config *conf.Config, logger log.Logger) *DiscoverService {
	return &DiscoverService{
		discover: discover,
		config:   config,
		log:      log.NewHelper(log.With(logger, "service", "DiscoverService")),
	}
}

// GetDiscoverList 查询发现列表
func (s *DiscoverService) GetDiscoverList(ctx context.Context, req *pb.GetDiscoverListRequest) (*pb.GetDiscoverListReply, error) {
	res, err := s.discover.GetDiscoverList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDiscover 查询发现内容
func (s *DiscoverService) GetDiscover(ctx context.Context, req *pb.GetDiscoverRequest) (*pb.GetDiscoverReply, error) {
	res, err := s.discover.GetDiscover(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDiscoverRecommend 查询推荐内容
func (s *DiscoverService) GetDiscoverRecommend(ctx context.Context, req *pb.GetDiscoverRecommendRequest) (*pb.GetDiscoverRecommendReply, error) {
	res, err := s.discover.GetDiscoverRecommend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetDiscoverRecommendExcludePostId 查询推荐内容，排除指定id
func (s *DiscoverService) GetDiscoverRecommendExcludePostId(ctx context.Context, req *pb.GetDiscoverRecommendRequest) (*pb.GetDiscoverRecommendReply, error) {
	res, err := s.discover.GetDiscoverRecommend(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
