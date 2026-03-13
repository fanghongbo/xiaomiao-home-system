package biz

import (
	"context"

	v1 "xiaomiao-home-system/api/user/group/v1"

	"github.com/go-kratos/kratos/v2/log"
)

var (
// ErrUserGroupNotFound is user group not found.
// ErrUserGroupNotFound = errors.NotFound(v1.ErrorReason_USER_GROUP_NOT_FOUND.String(), "user group not found")
)

// UserGroupRepo is a Greater repo.
type UserGroupRepo interface {
	// GetUserGroupList 查询用户组列表
	GetUserGroupList(ctx context.Context, req *v1.GetUserGroupListRequest) (*v1.GetUserGroupListReply, error)
	// CreateUserGroup 创建用户组
	CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (*v1.CreateUserGroupReply, error)
	// UpdateUserGroup 更新用户组
	UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (*v1.UpdateUserGroupReply, error)
	// DeleteUserGroup 删除用户组
	DeleteUserGroup(ctx context.Context, req *v1.DeleteUserGroupRequest) (*v1.DeleteUserGroupReply, error)
	// UpdateUserGroupStatus 更新用户组状态
	UpdateUserGroupStatus(ctx context.Context, req *v1.UpdateUserGroupStatusRequest) (*v1.UpdateUserGroupStatusReply, error)
	// GetUserGroups 查询所有用户组
	GetUserGroups(ctx context.Context, req *v1.GetUserGroupsRequest) (*v1.GetUserGroupsReply, error)
}

// UserGroupUsecase is a User usecase.
type UserGroupUsecase struct {
	repo UserGroupRepo
	log  *log.Helper
}

// NewUserGroupUsecase new a UserGroup usecase.
func NewUserGroupUsecase(repo UserGroupRepo, logger log.Logger) *UserGroupUsecase {
	return &UserGroupUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "UserGroupUsecase"))}
}

// GetUserGroupList 查询用户组列表
func (u *UserGroupUsecase) GetUserGroupList(ctx context.Context, req *v1.GetUserGroupListRequest) (*v1.GetUserGroupListReply, error) {
	res, err := u.repo.GetUserGroupList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUserGroup 创建用户组
func (u *UserGroupUsecase) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (*v1.CreateUserGroupReply, error) {
	res, err := u.repo.CreateUserGroup(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserGroup 更新用户组
func (u *UserGroupUsecase) UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (*v1.UpdateUserGroupReply, error) {
	res, err := u.repo.UpdateUserGroup(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteUserGroup 删除用户组
func (u *UserGroupUsecase) DeleteUserGroup(ctx context.Context, req *v1.DeleteUserGroupRequest) (*v1.DeleteUserGroupReply, error) {
	res, err := u.repo.DeleteUserGroup(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUserGroupStatus 更新用户组状态
func (u *UserGroupUsecase) UpdateUserGroupStatus(ctx context.Context, req *v1.UpdateUserGroupStatusRequest) (*v1.UpdateUserGroupStatusReply, error) {
	res, err := u.repo.UpdateUserGroupStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserGroups 查询所有用户组
func (u *UserGroupUsecase) GetUserGroups(ctx context.Context, req *v1.GetUserGroupsRequest) (*v1.GetUserGroupsReply, error) {
	res, err := u.repo.GetUserGroups(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
