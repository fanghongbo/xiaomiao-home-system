package service

import (
	"context"
	pb "xiaomiao-home-system/api/user/group/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type UserGroupService struct {
	pb.UnimplementedUserGroupServer

	usergroup *biz.UserGroupUsecase
	log       *log.Helper
	config    *conf.Config
}

func NewUserGroupService(usergroup *biz.UserGroupUsecase, config *conf.Config, logger log.Logger) *UserGroupService {
	return &UserGroupService{
		usergroup: usergroup,
		config:    config,
		log:       log.NewHelper(log.With(logger, "service", "UserGroupService")),
	}
}

// GetUserGroupList 查询用户组列表
func (s *UserGroupService) GetUserGroupList(ctx context.Context, req *pb.GetUserGroupListRequest) (*pb.GetUserGroupListReply, error) {
	res, err := s.usergroup.GetUserGroupList(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserGroup 创建用户组
func (s *UserGroupService) CreateUserGroup(ctx context.Context, req *pb.CreateUserGroupRequest) (*pb.CreateUserGroupReply, error) {
	res, err := s.usergroup.CreateUserGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserGroup 更新用户组
func (s *UserGroupService) UpdateUserGroup(ctx context.Context, req *pb.UpdateUserGroupRequest) (*pb.UpdateUserGroupReply, error) {
	res, err := s.usergroup.UpdateUserGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserGroup 删除用户组
func (s *UserGroupService) DeleteUserGroup(ctx context.Context, req *pb.DeleteUserGroupRequest) (*pb.DeleteUserGroupReply, error) {
	res, err := s.usergroup.DeleteUserGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserGroupStatus 更新用户组状态
func (s *UserGroupService) UpdateUserGroupStatus(ctx context.Context, req *pb.UpdateUserGroupStatusRequest) (*pb.UpdateUserGroupStatusReply, error) {
	res, err := s.usergroup.UpdateUserGroupStatus(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserGroups 查询所有用户组
func (s *UserGroupService) GetUserGroups(ctx context.Context, req *pb.GetUserGroupsRequest) (*pb.GetUserGroupsReply, error) {
	res, err := s.usergroup.GetUserGroups(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
