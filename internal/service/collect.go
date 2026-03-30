package service

import (
	"context"
	pb "xiaomiao-home-system/api/collect/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type CollectService struct {
	pb.UnimplementedCollectServer

	collect *biz.CollectUsecase
	log     *log.Helper
	config  *conf.Config
}

func NewCollectService(collect *biz.CollectUsecase, config *conf.Config, logger log.Logger) *CollectService {
	return &CollectService{
		collect: collect,
		config:  config,
		log:     log.NewHelper(log.With(logger, "service", "CollectService")),
	}
}

// GetCollectList 查询收藏列表
func (s *CollectService) GetCollectList(ctx context.Context, req *pb.GetCollectListRequest) (*pb.GetCollectListReply, error) {
	res, err := s.collect.GetCollectList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetCollectTypes 查询收藏分类
func (s *CollectService) GetCollectTypes(ctx context.Context, req *pb.GetCollectTypesRequest) (*pb.GetCollectTypesReply, error) {
	res, err := s.collect.GetCollectTypes(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
