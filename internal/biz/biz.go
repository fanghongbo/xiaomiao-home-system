package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase, NewRoleUsecase, NewPublishUsecase, NewUserNotificationUsecase, NewUserSettingUsecase, NewFileUsecase, NewCollectUsecase)
