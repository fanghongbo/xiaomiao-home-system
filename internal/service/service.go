package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService, NewUserNotificationService, NewUserSettingService, NewFileService, NewUserPostService, NewUserCollectService, NewDiscoverService, NewUserCatService, NewUserLikeService)
