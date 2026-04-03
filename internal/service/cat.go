package service

import (
	"context"
	pb "xiaomiao-home-system/api/cat/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type CatService struct {
	pb.UnimplementedCatServer

	cat    *biz.CatUsecase
	log    *log.Helper
	config *conf.Config
}

func NewCatService(cat *biz.CatUsecase, config *conf.Config, logger log.Logger) *CatService {
	return &CatService{
		cat:    cat,
		config: config,
		log:    log.NewHelper(log.With(logger, "service", "CatService")),
	}
}

// GetCatList 查询小猫列表
func (s *CatService) GetCatList(ctx context.Context, req *pb.GetCatListRequest) (*pb.GetCatListReply, error) {
	res, err := s.cat.GetCatList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
