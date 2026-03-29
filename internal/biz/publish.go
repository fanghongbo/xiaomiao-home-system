package biz

import (
	"context"
	v1 "xiaomiao-home-system/api/publish/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// PublishRepo is a Greater repo.
type PublishRepo interface {
	// GetPublishList 查询发布内容列表
	GetPublishList(ctx context.Context, req *v1.GetPublishListRequest) (*v1.GetPublishListReply, error)
	// CreatePublish 创建发布内容
	CreatePublish(ctx context.Context, req *v1.CreatePublishRequest) (*v1.CreatePublishReply, error)
	// UpdatePublish 更新发布内容
	UpdatePublish(ctx context.Context, req *v1.UpdatePublishRequest) (*v1.UpdatePublishReply, error)
	// DeletePublish 删除发布内容
	DeletePublish(ctx context.Context, req *v1.DeletePublishRequest) (*v1.DeletePublishReply, error)
	// UpdatePublishStatus 更新发布内容状态
	UpdatePublishStatus(ctx context.Context, req *v1.UpdatePublishStatusRequest) (*v1.UpdatePublishStatusReply, error)
}

// PublishUsecase is a Publish usecase.
type PublishUsecase struct {
	repo PublishRepo
	log  *log.Helper
}

// NewPublishUsecase new a Publish usecase.
func NewPublishUsecase(repo PublishRepo, logger log.Logger) *PublishUsecase {
	return &PublishUsecase{repo: repo, log: log.NewHelper(log.With(logger, "biz", "PublishUsecase"))}
}

// GetPublishList 查询发布内容列表
func (u *PublishUsecase) GetPublishList(ctx context.Context, req *v1.GetPublishListRequest) (*v1.GetPublishListReply, error) {
	res, err := u.repo.GetPublishList(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePublish 创建发布内容
func (u *PublishUsecase) CreatePublish(ctx context.Context, req *v1.CreatePublishRequest) (*v1.CreatePublishReply, error) {
	res, err := u.repo.CreatePublish(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePublish 更新发布内容
func (u *PublishUsecase) UpdatePublish(ctx context.Context, req *v1.UpdatePublishRequest) (*v1.UpdatePublishReply, error) {
	res, err := u.repo.UpdatePublish(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeletePublish 删除发布内容
func (u *PublishUsecase) DeletePublish(ctx context.Context, req *v1.DeletePublishRequest) (*v1.DeletePublishReply, error) {
	res, err := u.repo.DeletePublish(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdatePublishStatus 更新发布内容状态
func (u *PublishUsecase) UpdatePublishStatus(ctx context.Context, req *v1.UpdatePublishStatusRequest) (*v1.UpdatePublishStatusReply, error) {
	res, err := u.repo.UpdatePublishStatus(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
