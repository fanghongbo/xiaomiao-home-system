package service

import (
	"context"
	pb "xiaomiao-home-system/api/role/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type RoleService struct {
	pb.UnimplementedRoleServer

	role   *biz.RoleUsecase
	log    *log.Helper
	config *conf.Config
}

func NewRoleService(role *biz.RoleUsecase, config *conf.Config, logger log.Logger) *RoleService {
	return &RoleService{
		role:   role,
		config: config,
		log:    log.NewHelper(log.With(logger, "service", "RoleService")),
	}
}

// GetRoleList 查询角色列表
func (s *RoleService) GetRoleList(ctx context.Context, req *pb.GetRoleListRequest) (*pb.GetRoleListReply, error) {
	res, err := s.role.GetRoleList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleReply, error) {
	res, err := s.role.CreateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {
	res, err := s.role.UpdateRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleReply, error) {
	res, err := s.role.DeleteRole(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateRoleStatus 更新角色状态
func (s *RoleService) UpdateRoleStatus(ctx context.Context, req *pb.UpdateRoleStatusRequest) (*pb.UpdateRoleStatusReply, error) {
	res, err := s.role.UpdateRoleStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRoles 查询所有角色
func (s *RoleService) GetRoles(ctx context.Context, req *pb.GetRolesRequest) (*pb.GetRolesReply, error) {
	res, err := s.role.GetRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
