package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/role/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// RoleRepo is a Greater repo.
type RoleRepo interface {
	// GetRoleList 查询角色列表
	GetRoleList(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRoleListReply, error)
	// CreateRole 创建角色
	CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.CreateRoleReply, error)
	// UpdateRole 更新角色
	UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.UpdateRoleReply, error)
	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.DeleteRoleReply, error)
	// UpdateRoleStatus 更新角色状态
	UpdateRoleStatus(ctx context.Context, req *v1.UpdateRoleStatusRequest) (*v1.UpdateRoleStatusReply, error)
	// GetRoles 查询所有角色
	GetRoles(ctx context.Context, req *v1.GetRolesRequest) (*v1.GetRolesReply, error)
}

// RoleUsecase is a Role usecase.
type RoleUsecase struct {
	repo RoleRepo
	log  *log.Helper
}

// NewRoleUsecase new a Role usecase.
func NewRoleUsecase(repo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "RoleUsecase"))}
}

// GetRoleList 查询角色列表
func (u *RoleUsecase) GetRoleList(ctx context.Context, req *v1.GetRoleListRequest) (*v1.GetRoleListReply, error) {
	res, err := u.repo.GetRoleList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateRole 创建角色
func (u *RoleUsecase) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.CreateRoleReply, error) {
	res, err := u.repo.CreateRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateRole 更新角色
func (u *RoleUsecase) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.UpdateRoleReply, error) {
	res, err := u.repo.UpdateRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteRole 删除角色
func (u *RoleUsecase) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.DeleteRoleReply, error) {
	res, err := u.repo.DeleteRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateRoleStatus 更新角色状态
func (u *RoleUsecase) UpdateRoleStatus(ctx context.Context, req *v1.UpdateRoleStatusRequest) (*v1.UpdateRoleStatusReply, error) {
	res, err := u.repo.UpdateRoleStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRoles 查询所有角色
func (u *RoleUsecase) GetRoles(ctx context.Context, req *v1.GetRolesRequest) (*v1.GetRolesReply, error) {
	res, err := u.repo.GetRoles(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
