package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase, NewUserPostUsecase, NewUserNotificationUsecase, NewUserSettingUsecase, NewFileUsecase, NewUserCollectUsecase, NewDiscoverUsecase, NewUserCatUsecase, NewUserLikeUsecase)
