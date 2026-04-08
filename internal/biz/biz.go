package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase, NewRoleUsecase, NewUserPostUsecase, NewUserNotificationUsecase, NewUserSettingUsecase, NewFileUsecase, NewUserCollectUsecase, NewDiscoverUsecase, NewUserCatUsecase)
