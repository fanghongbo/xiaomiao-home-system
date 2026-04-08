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

// CreateCat 创建小猫
func (s *CatService) CreateCat(ctx context.Context, req *pb.CreateCatRequest) (*pb.CreateCatReply, error) {
	res, err := s.cat.CreateCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCat 更新小猫
func (s *CatService) UpdateCat(ctx context.Context, req *pb.UpdateCatRequest) (*pb.UpdateCatReply, error) {
	res, err := s.cat.UpdateCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteCat 删除小猫
func (s *CatService) DeleteCat(ctx context.Context, req *pb.DeleteCatRequest) (*pb.DeleteCatReply, error) {
	res, err := s.cat.DeleteCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetCat 查询小猫信息
func (s *CatService) GetCat(ctx context.Context, req *pb.GetCatRequest) (*pb.GetCatReply, error) {
	res, err := s.cat.GetCat(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
